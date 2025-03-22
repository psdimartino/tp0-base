package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config  ClientConfig
	conn    net.Conn
	running bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:  config,
		running: true,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *Client) closeClientSocket() error {
	if c.conn == nil {
		log.Infof("action: close_connection | result: success | client_id: %v | detail: already closed", c.config.ID)
		return nil
	}
	err := c.conn.Close()
	if err != nil {
		log.Infof("action: connection_close | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
	c.conn = nil
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {

	c.createSigtermHandler()

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount && c.running; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		err := c.createClientSocket()
		if err != nil {
			return
		}

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)

		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		err = c.closeClientSocket()
		if err != nil {
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) createSigtermHandler() {
	// Create the channel to listen to the signal
	sigChannel := make(chan os.Signal, 1)
	// Notify the channel: "Notify causes package signal to relay incoming signals to c."
	signal.Notify(sigChannel, syscall.SIGTERM)
	// Start the asynchronous task as a handler
	go c.handleSigterm(sigChannel)
}

func (c *Client) handleSigterm(sigChannel chan os.Signal) {
	// Wait for the signal. Doesn't catch the response as is always SIGTERM
	<-sigChannel
	log.Infof("action: handle_sigterm | result: success | client_id: %v", c.config.ID)
	c.running = false
	err := c.closeClientSocket()
	if err != nil {
		return
	}
}
