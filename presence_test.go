/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPresence_IsAvailable(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		BuildPresence(false)

	require.True(t, p.IsAvailable())
}

func TestPresence_IsUnavailable(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", UnavailableType).
		BuildPresence(false)

	require.True(t, p.IsUnavailable())
}

func TestPresence_IsSubscribe(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SubscribeType).
		BuildPresence(false)

	require.True(t, p.IsSubscribe())
}

func TestPresence_IsUnsubscribe(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", UnsubscribeType).
		BuildPresence(false)

	require.True(t, p.IsUnsubscribe())
}

func TestPresence_IsSubscribed(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SubscribedType).
		BuildPresence(false)

	require.True(t, p.IsSubscribed())
}

func TestPresence_IsUnsubscribed(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", UnsubscribedType).
		BuildPresence(false)

	require.True(t, p.IsUnsubscribed())
}

func TestPresence_IsProbe(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", ProbeType).
		BuildPresence(false)

	require.True(t, p.IsProbe())
}

func TestPresence_ValidateValidStatus(t *testing.T) {
	p, err := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("status").WithAttribute("xml:lang", "en").WithText("Away").Build()).
		BuildPresence(false)

	require.NotNil(t, p)
	require.Nil(t, err)

	require.Equal(t, "Away", p.Status())
}

func TestPresence_Capabilities(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(
			NewBuilder("c").
				WithAttribute(xmlNamespace, capabilitiesNamespace).
				WithAttribute("node", "https://psi-plus.com").
				WithAttribute("ver", "16e5SrD2bLTkZGpg9MXPyvyCWEk=").
				WithAttribute("hash", "sha1").
				Build(),
		).
		BuildPresence(false)

	c := p.Capabilities()

	require.NotNil(t, c)

	require.Equal(t, "https://psi-plus.com", c.Node)
	require.Equal(t, "16e5SrD2bLTkZGpg9MXPyvyCWEk=", c.Ver)
	require.Equal(t, "sha1", c.Hash)
}

func TestPresence_ValidateInvalidStatus(t *testing.T) {
	p, err := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("status").WithAttribute("foo", "value").WithText("Away").Build()).
		BuildPresence(false)

	require.Nil(t, p)
	require.NotNil(t, err)
}

func TestPresence_SetShow(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("show").WithText("chat").Build()).
		BuildPresence(false)

	err := p.setShow()

	require.Nil(t, err)
	require.Equal(t, ChatShowState, p.ShowState())
}

func TestPresence_SetInvalidShow(t *testing.T) {
	p, err := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("show").WithText("foo").Build()).
		BuildPresence(false)

	require.Nil(t, p)
	require.NotNil(t, err)
}

func TestPresence_SetPriority(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("priority").WithText("126").Build()).
		BuildPresence(false)

	err := p.setPriority()

	require.Nil(t, err)
	require.Equal(t, int8(126), p.Priority())
}

func TestPresence_SetInvalidPriority(t *testing.T) {
	p, err := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("priority").WithText("300").Build()).
		BuildPresence(false)

	require.Nil(t, p)
	require.NotNil(t, err)
}
