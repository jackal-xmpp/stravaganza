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

func TestStanza_FromJID(t *testing.T) {
	s := newTestStanza(t)
	fromJID := s.FromJID()

	require.NotNil(t, fromJID)
	require.Equal(t, "ortuman@jackal.im/yard", fromJID.String())
}

func TestStanza_ToJID(t *testing.T) {
	s := newTestStanza(t)
	toJID := s.ToJID()

	require.NotNil(t, toJID)
	require.Equal(t, "noelia@jackal.im/balcony", toJID.String())
}

func TestStanza_ID(t *testing.T) {
	s := newTestStanza(t)
	require.Equal(t, "s1234", s.ID())
}

func TestStanza_Namespace(t *testing.T) {
	s := newTestStanza(t)
	require.Equal(t, "ns-1", s.Namespace())
}

func TestStanza_Type(t *testing.T) {
	s := newTestStanza(t)
	require.Equal(t, "error", s.Type())
}

func TestStanza_IsError(t *testing.T) {
	s := newTestStanza(t)
	require.True(t, s.IsError())
}

func TestStanza_SetFromAndToJIDs(t *testing.T) {
	s := newTestStanza(t)
	err := s.setFromAndToJIDs(true)

	require.Nil(t, err)
	require.NotNil(t, s.fromJID)
	require.NotNil(t, s.toJID)

	require.Equal(t, "ortuman@jackal.im/yard", s.FromJID().String())
	require.Equal(t, "noelia@jackal.im/balcony", s.ToJID().String())
}

func TestStanza_SetInvalidFromAndToJIDs(t *testing.T) {
	s, err := NewBuilder("iq").
		WithAttribute("xmlns", "ns-1").
		WithAttribute("id", "s1234").
		WithAttribute("type", "error").
		BuildStanza(true)

	require.Nil(t, s)
	require.NotNil(t, err)

	s, err = NewBuilder("iq").
		WithAttribute("xmlns", "ns-1").
		WithAttribute("id", "s1234").
		WithAttribute("type", "error").
		WithAttribute("from", "ortuman").
		BuildStanza(true)

	require.Nil(t, s)
	require.NotNil(t, err)

	s, err = NewBuilder("iq").
		WithAttribute("xmlns", "ns-1").
		WithAttribute("id", "s1234").
		WithAttribute("type", "error").
		WithAttribute("from", "ortuman@").
		WithAttribute("to", "noelia@").
		BuildStanza(true)

	require.Nil(t, s)
	require.NotNil(t, err)

	s, err = NewBuilder("iq").
		WithAttribute("xmlns", "ns-1").
		WithAttribute("id", "s1234").
		WithAttribute("type", "error").
		WithAttribute("from", "ortuman@jackal.im").
		WithAttribute("to", "noelia@").
		BuildStanza(true)

	require.Nil(t, s)
	require.NotNil(t, err)
}

func newTestStanza(t *testing.T) *stanza {
	t.Helper()

	s, _ := NewBuilder("iq").
		WithAttribute("xmlns", "ns-1").
		WithAttribute("id", "s1234").
		WithAttribute("type", "error").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("error").Build()).
		BuildStanza(true)
	return s.(*stanza)
}
