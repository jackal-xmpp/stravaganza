/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

import (
	"errors"

	"github.com/jackal-xmpp/stravaganza/jid"
)

const (
	// ErrorType represents a generic 'error' type stanza.
	ErrorType = "error"
)

type stanza struct {
	element
	toJID, fromJID *jid.JID
}

func (s *stanza) ToJID() *jid.JID {
	return s.toJID
}

func (s *stanza) FromJID() *jid.JID {
	return s.fromJID
}

func (s *stanza) ID() string {
	return s.Attribute("id").Value
}

func (s *stanza) Namespace() string {
	return s.Attribute("xmlns").Value
}

func (s *stanza) Type() string {
	return s.Attribute("type").Value
}

func (s *stanza) IsError() bool {
	return s.Type() == "error"
}

func (s *stanza) Error() Element {
	return s.Child("error")
}

func (s *stanza) setFromAndToJIDs(validateJIDs bool) error {
	pbElem := s.element.pb
	fromAttr := getProtoElementAttribute(pbElem, "from")
	if len(fromAttr) == 0 {
		return errors.New("stravaganza: stanza 'from' attribute is required")
	}
	toAttr := getProtoElementAttribute(pbElem, "to")
	if len(toAttr) == 0 {
		return errors.New("stravaganza: stanza 'to' attribute is required")
	}
	fromJID, err := jid.NewWithString(fromAttr, !validateJIDs)
	if err != nil {
		return err
	}
	toJID, err := jid.NewWithString(toAttr, !validateJIDs)
	if err != nil {
		return err
	}
	s.fromJID = fromJID
	s.toJID = toJID
	return nil
}
