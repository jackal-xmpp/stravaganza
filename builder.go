/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
)

// Builder builds generic XML node elements.
type Builder struct {
	attrs      []*PBAttribute
	elements   []*PBElement
	name, text string
}

// NewBuilder returns a name initialized builder instance.
func NewBuilder(name string) *Builder {
	return &Builder{name: name}
}

// NewMessageBuilder returns a 'message' stanza builder instance.
func NewMessageBuilder() *Builder {
	return NewBuilder("message")
}

// NewPresenceBuilder returns a 'presence' stanza builder instance.
func NewPresenceBuilder() *Builder {
	return NewBuilder("presence")
}

// NewPresenceBuilder returns an 'iq' stanza builder instance.
func NewIQBuilder() *Builder {
	return NewBuilder("iq")
}

// NewBuilderFromElement returns an element builder derived from a copied element.
func NewBuilderFromElement(element Element) *Builder {
	if element == nil {
		return &Builder{}
	}
	protoFrom := element.Proto()
	return &Builder{
		name:     protoFrom.GetName(),
		text:     protoFrom.GetText(),
		attrs:    copyProtoAttributes(protoFrom.GetAttributes()),
		elements: protoFrom.GetElements(),
	}
}

// NewBuilderFromBinary returns an element builder derived from an element binary representation.
func NewBuilderFromBinary(b []byte) (*Builder, error) {
	var protoFrom PBElement
	if err := proto.Unmarshal(b, &protoFrom); err != nil {
		return nil, err
	}
	return &Builder{
		name:     protoFrom.GetName(),
		text:     protoFrom.GetText(),
		attrs:    protoFrom.GetAttributes(),
		elements: protoFrom.GetElements(),
	}, nil
}

// NewBuilderFromProto returns an element builder derived from proto type.
func NewBuilderFromProto(protoFrom *PBElement) *Builder {
	return &Builder{
		name:     protoFrom.GetName(),
		text:     protoFrom.GetText(),
		attrs:    protoFrom.GetAttributes(),
		elements: protoFrom.GetElements(),
	}
}

// WithName sets XML node name.
func (b *Builder) WithName(name string) *Builder {
	b.name = name
	return b
}

// WithAttribute sets an XML node attribute (label=value).
func (b *Builder) WithAttribute(label, value string) *Builder {
	for i, pbAttr := range b.attrs {
		if pbAttr.Label == label {
			b.attrs[i].Value = value
			return b
		}
	}
	b.attrs = append(b.attrs, &PBAttribute{Label: label, Value: value})
	return b
}

// WithAttributes sets all XML node attributes.
func (b *Builder) WithAttributes(attributes ...Attribute) *Builder {
	for _, attr := range attributes {
		b.WithAttribute(attr.Label, attr.Value)
	}
	return b
}

// WithoutAttribute removes an XML node attribute.
func (b *Builder) WithoutAttribute(label string) *Builder {
	for i, pbAttr := range b.attrs {
		if pbAttr.Label == label {
			b.attrs = append(b.attrs[:i], b.attrs[i+1:]...)
			return b
		}
	}
	return b
}

// WithChild appends a new sub element.
func (b *Builder) WithChild(child Element) *Builder {
	b.elements = append(b.elements, child.Proto())
	return b
}

// WithChildren appends all new sub elements.
func (b *Builder) WithChildren(children ...Element) *Builder {
	for _, child := range children {
		b.elements = append(b.elements, child.Proto())
	}
	return b
}

// WithoutChildren removes all elements with a given name.
func (b *Builder) WithoutChildren(name string) *Builder {
	filtered := b.elements[:0]
	for _, pbElem := range b.elements {
		if pbElem.Name != name {
			filtered = append(filtered, pbElem)
		}
	}
	b.elements = filtered
	return b
}

// WithoutChildrenNamespace removes all elements with a given name and namespace.
func (b *Builder) WithoutChildrenNamespace(name, ns string) *Builder {
	filtered := b.elements[:0]
	for _, pbElem := range b.elements {
		if pbElem.Name != name && getProtoElementAttribute(pbElem, xmlNamespace) == ns {
			filtered = append(filtered, pbElem)
		}
	}
	b.elements = filtered
	return b
}

// WithText sets XML node text value.
func (b *Builder) WithText(text string) *Builder {
	b.text = text
	return b
}

// Build returns a new element instance.
func (b *Builder) Build() Element {
	return &element{pb: b.buildProtoElement()}
}

// BuildStanza validates and returns a generic stanza instance.
func (b *Builder) BuildStanza(validateJIDs bool) (Stanza, error) {
	pbElem := b.buildProtoElement()
	s := &stanza{
		element: element{pb: pbElem},
	}
	// validate 'to' and 'from' JIDs...
	if err := s.setFromAndToJIDs(validateJIDs); err != nil {
		return nil, err
	}
	return s, nil
}

