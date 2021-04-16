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

func TestIQ_IsStanza(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", GetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ()

	require.True(t, IsStanza(iq))
}

func TestIQ_IsGet(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", GetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ()

	require.True(t, iq.IsGet())
}

func TestIQ_IsSet(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ()

	require.True(t, iq.IsSet())
}

func TestIQ_IsResult(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", ResultType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ()

	require.True(t, iq.IsResult())
}

func TestIQ_ResultIQ(t *testing.T) {
	iq, _ := NewBuilder("iq").
		WithAttribute("id", "1234").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "noelia@jackal.im/balcony").
		WithAttribute("type", SetType).
		WithChild(NewBuilder("q").Build()).
		BuildIQ()

	resultIQ, _ := iq.ResultBuilder().BuildIQ()

	require.Equal(t, iq.ID(), resultIQ.ID())
	require.Equal(t, ResultType, resultIQ.Type())
	require.Equal(t, iq.ToJID().String(), resultIQ.FromJID().String())
	require.Equal(t, iq.FromJID().String(), resultIQ.ToJID().String())
}
