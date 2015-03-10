package soap

import (
	"bytes"
	"encoding/xml"
)

type KV struct {
	Key   string
	Value interface{}
}

func (p *KV) writeXml(b *bytes.Buffer) *bytes.Buffer {
	if b == nil {
		b = &bytes.Buffer{}
	}
	p.open(b)
	p.content(b)
	p.close(b)
	return b
}

func (p *KV) open(b *bytes.Buffer) *bytes.Buffer {
	b.WriteByte('<')
	b.WriteString(p.Key)
	b.WriteByte('>')
	return b
}

func (p *KV) close(b *bytes.Buffer) *bytes.Buffer {
	b.WriteString("</")
	b.WriteString(p.Key)
	b.WriteByte('>')
	return b
}

func (p *KV) content(b *bytes.Buffer) *bytes.Buffer {
	switch p.Value.(type) {
	case string:
		xml.EscapeText(b, []byte(p.Value.(string)))
	case XE:
		np := p.Value.(KV)
		np.writeXml(b)
	case []XE:
		for _, np := range p.Value.([]KV) {
			np.writeXml(b)
		}
	}
	return b
}
