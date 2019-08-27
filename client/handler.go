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
			Type: relations.TypeErrorResult,
			Data: &relations.ErrorResult{
				Code: 1,
				Data: err.Error(),
			},
		}
	} else {
		resp = &relations.Response{
			Type: relations.TypeShellResult,
			Data: res,
		}
	}

	return resp
}

func (h *Handler) exec() (*relations.ShellResult, error) {
	data := strings.Split(h.Command.Data, " ")
	cmd := data[0]
	args := append(data[:0], data[0+1:]...)

	e := exec.Command(cmd, args...)
	o, err := e.Output()

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			res := &relations.ShellResult{
				Exit:   ee.ExitCode(),
				Stderr: ee.Stderr,
			}

			return res, nil
		}

		return nil, err
	}

	res := &relations.ShellResult{
		Exit:   0,
		Stdout: o,
	}

	return res, nil
}
