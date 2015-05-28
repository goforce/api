package soap

import (
	"bytes"
	"encoding/xml"
)

type XE struct {
	Tag   string
	Value interface{}
}

func (e *XE) write(b *bytes.Buffer) *bytes.Buffer {
	if b == nil {
		b = &bytes.Buffer{}
	}
	e.open(b)
	e.content(b)
	e.close(b)
	return b
}

func (e *XE) open(b *bytes.Buffer) *bytes.Buffer {
	b.WriteByte('<')
	b.WriteString(e.Tag)
	b.WriteByte('>')
	return b
}

func (e *XE) close(b *bytes.Buffer) *bytes.Buffer {
	b.WriteString("</")
	b.WriteString(e.Tag)
	b.WriteByte('>')
	return b
}

func (e *XE) content(b *bytes.Buffer) *bytes.Buffer {
	switch e.Value.(type) {
	case string:
		xml.EscapeText(b, []byte(e.Value.(string)))
	case XE:
		ne := e.Value.(XE)
		ne.write(b)
	case []XE:
		for _, ne := range e.Value.([]XE) {
			ne.write(b)
		}
	}
	return b
}
