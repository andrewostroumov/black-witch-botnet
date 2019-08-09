package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Payload struct {
	Addr net.Addr
	Conn net.Conn
}

func (p *Payload) Activate(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		pre := fmt.Sprintf("<%s:#> ", p.Addr.String())
		conn.Write([]byte(pre))

		text, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("[PL] Reading unix input", err)
			continue
		}
		text = strings.TrimSpace(text)

		switch {
		case text == "exit":
			return
		default:
			p.execCommand(text, conn)
		}
	}
}

func (p *Payload) execCommand(text string, conn net.Conn) {
	reader := bufio.NewReader(p.Conn)
	_, err := p.Conn.Write([]byte(text + "\r"))

	if err != nil {
		s := fmt.Sprintf("Write connection %s\n", err)
		conn.Write([]byte(s))
		return
	}

	res, err := reader.ReadString('\r')

	if err != nil {
		s := fmt.Sprintf("Read connection %s\n", err)
		conn.Write([]byte(s))
		return
	}

	conn.Write([]byte(res))
}