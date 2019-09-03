package server

// TODO: next
// When client connect it should send info about itself
// Create a new response type - change dir and attach it to unix sock prefix
// hello message from server (as part as internal monitor maybe?)
// save logs to file
// add daemon file
// implement internal loop to check connections

import (
	"errors"
	"strconv"
	"sync"
)

type Runner struct {
	Accept   *AcceptServer
	Control  *ControlServer
	Payloads []*Payload
}

func (r *Runner) Run() {
	var wg sync.WaitGroup
	wg.Add(2)

	go r.Accept.Run(r, wg)
	go r.Control.Run(r, wg)

	wg.Wait()
}

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
