package gosaxml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Decoder decodes an XML input stream into Token values.
type Decoder interface {
	// NextToken decodes and stores the next Token into
	// the provided Token pointer.
	// Only the fields relevant for the decoded token type
	// are written to the Token. Other fields may have previous
	// values. The caller should thus determine the Token.Kind
	// and then only read/touch the fields relevant for that kind.
	NextToken(t *Token) error

	// InputOffset returns the current offset in the input stream.
	InputOffset() int

	// Reset resets the Decoder to the given io.Reader.
	Reset(r io.Reader)
}

type decoder struct {
	rb                  [2048]byte
	bbOffset            [256]int32
	numAttributes       [256]byte
	lastOpen            Name
	preserveWhitespaces [32]bool
	rd                  io.Reader
	bb                  []byte
	attrs               []Attr
	r                   int
	w                   int
	top                 byte
	lastStartElement    bool
}

var (
	bsxml      = []byte("xml")
	bsspace    = []byte("space")
	bspreserve = []byte("preserve")
	simdWidth  int
)

// NewDecoder creates a new Decoder.
func NewDecoder(r io.Reader) Decoder {
	return &decoder{
		rd:    r,
		bb:    make([]byte, 0, 256),
		attrs: make([]Attr, 0, 256),
	}
}

func isWhitespace(b byte) bool {
	return b <= ' '
}

func (d *decoder) read0() error {
	if d.r > 0 {
		copy(d.rb[:], d.rb[d.r:d.w])
		d.w -= d.r
		d.r = 0
	}
	n, err := d.rd.Read(d.rb[d.w : cap(d.rb)-simdWidth])
	d.w += n
	if n <= 0 && err != nil {
		return err
	}
	return nil
}

func (d *decoder) unreadByte() {
	d.r--
}

func (d *decoder) readByte() (byte, error) {
	for d.r == d.w {
		err := d.read0()
		if err != nil {
			return 0, err
		}
	}
	c := d.rb[d.r]
	d.r++
	return c, nil
}

func (d *decoder) discardBuffer() {
	d.r = d.w
}

func (d *decoder) discard(n int) (int, error) {
	for d.r+n > d.w {
		err := d.read0()
		if err != nil {
			return 0, err
		}
	}
	d.r += n
	return n, nil
}

func (d *decoder) InputOffset() int {
	return d.r
}

func (d *decoder) Reset(r io.Reader) {
	d.rd = r
	d.r = 0
	d.w = 0
	d.attrs = d.attrs[:0]
	d.bb = d.bb[:0]
	d.top = 0
	d.lastStartElement = false
}

func (d *decoder) skipWhitespacesGeneric(b byte) (byte, error) {
	for {
		if !isWhitespace(b) {
			return b, nil
		}
		var err error
		b, err = d.readByte()
		if err != nil {
			return 0, err
		}
	}
}

func (d *decoder) NextToken(t *Token) error {
	for {
		// read next character
		b, err := d.readByte()
		if err != nil {
			return err
		}
		switch b {
		case '>':
			// Previous StartElement now got properly ended.
			// That's fine. We just did not consume the end token
			// because there could have been an implicit
			// "/>" close at the end of the start element.
			d.lastStartElement = false
		case '/':
			if d.lastStartElement {
				// Immediately closing last openend StartElement.
				// This will generate an EndElement with the same
				// name that we used in the previous StartElement.
				_, err = d.discard(1)
				if err != nil {
					return err
				}
				d.lastStartElement = false
				return d.decodeEndElement(t, d.lastOpen)
			}
			d.unreadByte()
			cntn, err := d.decodeText(t)
			if err != nil || !cntn {
				return err
			}
		case '<':
			b, err = d.readByte()
			if err != nil {
				return err
			}
			switch b {
			case '?':
				d.lastStartElement = false
				err = d.decodeProcInst(t)
				d.unreadByte()
				return err
			case '!':
				// CDATA or comment
				b, err = d.readByte()
				if err != nil {
					return err
				}
				switch b {
				case '-':
					err = d.ignoreComment()
					if err != nil {
						return err
					}
				case '[':
					d.lastStartElement = false
					return d.readCDATA()
				default:
					return errors.New("invalid XML: comment or CDATA expected")
				}
			case '/':
				var name Name
				name, _, err = d.readName()
				if err != nil {
					return err
				}
				d.lastStartElement = false
				return d.decodeEndElement(t, name)
			default:
				d.lastStartElement = true
				return d.decodeStartElement(t)
			}
		default:
			d.lastStartElement = false
			d.unreadByte()
			cntn, err := d.decodeText(t)
			if err != nil || !cntn {
				return err
			}
		}
	}
}

