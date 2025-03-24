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
	Nombre        string
	Apellido      string
	Documento     string
	Nacimiento    string
	Numero        string
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

	err := c.createClientSocket()
	if err != nil {
		return
	}

	err = c.sendMessage()
	if err != nil {
		return
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
		c.config.Documento,
		c.config.Numero,
	)

	msg := c.receiveResponse()

	err = c.closeClientSocket()

	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		msg,
	)

	log.Infof("action: send_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) receiveResponse() string {
	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return ""
	}
	return msg
}

func (c *Client) sendMessage() error {
	message := fmt.Sprintf("%s,%s,%s,%s,%s,%s\n",
		c.config.ID, c.config.Nombre, c.config.Apellido, c.config.Documento, c.config.Nacimiento, c.config.Numero)
	totalWritten := 0
	messageBytes := []byte(message)

	for totalWritten < len(messageBytes) {
		n, err := c.conn.Write(messageBytes[totalWritten:])
		if err != nil {
			return err
		}
		totalWritten += n
	}

	return nil
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
