package common

import (
	"bufio"
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
	Bets          string
	MaxBatchSize  int
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

	file, err := os.Open(c.config.Bets)
	if err != nil {
		log.Fatalf("action: open_file | result: error | error: %v", err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("action: open_file | result: error | error: %v", err)
	}
	if fileInfo.Size() == 0 {
		log.Fatalf("action: open_file | result: error | msg: Empty file")
	}
	log.Infof("action: open_file | result: success | size: %d", fileInfo.Size())

	scanner := bufio.NewScanner(file)
	// Close the file when the function returns
	defer file.Close()

	err = c.createClientSocket()
	if err != nil && scanner.Scan() == false {
		log.Fatalf("action: create_socket | result: error | error: %v", err)
		return
	}
	var line string
	var maxBatch = c.config.MaxBatchSize
	var end = false
	// For all the file
	for !end {
		// While file didn't end and below batch max
		var i = 0
		var msg = ""
		for i < maxBatch {
			scanner.Scan()
			line := scanner.Text()
			if err := scanner.Err(); err != nil {
				log.Fatalf("action: scan_line | result: fail | error: %v", err)
				return
			}
			if len(line) == 0 {
				end = true
				break
			}
			log.Infof("action: scan_line | result: success | line: %s", line)
			msg = msg + c.config.ID + "," + line + "\n"
			i++
		}

		log.Infof("action: batch_leido | result: success | amount: %v", i)
		err = c.sendMessage(msg)

		if err != nil {
			log.Fatalf("action: apuesta_enviada | result: fail | cantidad: %v", i)
			break
		}

		log.Infof("action: batch_enviado | result: success | amount: %v", i)
		msg = c.receiveResponse()

		if msg == "err" {
			log.Fatalf("action: respuesta_recibida | result: fail")
			break
		}

		log.Infof("action: respuesta_recibida | result: success | response: %v", msg)

	}
	err = c.sendMessage("end")
	if err != nil {
		log.Fatalf("action: terminar_envio | result: fail | msg: %v", line)
		return
	}
	err = c.closeClientSocket()
	if err != nil {
		log.Fatalf("action: cerrar_socket | result: fail | msg: %v", line)
		return
	}
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

func (c *Client) sendMessage(message string) error {
	totalWritten := 0
	msg := message + "\n"
	messageBytes := []byte(msg)

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