func (d *decoder) decodeProcInst(t *Token) error {
	name, b, err := d.readName()
	if err != nil {
		return err
	}
	b, err = d.skipWhitespaces(b)
	if err != nil {
		return err
	}
	i := len(d.bb)
	j := i
	for {
		if b == '?' {
			for {
				var b2 byte
				b2, err = d.readByte()
				if err != nil {
					return err
				}
				if b2 == '>' {
					t.Kind = TokenTypeProcInst
					t.Name = name
					t.ByteData = d.bb[i:j]
					return nil
				} else if b2 != '?' {
					d.bb = append(d.bb, b, b2)
					if !isWhitespace(b2) {
						j = len(d.bb)
					}
					break
				}
				d.bb = append(d.bb, b2)
				if !isWhitespace(b2) {
					j = len(d.bb)
				}
			}
		} else {
			d.bb = append(d.bb, b)
			if !isWhitespace(b) {
				j = len(d.bb)
			}
		}
		b, err = d.readByte()
		if err != nil {
			return err
		}
	}
}

func (d *decoder) ignoreComment() error {
	_, err := d.discard(1)
	if err != nil {
		return err
	}
	for {
		for d.w > d.r {
			k := bytes.IndexByte(d.rb[d.r:d.w], '-')
			if k > -1 {
				_, err = d.discard(k + 1)
				if err != nil {
					return err
				}
				var b2 byte
				b2, err = d.readByte()
				if err != nil {
					return err
				}
				if b2 == '-' {
					for {
						var b3 byte
						b3, err = d.readByte()
						if err != nil {
							return err
						}
						if b3 == '>' {
							return nil
						} else if b3 != '-' {
							break
						}
					}
				}
			} else {
				d.discardBuffer()
			}
		}
		err := d.read0()
		if err != nil {
			return err
		}
	}
}

func (d *decoder) decodeEndElement(t *Token, name Name) error {
	end := len(d.attrs) - int(d.numAttributes[d.top])
	d.attrs = d.attrs[0:end]
	d.bb = d.bb[:d.bbOffset[d.top]]
	t.Kind = TokenTypeEndElement
	t.Name = name
	d.top--
	return nil
}

func (d *decoder) decodeStartElement(t *Token) error {
	d.top++
	d.numAttributes[d.top] = 0
	d.bbOffset[d.top] = int32(len(d.bb))
	d.preserveWhitespaces[d.top+1] = d.preserveWhitespaces[d.top]
	d.unreadByte()
	name, b, err := d.readName()
	if err != nil {
		return err
	}
	var attributes []Attr
	attributes, err = d.decodeAttributes(b)
	if err != nil {
		return err
	}
	d.lastOpen = name
	t.Kind = TokenTypeStartElement
	t.Name = name
	t.Attr = attributes
	d.unreadByte()
	return nil
}

func (d *decoder) decodeTextGeneric(t *Token) (bool, error) {
	i := len(d.bb)
	onlyWhitespaces := true
	for {
		j := d.r
		for k := j; k < d.w; k++ {
			b := d.rb[k]
			if b == '<' {
				_, err := d.discard(k - j)
				if err != nil {
					return false, err
				}
				if onlyWhitespaces && !d.preserveWhitespaces[d.top] {
					return true, nil
				}
				d.bb = append(d.bb, d.rb[j:k]...)
				t.Kind = TokenTypeTextElement
				t.ByteData = d.bb[i:len(d.bb)]
				return false, nil
			}
			onlyWhitespaces = onlyWhitespaces && isWhitespace(b)
		}
		d.bb = append(d.bb, d.rb[j:d.w]...)
		d.discardBuffer()
		err := d.read0()
		if err != nil {
			return false, err
		}
	}
}

