/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza


const (
	// GetType represents a 'get' IQ type.
	GetType = "get"

	// SetType represents a 'set' IQ type.
	SetType = "set"

	// ResultType represents a 'result' IQ type.
	ResultType = "result"
)

// IQ type represents an <iq> element.
type IQ struct {
	stanza
}

// IsGet returns true if this is a 'get' type IQ.
func (iq *IQ) IsGet() bool {
	return iq.Type() == GetType
}

// IsSet returns true if this is a 'set' type IQ.
func (iq *IQ) IsSet() bool {
	return iq.Type() == SetType
}

// IsResult returns true if this is a 'result' type IQ.
func (iq *IQ) IsResult() bool {
	return iq.Type() == ResultType
}

// ResultIQ returns the instance associated result IQ.
func (iq *IQ) ResultIQ() *IQ {
	pbElem := &PBElement{
		Name: "iq",
		Attributes: []*PBAttribute{
			{Label: "xmlns", Value: iq.Namespace()},
			{Label: "id", Value: iq.ID()},
			{Label: "type", Value: ResultType},
			{Label: "from", Value: iq.ToJID().String()},
			{Label: "to", Value: iq.FromJID().String()},
		},
	}
	return &IQ{
		stanza: stanza{
			element: element{pb: pbElem},
			fromJID: iq.toJID,
			toJID:   iq.fromJID,
		},
	}
}
