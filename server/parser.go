package server

import (
	"soulless_network/relations"
	"strings"
)

type Parser struct {
	Text string
}

func (p *Parser) Parse() (*relations.Message) {
	msg := &relations.Message{}
	seg := strings.Split(p.Text, " ")

	i := 0

	if len(seg) == 0 {
		return msg
	}

	res, ok := p.Normalize(seg, relations.Types, &i, "shell")

	if ok {
		msg.Type = res
	} else {
		return msg
	}

	res, ok = p.Normalize(seg, relations.Domains, &i, "exec")

	if ok {
		msg.Domain = res
	} else {
		return msg
	}

	msg.Data = strings.Join(seg[i:], " ")

	return msg
}

func (p *Parser) Normalize(seg []string, enum []string, i *int, def string) (string, bool) {
	var res string

	if len(seg) < *i + 1 {
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
		res = def
	}

	return res, true
}
