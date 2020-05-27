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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestElement_Name(t *testing.T) {
	el := &element{pb: &PBElement{Name: "n1"}}

	require.Equal(t, "n1", el.Name())
}

func TestElement_Text(t *testing.T) {
	el := &element{pb: &PBElement{Name: "n1", Text: "some text"}}

	require.Equal(t, "some text", el.Text())
}

func TestElement_Attribute(t *testing.T) {
	el := &element{pb: &PBElement{
		Name:       "n1",
		Attributes: []*PBAttribute{{Label: "l1", Value: "v1"}},
	}}
	require.Equal(t, "v1", el.Attribute("l1"))
}

func TestElement_AttributeCount(t *testing.T) {
	el := &element{pb: &PBElement{
		Name:       "n1",
		Attributes: []*PBAttribute{{Label: "l1", Value: "v1"}},
	}}
	require.Equal(t, 1, el.AttributeCount())
}

func TestElement_AllAttributes(t *testing.T) {
	el := &element{pb: &PBElement{
		Name:       "n1",
		Attributes: []*PBAttribute{{Label: "l1", Value: "v1"}},
	}}
	attrs := el.AllAttributes()
	require.Len(t, attrs, 1)
	require.Equal(t, "v1", attrs[0].Value)
}

func TestElement_Children(t *testing.T) {
	el := &element{pb: &PBElement{
		Name:     "n1",
		Elements: []*PBElement{{Name: "n2"}, {Name: "n2", Attributes: []*PBAttribute{{Label: xmlNamespace, Value: "ns2"}}}},
	}}

	children := el.AllChildren()
	require.Len(t, children, 2)

	require.Equal(t, 2, el.ChildrenCount())

	n2 := el.Child("n2")
	require.NotNil(t, n2)

	n2s := el.Children("n2")
	require.NotNil(t, n2)
	require.Len(t, n2s, 2)

	n2 = el.ChildNamespace("n2", "ns2")
	require.NotNil(t, n2)
	require.Equal(t, "ns2", n2.Attribute(xmlNamespace))

	n2s = el.ChildrenNamespace("n2", "ns2")
	require.NotNil(t, n2)
	require.Len(t, n2s, 1)
}

func TestElement_ToXML(t *testing.T) {
	el := &element{pb: &PBElement{
		Name:       "n1",
		Attributes: []*PBAttribute{{Label: "l1", Value: "v1"}},
		Elements:   []*PBElement{{Name: "n2", Attributes: []*PBAttribute{{Label: xmlNamespace, Value: "ns2"}}}},
		Text:       "Where arth thou, my Juliet?",
	}}

	buf1 := bytes.NewBuffer(nil)
	buf2 := bytes.NewBuffer(nil)

	_ = el.ToXML(buf1, true)
	_ = el.ToXML(buf2, false)

	str := el.String()
	goStr := el.GoString()

	require.Equal(t, "<n1 l1=\"v1\">Where arth thou, my Juliet?<n2 xmlns=\"ns2\"/></n1>", buf1.String())
	require.Equal(t, "<n1 l1=\"v1\">Where arth thou, my Juliet?<n2 xmlns=\"ns2\"/>", buf2.String())

	require.Equal(t, buf1.String(), str)
	require.Equal(t, str, goStr)
}

func TestElement_MarshalBinary(t *testing.T) {
	el := NewBuilderFromElement(nil).
		WithName("n1").
		WithAttribute("id", "id1234").
		WithAttribute(xmlNamespace, "com.stravaganza.ns").
		Build()

	protoBytes, _ := el.MarshalBinary()

	b, _ := NewBuilderFromBinary(protoBytes)
	el2 := b.Build()

	require.Equal(t, el.Proto().GetName(), el2.Proto().GetName())
	require.Equal(t, el.Proto().GetText(), el2.Proto().GetText())
	require.Equal(t, el.Proto().GetElements(), el2.Proto().GetElements())
	require.Equal(t, el.Proto().GetText(), el2.Proto().GetText())
}
