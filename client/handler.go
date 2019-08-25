package client

import (
	"os/exec"
	"soulless_network/relations"
	"strings"
)

type Handler struct {
	Command *relations.Command
}

func (h *Handler) handle() *relations.Response {
	var resp *relations.Response
	res, err := h.exec()

	if err != nil {
		resp = &relations.Response{
			Error: &relations.Error{
				Code: 1,
				Data: err.Error(),
			},
		}
	} else {
		resp = &relations.Response{
			Result: res,
		}
	}

	return resp
}

func (h *Handler) exec() (*relations.Result, error) {
	data := strings.Split(h.Command.Data, " ")
	cmd := data[0]
	args := append(data[:0], data[0+1:]...)

	e := exec.Command(cmd, args...)
	o, err := e.Output()

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			res := &relations.Result{
				Exit:   ee.ExitCode(),
				Stderr: ee.Stderr,
			}

			return res, nil
		}

		return nil, err
	}

	res := &relations.Result{
		Exit:   0,
		Stdout: o,
	}

	return res, nil
}
