package server

import (
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net"
	"soulless_network/relations"
	"strings"
)

type Payload struct {
	Addr net.Addr
	Conn net.Conn
}

func (p *Payload) Activate(conn net.Conn) {
	for {
		pre := color.Blue.Text(fmt.Sprintf("<%s:#> ", p.Addr.String()))
		// 	writer := bufio.NewWriter(conn) doesn't work
		conn.Write([]byte(pre))

		text, err := receive(conn)

		if text == "exit" {
			break
		}

		msg := parse(text)

		// TODO: check if message is full
		log.Println(msg)

		if err != nil {
			log.Println("[PL] Reading unix input", err)
			break
		}

		err = p.handle(msg, conn)

		if err != nil {
			log.Println("[TCP] Handling", err)
			conn.Write([]byte(err.Error()))
			break
		}
	}
}

func (p *Payload) Write(msg *relations.Message) (error) {
	b, err := bson.Marshal(msg)

	if err != nil {
		return err
	}

	_, err = p.Conn.Write(append(b, '\r'))

	if err != nil {
		return err
	}

	return nil
}

func (p *Payload) Read() (*relations.Result, error) {
	reader := bufio.NewReader(p.Conn)
	b, err := reader.ReadBytes('\r')

	if err != nil {
		return nil, err
	}

	res := &relations.Result{}
	err = bson.Unmarshal(b, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Payload) handle(msg *relations.Message, conn net.Conn) (error) {
	err := p.Write(msg)

	if err != nil {
		return err
	}

	res, err := p.Read()

	if err != nil {
		return err
	}


	// TODO: stringify struct
	_, err = conn.Write(append([]byte(res.Data), '\n'))

	if err != nil {
		return err
	}

	return nil
}

func receive(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	text, err := reader.ReadString('\n')

	if err == io.EOF {
		return "", err
	}

	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	return text, nil
}

func parse(text string) (*relations.Message) {
	p := &Parser{text}
	return p.Parse()
}