func (d *decoder) readCDATA() error {
	// discard "CDATA["
	_, err := d.discard(6)
	if err != nil {
		return err
	}
	return errors.New("NYI")
}

func (d *decoder) readName() (Name, byte, error) {
	localOrPrefix, b, err := d.readSimpleName()
	if err != nil {
		return Name{}, 0, err
	}
	if b == ':' {
		var local []byte
		local, b, err = d.readSimpleName()
		if err != nil {
			return Name{}, 0, err
		}
		return Name{
			Local:  local,
			Prefix: localOrPrefix,
		}, b, nil
	}
	return Name{
		Local: localOrPrefix,
	}, b, nil
}

var seps = generateTable()

func generateTable() ['>' + 1]bool {
	var s ['>' + 1]bool
	s['\t'] = true
	s['\n'] = true
	s['\r'] = true
	s[' '] = true
	s['/'] = true
	s[':'] = true
	s['='] = true
	s['>'] = true
	return s
}

func isSeparator(b byte) bool {
	return int(b) < len(seps) && seps[b]
}

func (d *decoder) readSimpleNameGeneric() ([]byte, byte, error) {
	i := len(d.bb)
	for {
		j := d.r
		for k := j; k < d.w; k++ {
			if isSeparator(d.rb[k]) {
				d.bb = append(d.bb, d.rb[j:k]...)
				_, err := d.discard(k - j + 1)
				if err != nil {
					return nil, 0, err
				}
				return d.bb[i:len(d.bb)], d.rb[k], nil
			}
		}
		d.bb = append(d.bb, d.rb[j:d.w]...)
		d.discardBuffer()
		err := d.read0()
		if err != nil {
			return nil, 0, err
		}
	}
}

func (d *decoder) decodeAttributes(b byte) ([]Attr, error) {
	i := len(d.attrs)
	for {
		var err error
		b, err = d.skipWhitespaces(b)
		if err != nil {
			return nil, err
		}
		switch b {
		case '/', '>':
			return d.attrs[i:len(d.attrs)], nil
		default:
			i := len(d.attrs)
			d.attrs = d.attrs[:i+1]
			err = d.decodeAttribute(&d.attrs[i])
			if err != nil {
				return nil, err
			}
			b, err = d.readByte()
			if err != nil {
				return nil, err
			}
			d.numAttributes[d.top]++
		}
	}
}

// decodeAttribute parses a single XML attribute.
// After this function returns, the next reader symbol
// is the byte after the closing single or double quote
// of the attribute's value.
func (d *decoder) decodeAttribute(attr *Attr) error {
	d.unreadByte()
	name, b, err := d.readName()
	if err != nil {
		return err
	}
	b, err = d.skipWhitespaces(b)
	if err != nil {
		return err
	}
	if b != '=' {
		return fmt.Errorf("expected '=' character following attribute %+v", name)
	}
	b, err = d.readByte()
	if err != nil {
		return err
	}
	b, err = d.skipWhitespaces(b)
	if err != nil {
		return err
	}
	value, singleQuote, err := d.readString(b)
	if err != nil {
		return err
	}
	// xml:space?
	if bytes.Equal(name.Prefix, bsxml) && bytes.Equal(name.Local, bsspace) {
		d.preserveWhitespaces[d.top] = bytes.Equal(value, bspreserve)
	}
	attr.Name = name
	attr.SingleQuote = singleQuote
	attr.Value = value
	return nil
}

// readString parses a single string (in single or double quotes)
func (d *decoder) readString(b byte) ([]byte, bool, error) {
	i := len(d.bb)
	singleQuote := b == '\''
	for {
		j := d.r
		k := bytes.IndexByte(d.rb[j:d.w], b)
		if k > -1 {
			d.bb = append(d.bb, d.rb[j:j+k]...)
			_, err := d.discard(k + 1)
			if err != nil {
				return nil, false, err
			}
			return d.bb[i:len(d.bb)], singleQuote, nil
		}
		d.bb = append(d.bb, d.rb[j:d.w]...)
		d.discardBuffer()
		err := d.read0()
		if err != nil {
			return nil, false, err
		}
	}
}
