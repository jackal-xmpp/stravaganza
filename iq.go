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

// IQName represents 'iq' stanza type name.
const IQName = "iq"

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

// ResultBuilder returns a builder instance associated to iq result stanza.
func (iq *IQ) ResultBuilder() *Builder {
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
	return NewBuilderFromProto(pbElem)
}
