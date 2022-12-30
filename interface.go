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
	"encoding"
	"fmt"
	"io"

	"github.com/jackal-xmpp/stravaganza/jid"
)

// Attribute represents an XML node attribute (label=value).
type Attribute struct {
	Label string
	Value string
}

// AttributeReader defines an XML attributes read-only interface.
type AttributeReader interface {
	// AllAttributes returns a list of all node attributes.
	AllAttributes() []Attribute

	// AttributeCount returns node total attribute count.
	AttributeCount() int

	// Attribute returns XML node attribute value.
	Attribute(label string) string
}

// ElementReader defines an XML sub elements read-only interface.
type ElementReader interface {
	// AllChildren returns a list of all child nodes.
	AllChildren() []Element

	// ChildrenCount returns child elements count.
	ChildrenCount() int

	// Child returns first element identified by name.
	Child(name string) Element

	// Children returns all elements identified by name.
	// Returns an empty array if no elements are found.
	Children(name string) []Element

	// ChildNamespace returns first element identified by name and namespace.
	// Returns nil if no element is found.
	ChildNamespace(name, ns string) Element

	// ChildrenNamespace returns all elements identified by name and namespace.
	ChildrenNamespace(name, ns string) []Element
}

// XMLSerializer represents element common XML serializer interface.
type XMLSerializer interface {
	// ToXML serializes element to a raw XML representation.
	// includeClosing determines if closing tag should be attached.
	ToXML(w io.Writer, includeClosing bool) error
}

// Element represents a generic XML node element.
type Element interface {
	AttributeReader
	ElementReader
	XMLSerializer
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	fmt.Stringer
	fmt.GoStringer

	// Name returns XML node name.
	Name() string

	// Text returns XML node text value.
	Text() string

	// Proto returns element protobuf message.
	Proto() *PBElement
}

// Stanza represents an XMPP stanza element.
type Stanza interface {
	Element

	// ToJID returns stanza 'to' JID value.
	ToJID() *jid.JID

	// FromJID returns stanza 'from' JID value.
	FromJID() *jid.JID

	// ID returns 'id' node attribute.
	ID() string

	// Namespace returns 'xmlns' node attribute.
	Namespace() string

	// Type returns 'type' node attribute.
	Type() string

	// IsError returns true if stanza has a 'type' attribute of value 'error'.
	IsError() bool

	// Error returns stanza error sub element.
	Error() Element
}
