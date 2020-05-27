/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
