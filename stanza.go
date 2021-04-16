// Copyright 2020 The jackal Authors
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

package stravaganza

import (
	"errors"

	"github.com/jackal-xmpp/stravaganza/v2/jid"
)

const (
	// ErrorType represents a generic 'error' type stanza.
	ErrorType = "error"
)

// IsStanza tells whether or not el element is of type stanza.
func IsStanza(el Element) bool {
	switch el.Name() {
	case IQName, PresenceName, MessageName:
		return true
	}
	return false
}

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
	return s.Attribute("id")
}

func (s *stanza) Namespace() string {
	return s.Attribute("xmlns")
}

func (s *stanza) Type() string {
	return s.Attribute("type")
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
