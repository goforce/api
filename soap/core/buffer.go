package core

import (
	"bytes"
	"encoding/xml"
	"github.com/goforce/api/commons"
)

type Buffer struct{ bytes.Buffer }

func (buf *Buffer) writeValue(tag string, value interface{}) error {
	buf.openTag(tag)
	s, err := commons.String(value)
	if err != nil {
		return err
	}
	err = xml.EscapeText(buf, []byte(s))
	if err != nil {
		return err
	}
	buf.closeTag(tag)
	return nil
}

func (buf *Buffer) writeElement(tag string, value string) {
	buf.openTag(tag)
	buf.WriteString(value)
	buf.closeTag(tag)
}

func (buf *Buffer) openTag(tag string) {
	buf.WriteByte('<')
	buf.WriteString(tag)
	buf.WriteByte('>')
}

func (buf *Buffer) closeTag(tag string) {
	buf.WriteString("</")
	buf.WriteString(tag)
	buf.WriteByte('>')
}
