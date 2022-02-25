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
	"bytes"
	"io"
	"sync"

	"github.com/golang/protobuf/proto"
)

const xmlNamespace = "xmlns"

var bufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

type element struct {
	pb *PBElement
}

func (e *element) AllAttributes() []Attribute {
	attributes := make([]Attribute, len(e.pb.GetAttributes()))
	for i, pbAttr := range e.pb.GetAttributes() {
		attributes[i] = Attribute{Label: pbAttr.Label, Value: pbAttr.Value}
	}
	return attributes
}

func (e *element) AttributeCount() int {
	return len(e.pb.GetAttributes())
}

func (e *element) Attribute(label string) string {
	for _, pbAttr := range e.pb.GetAttributes() {
		if pbAttr.Label == label {
			return pbAttr.Value
		}
	}
	return ""
}

func (e *element) AllChildren() []Element {
	elements := make([]Element, len(e.pb.GetElements()))
	for i, pbElement := range e.pb.GetElements() {
		elements[i] = &element{pb: pbElement}
	}
	return elements
}

func (e *element) ChildrenCount() int {
	return len(e.pb.GetElements())
}

func (e *element) Child(name string) Element {
	for _, pbElement := range e.pb.GetElements() {
		if pbElement.Name == name {
			return &element{pb: pbElement}
		}
	}
	return nil
}

func (e *element) Children(name string) []Element {
	var elements []Element
	for _, pbElement := range e.pb.GetElements() {
		if pbElement.Name == name {
			elements = append(elements, &element{pb: pbElement})
		}
	}
	return elements
}

func (e *element) ChildNamespace(name, ns string) Element {
	for _, pbElement := range e.pb.GetElements() {
		if pbElement.Name == name && getProtoElementAttribute(pbElement, xmlNamespace) == ns {
			return &element{pb: pbElement}
		}
	}
	return nil
}

func (e *element) ChildrenNamespace(name, ns string) []Element {
	var elements []Element
	for _, pbElement := range e.pb.GetElements() {
		if pbElement.Name == name && getProtoElementAttribute(pbElement, xmlNamespace) == ns {
			elements = append(elements, &element{pb: pbElement})
		}
	}
	return elements
}

func (e *element) ToXML(w io.Writer, includeClosing bool) error {
	if _, err := io.WriteString(w, "<"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, e.Name()); err != nil {
		return err
	}

	// serialize attributes
	for _, attr := range e.AllAttributes() {
		if len(attr.Value) == 0 {
			continue
		}
		if _, err := io.WriteString(w, " "); err != nil {
			return err
		}
		if _, err := io.WriteString(w, attr.Label); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "='"); err != nil {
			return err
		}
		if _, err := io.WriteString(w, attr.Value); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "'"); err != nil {
			return err
		}
	}

	// serialize elements
	if e.ChildrenCount() > 0 || len(e.Text()) > 0 {
		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}
		if len(e.Text()) > 0 {
			if err := escapeText(w, []byte(e.Text()), false); err != nil {
				return err
			}
		}
		for _, elem := range e.AllChildren() {
			if err := elem.ToXML(w, true); err != nil {
				return err
			}
		}

		if includeClosing {
			if _, err := io.WriteString(w, "</"); err != nil {
				return err
			}
			if _, err := io.WriteString(w, e.Name()); err != nil {
				return err
			}
			if _, err := io.WriteString(w, ">"); err != nil {
				return err
			}
		}
		return nil
	}
	if includeClosing {
		_, err := io.WriteString(w, "/>")
		return err
	}
	_, err := io.WriteString(w, ">")
	return err
}

func (e *element) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e.pb)
}

func (e *element) UnmarshalBinary(data []byte) error {
	return proto.Unmarshal(data, e.pb)
}

func (e *element) String() string {
	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	_ = e.ToXML(buf, true)
	return buf.String()
}

func (e *element) GoString() string {
	return e.String()
}

func (e *element) Name() string {
	return e.pb.GetName()
}

func (e *element) Text() string {
	return e.pb.GetText()
}

func (e *element) Proto() *PBElement {
	return e.pb
}
