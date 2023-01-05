package gosaxml

func (d *decoder) skipWhitespaces(b byte) (byte, error) {
	return d.skipWhitespacesGeneric(b)
}

func (d *decoder) decodeText(t *Token) (bool, error) {
	return d.decodeTextGeneric(t)
}

func (d *decoder) readSimpleName() ([]byte, byte, error) {
	return d.readSimpleNameGeneric()
}
