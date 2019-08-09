package server

import (
	"errors"
	"strconv"
)

type Runner struct {
	Accept  *AcceptServer
	Control *ControlServer
	Payloads []*Payload
}

func (r *Runner) Run() {
	go r.Accept.Run(r)
	r.Control.Run(r)
    // Wait for exit
}

//func (c *Server) handleConnections() {
//	for {
//		fmt.Print("<CC:#> ")
//		// Read the stdin
//		stdreader := bufio.NewReader(os.Stdin)
//		text, err := stdreader.ReadString('\n')
//		if err != nil {
//			fmt.Println("[ERROR] reading std input", err)
//			continue
//		}
//
//		// Check the command issued
//		switch {
//		case strings.TrimSpace(text) == "show":
//			for i, p := range c.Payloads {
//				fmt.Printf("ID: %d Address: %s\n", i, p.Addr.String())
//			}
//		case strings.Contains(strings.TrimSpace(text), "use"):
//			index := strings.Split(text, " ")[1]
//			p, err := c.getPayload(strings.TrimSpace(index))
//			if err != nil {
//				fmt.Println("[ERROR] getting payload", err)
//				continue
//			}
//
//			p.Activate()
//		}
//	}
//}

func (r *Runner) getPayload(index string) (*Payload, error) {
	i, err := strconv.Atoi(index)


	if err != nil {
		return nil, err
	}

	if i < 0 || i >= len(r.Payloads) {
		return nil, errors.New("index out of range")
	}

	p := r.Payloads[i]

	if p == nil {
		return nil, errors.New("payload not found")
	}

	return p, nil
}