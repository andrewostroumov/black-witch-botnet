package proto

import (
	"encoding/binary"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"soulless_network/relations"
)

const (
	bufSize = 256
)

const (
	typeShellCommand = iota
	typeEventMessage
	typeShellResult
	typeEventResult
	typeErrorResult
)

type Package struct {
	Type  uint8
	Bytes []byte
}

type Reader struct {
	rd io.Reader
}

type Writer struct {
	wr io.Writer
}

func NewReader(rd io.Reader) *Reader {
	r := &Reader{rd}
	return r
}

func NewWriter(wr io.Writer) *Writer {
	w := &Writer{wr}
	return w
}

// TODO: add timeout to read tail
func (r *Reader) Read() (*Package, error) {
	buf := make([]byte, 4)
	_, err := r.rd.Read(buf)

	if err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(buf)

	buf = make([]byte, 1)
	_, err = r.rd.Read(buf)

	if err != nil {
		return nil, err
	}

	t := buf[0]

	pack := &Package{
		Type: t,
	}

	count := size / bufSize

	log.Println("proto read frame count", count)

	for i := uint32(0); i < count; i++ {
		buf = make([]byte, bufSize)
		_, err = r.rd.Read(buf)

		if err != nil {
			return nil, err
		}

		pack.Bytes = append(pack.Bytes, buf...)
	}

	mod := size % bufSize

	if mod != 0 {
		buf = make([]byte, mod)
		_, err = r.rd.Read(buf)

		if err != nil {
			return nil, err
		}

		pack.Bytes = append(pack.Bytes, buf...)
	}

	return pack, nil
}

func (w *Writer) Write(r *Package) error {
	buf := make([]byte, 4)
	size := uint32(len(r.Bytes))
	binary.BigEndian.PutUint32(buf, size)

	buf = append(buf, r.Type)
	buf = append(buf, r.Bytes...)

	_, err := w.wr.Write(buf)

	if err != nil {
		return err
	}

	return nil
}

func Unmarshal(p *Package) (interface{}, error) {
	var i interface{}

	switch p.Type {
	case typeShellCommand:
		i = &relations.ShellCommand{}
	case typeEventMessage:
		i = &relations.EventMessage{}
	case typeShellResult:
		i = &relations.ShellResult{}
	case typeEventResult:
		i = &relations.EventResult{}
	case typeErrorResult:
		i = &relations.ErrorResult{}
	default:
		return nil, errors.New("unknown package type")
	}

	err := bson.Unmarshal(p.Bytes, i)

	if err != nil {
		return nil, err
	}

	return i, nil
}

func Marshal(i interface{}) (*Package, error) {
	p := &Package{}

	if _, ok := i.(*relations.ShellCommand); ok {
		p.Type = typeShellCommand
	} else if _, ok := i.(*relations.EventMessage); ok {
		p.Type = typeEventMessage
	} else if _, ok := i.(*relations.ShellResult); ok {
		p.Type = typeShellResult
	} else if _, ok := i.(*relations.EventResult); ok {
		p.Type = typeEventResult
	} else if _, ok := i.(*relations.ErrorResult); ok {
		p.Type = typeErrorResult
	} else {
		return nil, errors.New("unknown package type")
	}

	b, err := bson.Marshal(i)

	if err != nil {
		return nil, err
	}

	p.Bytes = b

	return p, nil
}
