// Copyright 2022 The jackal Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xmppparser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jackal-xmpp/stravaganza/parser/internal/gosaxml"

	"github.com/jackal-xmpp/stravaganza"
)

const rootElementIndex = -1

var (
	streamName = []byte("stream")
)

// ParsingMode defines the way in which special parsed element
// should be considered or not according to the reader nature.
type ParsingMode int

const (
	// DefaultMode treats incoming elements as provided from raw byte reader.
	DefaultMode = ParsingMode(iota)

	// SocketStream treats incoming elements as provided from a socket transport.
	SocketStream
)

// ErrTooLargeStanza will be returned Parse when the size of the incoming stanza is too large.
var ErrTooLargeStanza = errors.New("parser: too large stanza")

// ErrStreamClosedByPeer will be returned by Parse when stream closed element is parsed.
var ErrStreamClosedByPeer = errors.New("parser: stream closed by peer")

// Parser parses arbitrary XML input and builds an array with the structure of all tag and data elements.
type Parser struct {
	dec           gosaxml.Decoder
	mode          ParsingMode
	nextElement   stravaganza.Element
	stack         []*stravaganza.Builder
	inElement     bool
	pIndex        int
	lastOffset    int
	maxStanzaSize int
}

// New creates an empty Parser instance.
func New(reader io.Reader, mode ParsingMode, maxStanzaSize int) *Parser {
	return &Parser{
		mode:          mode,
		dec:           gosaxml.NewDecoder(reader),
		pIndex:        rootElementIndex,
		maxStanzaSize: maxStanzaSize,
	}
}

// Parse parses next available XML element from reader.
func (p *Parser) Parse() (stravaganza.Element, error) {
	var tk gosaxml.Token

	if err := p.dec.NextToken(&tk); err != nil {
		return nil, err
	}
	for {
		// check max stanza size limit
		off := p.dec.InputOffset()
		if p.maxStanzaSize > 0 && off-p.lastOffset > p.maxStanzaSize {
			return nil, ErrTooLargeStanza
		}
		switch tk.Kind {
		case gosaxml.TokenTypeStartElement:
			p.startElement(&tk)
			if p.mode == SocketStream && isStreamName(tk.Name) {
				if err := p.closeElement(xmlName(tk.Name)); err != nil {
					return nil, err
				}
				goto done
			}

		case gosaxml.TokenTypeCharData:
			if p.inElement {
				p.setElementText(&tk)
			}

		case gosaxml.TokenTypeEndElement:
			if p.mode == SocketStream && isStreamName(tk.Name) {
				return nil, ErrStreamClosedByPeer
			}
			if err := p.endElement(&tk); err != nil {
				return nil, err
			}
			if p.pIndex == rootElementIndex {
				goto done
			}
		}
		if err := p.dec.NextToken(&tk); err != nil {
			return nil, err
		}
	}

done:
	p.lastOffset = p.dec.InputOffset()
	elem := p.nextElement
	p.nextElement = nil

	return elem, nil
}

func (p *Parser) startElement(tk *gosaxml.Token) {
	name := xmlName(tk.Name)

	attrs := make([]stravaganza.Attribute, 0, len(tk.Attr))
	for _, a := range tk.Attr {
		name := xmlName(a.Name)
		attrs = append(attrs, stravaganza.Attribute{Label: name, Value: string(a.Value)})
	}
	builder := stravaganza.NewBuilder(name).WithAttributes(attrs...)
	p.stack = append(p.stack, builder)

	p.pIndex = len(p.stack) - 1
	p.inElement = true
}

func (p *Parser) setElementText(tk *gosaxml.Token) {
	p.stack[p.pIndex] = p.stack[p.pIndex].WithText(string(tk.ByteData))
}

func (p *Parser) endElement(tk *gosaxml.Token) error {
	return p.closeElement(xmlName(tk.Name))
}

func (p *Parser) closeElement(name string) error {
	if p.pIndex == rootElementIndex {
		return errUnexpectedEnd(name)
	}
	builder := p.stack[p.pIndex]
	p.stack = p.stack[:p.pIndex]

	element := builder.Build()

	if name != element.Name() {
		return errUnexpectedEnd(name)
	}
	p.pIndex = len(p.stack) - 1
	if p.pIndex == rootElementIndex {
		p.nextElement = element
	} else {
		p.stack[p.pIndex] = p.stack[p.pIndex].WithChild(element)
	}
	p.inElement = false
	return nil
}

func xmlName(name gosaxml.Name) string {
	if len(name.Prefix) > 0 {
		var sb strings.Builder
		sb.Write(name.Prefix)
		sb.WriteRune(':')
		sb.Write(name.Local)
		return sb.String()
	}
	return string(name.Local)
}

func isStreamName(name gosaxml.Name) bool {
	return bytes.Compare(name.Local, streamName) == 0 && bytes.Compare(name.Prefix, streamName) == 0
}

func errUnexpectedEnd(name string) error {
	return fmt.Errorf("xmppparser: unexpected end element </%s>", name)
}
