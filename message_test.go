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

func TestMessage_IsStanza(t *testing.T) {
	m, _ := NewBuilder("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", NormalType).
		BuildMessage(false)

	require.True(t, IsStanza(m))
}

func TestMessage_IsNormal(t *testing.T) {
	m1, _ := NewBuilder("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", NormalType).
		BuildMessage(false)

	m2, _ := NewBuilder("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		BuildMessage(false)

	require.True(t, m1.IsNormal())
	require.True(t, m2.IsNormal())
}

func TestMessage_IsHeadline(t *testing.T) {
	m, _ := NewBuilder("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", HeadlineType).
		BuildMessage(false)

	require.True(t, m.IsHeadline())
}

func TestMessage_IsChat(t *testing.T) {
	m, _ := NewBuilder("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", ChatType).
		BuildMessage(false)

	require.True(t, m.IsChat())
}

func TestMessage_IsGroupChat(t *testing.T) {
	m, _ := NewBuilder("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", GroupChatType).
		BuildMessage(false)

	require.True(t, m.IsGroupChat())
}

func TestMessage_IsMessageWithBody(t *testing.T) {
	m, _ := NewBuilder("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", NormalType).
		WithChild(NewBuilder("body").Build()).
		BuildMessage(false)

	require.True(t, m.IsMessageWithBody())
}
