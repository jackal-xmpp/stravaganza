package gosaxml

import "github.com/klauspost/cpuid/v2"

var canUseSSE = cpuid.CPU.Has(cpuid.SSE2) && cpuid.CPU.Has(cpuid.BMI1)
var canUseAVX2 = canUseSSE && cpuid.CPU.Has(cpuid.AVX2)

func init() {
	if canUseAVX2 {
		simdWidth = 32
	} else if canUseSSE {
		simdWidth = 16
	}
}

func (d *decoder) skipWhitespaces(b byte) (byte, error) {
	if canUseAVX2 {
		return d.skipWhitespacesAVX2(b)
	} else if canUseSSE {
		return d.skipWhitespacesSSE(b)
	}
	return d.skipWhitespacesGeneric(b)
}

func (d *decoder) skipWhitespacesAVX2(b byte) (byte, error) {
	if !isWhitespace(b) {
		return b, nil
	}
	for {
		for d.w > d.r {
			sidx := int(onlySpaces32(d.rb[d.r:d.w]))
			_, err := d.discard(sidx)
			if err != nil {
				return 0, err
			}
			if sidx != 32 {
				newB, err := d.readByte()
				if err != nil {
					return 0, err
				}
				return newB, nil
			}
		}
		d.discardBuffer()
		err := d.read0()
		if err != nil {
			return 0, err
		}
	}
}

func (d *decoder) skipWhitespacesSSE(b byte) (byte, error) {
	if !isWhitespace(b) {
		return b, nil
	}
	for {
		j := d.r
		c := 0
		for d.w > d.r+c {
			sidx := onlySpaces32(d.rb[j+c : d.w])
			c += int(sidx)
			if sidx != 16 {
				_, err := d.discard(c)
				if err != nil {
					return 0, err
				}
				newB, err := d.readByte()
				if err != nil {
					return 0, err
				}
				return newB, nil
			}
		}
		d.discardBuffer()
		err := d.read0()
		if err != nil {
			return 0, err
		}
	}
}

func (d *decoder) decodeText(t *Token) (bool, error) {
	if canUseAVX2 {
		return d.decodeTextAVX2(t)
	} else if canUseSSE {
		return d.decodeTextSSE(t)
	}
	return d.decodeTextGeneric(t)
}

func (d *decoder) decodeTextSSE(t *Token) (bool, error) {
	i := len(d.bb)
	onlyWhitespaces := true
	for {
		j := d.r
		c := 0
		for d.w > d.r+c {
			sidx := openAngleBracket16(d.rb[j+c : d.w])
			onlyWhitespaces = onlyWhitespaces && onlySpaces16(d.rb[j+c:d.w]) >= sidx
			c += int(sidx)
			if sidx != 16 {
				_, err := d.discard(c)
				if err != nil {
					return false, err
				}
				if onlyWhitespaces && !d.preserveWhitespaces[d.top] {
					return true, nil
				}
				d.bb = append(d.bb, d.rb[j:j+c]...)
				t.Kind = TokenTypeTextElement
				t.ByteData = d.bb[i:len(d.bb)]
				return false, nil
			}
		}
		d.bb = append(d.bb, d.rb[j:d.w]...)
		d.discardBuffer()
		err := d.read0()
		if err != nil {
			return false, err
		}
	}
}

func (d *decoder) decodeTextAVX2(t *Token) (bool, error) {
	i := len(d.bb)
	onlyWhitespaces := true
	for {
		j := d.r
		c := 0
		for d.w > d.r+c {
			sidx := openAngleBracket32(d.rb[j+c : d.w])
			onlyWhitespaces = onlyWhitespaces && onlySpaces32(d.rb[j+c:d.w]) >= sidx
			c += int(sidx)
			if sidx != 32 {
				_, err := d.discard(c)
				if err != nil {
					return false, err
				}
				if onlyWhitespaces && !d.preserveWhitespaces[d.top] {
					return true, nil
				}
				d.bb = append(d.bb, d.rb[j:j+c]...)
				t.Kind = TokenTypeTextElement
				t.ByteData = d.bb[i:len(d.bb)]
				return false, nil
			}
		}
		d.bb = append(d.bb, d.rb[j:d.w]...)
		d.discardBuffer()
		err := d.read0()
		if err != nil {
			return false, err
		}
	}
}

func (d *decoder) readSimpleName() ([]byte, byte, error) {
	if canUseAVX2 {
		return d.readSimpleNameAVX()
	}
	return d.readSimpleNameGeneric()
}

func (d *decoder) readSimpleNameAVX() ([]byte, byte, error) {
	i := len(d.bb)
	for {
		j := d.r
		c := 0
		for d.w > d.r+c {
			sidx := int(seperator32(d.rb[j+c : d.w]))
			c += sidx
			if sidx != 32 {
				_, err := d.discard(c + 1)
				if err != nil {
					return nil, 0, err
				}
				d.bb = append(d.bb, d.rb[j:j+c]...)
				return d.bb[i:len(d.bb)], d.rb[j+c], nil
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
