package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type AcceptServer struct {
	Addr string
	Cert string
	Key  string
}

func (s *AcceptServer) Run(r *Runner, wg sync.WaitGroup) {
	l := s.listen()
	s.accept(l, r)
	wg.Done()
	defer l.Close()
}

func (s *AcceptServer) listen() (l net.Listener) {
	cer, err := tls.LoadX509KeyPair(s.Cert, s.Key)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	l, err = tls.Listen("tcp", s.Addr, config)
	if err != nil {
		log.Printf("[TCP] Listen on %s: %v", s.Addr, err)
		os.Exit(1)
	}

	log.Printf("[TCP] Listening on %s", s.Addr)

	return
}

func (s *AcceptServer) accept(l net.Listener, r *Runner) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("[TCP] Accept connection", err)
			continue
		}

		p := &Payload{
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		r.Payloads = append(r.Payloads, p)

		fmt.Printf("[*] New connection %s. Total connections: %d\n", p.Addr, len(r.Payloads))
	}
}
