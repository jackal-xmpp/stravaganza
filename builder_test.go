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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilder_MessageBuilder(t *testing.T) {
	b := NewMessageBuilder()
	require.Equal(t, "message", b.name)
}

func TestBuilder_PresenceBuilder(t *testing.T) {
	b := NewPresenceBuilder()
	require.Equal(t, "presence", b.name)
}

func TestBuilder_IQBuilder(t *testing.T) {
	b := NewIQBuilder()
	require.Equal(t, "iq", b.name)
}

func TestBuilder_WithName(t *testing.T) {
	el := NewBuilderFromElement(nil).
		WithName("node").
		WithText("some text").
		WithAttribute("id", "1234").
		Build()
	require.Equal(t, "node", el.Name())
}

func TestBuilder_WithText(t *testing.T) {
	el := NewBuilderFromElement(nil).
		WithName("node").
		WithText("some text").
		WithAttribute("id", "1234").
		Build()
	require.Equal(t, "some text", el.Text())
}

func TestBuilder_WithAttribute(t *testing.T) {
	el := NewBuilderFromElement(nil).
		WithName("node").
		WithAttribute("id", "id1234").
		WithAttribute("id", "id5678").
		Build()

	attr := el.Attribute("id")
	require.NotNil(t, attr)
	require.Equal(t, "id5678", attr)
}

func TestBuilder_WithAttributes(t *testing.T) {
	el := NewBuilderFromElement(nil).
		WithName("node").
		WithAttributes(
			Attribute{Label: "id1", Value: "id1234"},
			Attribute{Label: "id2", Value: "id5678"},
		).
		Build()

	require.Equal(t, "id1234", el.Attribute("id1"))
	require.Equal(t, "id5678", el.Attribute("id2"))
}

func TestBuilder_WithoutAttribute(t *testing.T) {
	el1 := NewBuilderFromElement(nil).
		WithName("node").
		WithAttribute("id", "id1234").
		Build()
	el2 := NewBuilderFromElement(el1).
		WithoutAttribute("id").
		Build()

	attr := el2.Attribute("id")
	require.Len(t, attr, 0)
}

func TestBuilder_WithChildren(t *testing.T) {
	el1 := NewBuilderFromElement(nil).
		WithName("n1").
		WithAttribute("id", "id1234").
		Build()
	el2 := NewBuilderFromElement(nil).
		WithName("n2").
		WithChildren(el1).
		Build()

	child := el2.Child("n1")
	require.NotNil(t, child)
	require.Equal(t, "n1", child.Name())
}

func TestBuilder_WithoutChildren(t *testing.T) {
	el1 := NewBuilderFromElement(nil).
		WithName("n1").
		WithAttribute("id", "id1234").
		Build()
	el2 := NewBuilderFromElement(nil).
		WithName("n2").
		WithChildren(el1).
		Build()
	el3 := NewBuilderFromElement(el2).
		WithoutChildren("n1").
		Build()

	require.Nil(t, el3.Child("n1"))
}

func TestBuilder_WithoutChildrenNamespace(t *testing.T) {
	el0 := NewBuilderFromElement(nil).
		WithName("n2").
		WithChild(
			NewBuilder("n0").
				Build(),
		).
		WithChild(
			NewBuilder("n1").
				WithAttribute(xmlNamespace, "com.stravaganza.ns").
				Build(),
		).
		Build()
	el1 := NewBuilderFromElement(el0).
		WithoutChildrenNamespace("n1", "com.stravaganza.ns").
		Build()

	require.Nil(t, el1.Child("n1"))
	require.Equal(t, 1, el1.ChildrenCount())
}
