/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIQ_IsStanza(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", GetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ(false)

	require.True(t, IsStanza(iq))
}

func TestIQ_IsGet(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", GetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ(false)

	require.True(t, iq.IsGet())
}

func TestIQ_IsSet(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ(false)

	require.True(t, iq.IsSet())
}

func TestIQ_IsResult(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", ResultType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ(false)

	require.True(t, iq.IsResult())
}

func TestIQ_ResultIQ(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ(false)

	resultIQ := iq.ResultIQ()

	require.Equal(t, iq.ID(), resultIQ.ID())
	require.Equal(t, ResultType, resultIQ.Type())
	require.Equal(t, iq.ToJID().String(), resultIQ.FromJID().String())
	require.Equal(t, iq.FromJID().String(), resultIQ.ToJID().String())
}
