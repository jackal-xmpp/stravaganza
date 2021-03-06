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

func TestPresence_IsStanza(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		BuildPresence()

	require.True(t, IsStanza(p))
}

func TestPresence_IsAvailable(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		BuildPresence()

	require.True(t, p.IsAvailable())
}

func TestPresence_IsUnavailable(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", UnavailableType).
		BuildPresence()

	require.True(t, p.IsUnavailable())
}

func TestPresence_IsSubscribe(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SubscribeType).
		BuildPresence()

	require.True(t, p.IsSubscribe())
}

func TestPresence_IsUnsubscribe(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", UnsubscribeType).
		BuildPresence()

	require.True(t, p.IsUnsubscribe())
}

func TestPresence_IsSubscribed(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SubscribedType).
		BuildPresence()

	require.True(t, p.IsSubscribed())
}

func TestPresence_IsUnsubscribed(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", UnsubscribedType).
		BuildPresence()

	require.True(t, p.IsUnsubscribed())
}

func TestPresence_IsProbe(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", ProbeType).
		BuildPresence()

	require.True(t, p.IsProbe())
}

func TestPresence_ValidateValidStatus(t *testing.T) {
	p, err := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("status").WithAttribute("xml:lang", "en").WithText("Away").Build()).
		BuildPresence()

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
		BuildPresence()

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
		BuildPresence()

	require.Nil(t, p)
	require.NotNil(t, err)
}

func TestPresence_SetShow(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("show").WithText("chat").Build()).
		BuildPresence()

	err := p.setShow()

	require.Nil(t, err)
	require.Equal(t, ChatShowState, p.ShowState())
}

func TestPresence_SetInvalidShow(t *testing.T) {
	p, err := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("show").WithText("foo").Build()).
		BuildPresence()

	require.Nil(t, p)
	require.NotNil(t, err)
}

func TestPresence_SetPriority(t *testing.T) {
	p, _ := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("priority").WithText("126").Build()).
		BuildPresence()

	err := p.setPriority()

	require.Nil(t, err)
	require.Equal(t, int8(126), p.Priority())
}

func TestPresence_SetInvalidPriority(t *testing.T) {
	p, err := NewBuilder("presence").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithChild(NewBuilder("priority").WithText("300").Build()).
		BuildPresence()

	require.Nil(t, p)
	require.NotNil(t, err)
}