// BuildIQ validates and returns a new IQ stanza.
func (b *Builder) BuildIQ(validateJIDs bool) (*IQ, error) {
	if b.name != "iq" {
		return nil, fmt.Errorf("stravaganza: wrong iq element name: %s", b.name)
	}
	pbElem := b.buildProtoElement()
	if len(getProtoElementAttribute(pbElem, "id")) == 0 {
		return nil, errors.New(`stravaganza: iq "id" attribute is required`)
	}
	iqType := getProtoElementAttribute(pbElem, "type")
	if len(iqType) == 0 {
		return nil, errors.New(`stravaganza: iq "type" attribute is required`)
	}
	if !isIQType(iqType) {
		return nil, fmt.Errorf(`stravaganza: invalid iq "type" attribute: %s`, iqType)
	}
	elemCount := len(pbElem.Elements)
	if (iqType == GetType || iqType == SetType) && elemCount != 1 {
		return nil, errors.New(`stravaganza: an iq stanza of type "get" or "set" must contain one and only one child element`)
	}
	if iqType == ResultType && elemCount > 1 {
		return nil, errors.New(`stravaganza: an iq stanza of type "result" must include zero or one child elements`)
	}
	iq := &IQ{
		stanza: stanza{
			element: element{pb: pbElem},
		},
	}
	// validate 'to' and 'from' JIDs...
	if err := iq.setFromAndToJIDs(validateJIDs); err != nil {
		return nil, err
	}
	return iq, nil
}

// BuildMessage validates and returns a new Message stanza.
func (b *Builder) BuildMessage(validateJIDs bool) (*Message, error) {
	if b.name != "message" {
		return nil, fmt.Errorf(`stravaganza: wrong message element name: %s`, b.name)
	}
	pbElem := b.buildProtoElement()
	messageType := getProtoElementAttribute(pbElem, "type")
	if !isMessageType(messageType) {
		return nil, fmt.Errorf(`stravaganza: invalid message "type" attribute: %s`, messageType)
	}
	m := &Message{
		stanza: stanza{
			element: element{pb: pbElem},
		},
	}
	// validate 'to' and 'from' JIDs...
	if err := m.setFromAndToJIDs(validateJIDs); err != nil {
		return nil, err
	}
	return m, nil
}

// BuildPresence validates and returns a new Presence stanza.
func (b *Builder) BuildPresence(validateJIDs bool) (*Presence, error) {
	if b.name != "presence" {
		return nil, fmt.Errorf("stravaganza: wrong presence element name: %s", b.name)
	}
	pbElem := b.buildProtoElement()
	presenceType := getProtoElementAttribute(pbElem, "type")
	if !isPresenceType(presenceType) {
		return nil, fmt.Errorf(`stravaganza: invalid presence "type" attribute: %s`, presenceType)
	}
	p := &Presence{
		stanza: stanza{
			element: element{pb: pbElem},
		},
	}
	// validate 'to' and 'from' JIDs...
	if err := p.setFromAndToJIDs(validateJIDs); err != nil {
		return nil, err
	}
	// validate presence status...
	if err := p.validateStatus(); err != nil {
		return nil, err
	}
	// set show and priority values...
	if err := p.setShow(); err != nil {
		return nil, err
	}
	if err := p.setPriority(); err != nil {
		return nil, err
	}
	return p, nil
}

func (b *Builder) buildProtoElement() *PBElement {
	return &PBElement{
		Name:       b.name,
		Attributes: b.attrs,
		Elements:   b.elements,
		Text:       b.text,
	}
}

func copyProtoAttributes(pbAttrs []*PBAttribute) []*PBAttribute {
	cp := make([]*PBAttribute, len(pbAttrs))
	for i, pbAttr := range pbAttrs {
		var cpAttr PBAttribute
		cpAttr.Label = pbAttr.Label
		cpAttr.Value = pbAttr.Value
		cp[i] = &cpAttr
	}
	return cp
}

func isIQType(tp string) bool {
	switch tp {
	case ErrorType, GetType, SetType, ResultType:
		return true
	}
	return false
}

func isMessageType(messageType string) bool {
	switch messageType {
	case "", ErrorType, NormalType, HeadlineType, ChatType, GroupChatType:
		return true
	default:
		return false
	}
}

func isPresenceType(presenceType string) bool {
	switch presenceType {
	case ErrorType, AvailableType, UnavailableType, SubscribeType,
		UnsubscribeType, SubscribedType, UnsubscribedType, ProbeType:
		return true
	default:
		return false
	}
}
