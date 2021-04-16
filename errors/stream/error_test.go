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

package streamerror

import (
	"testing"

	"github.com/jackal-xmpp/stravaganza/v2"
	"github.com/stretchr/testify/require"
)

func TestStreamError_Error(t *testing.T) {
	// given
	var seInvalidXML interface{} = E(InvalidXML)

	// when
	err, ok := seInvalidXML.(error)

	// then
	require.NotNil(t, err)
	require.True(t, ok)
}

func TestStreamError_Name(t *testing.T) {
	// given
	seInvalidXML := E(InvalidXML)
	seInvalidNamespace := E(InvalidNamespace)
	seHostUnknown := E(HostUnknown)
	seInvalidFrom := E(InvalidFrom)
	sePolicyViolation := E(PolicyViolation)
	seRemoteConnectionFailed := E(RemoteConnectionFailed)
	seConflict := E(Conflict)
	seConnectionTimeout := E(ConnectionTimeout)
	seUnsupportedStanzaType := E(UnsupportedStanzaType)
	seUnsupportedVersion := E(UnsupportedVersion)
	seNotAuthorized := E(NotAuthorized)
	seResourceConstraint := E(ResourceConstraint)
	seSystemShutdown := E(SystemShutdown)
	seUndefinedCondition := E(UndefinedCondition)
	seInternalServerError := E(InternalServerError)

	// then
	require.Equal(t, "invalid-xml", seInvalidXML.Error())
	require.Equal(t, "invalid-namespace", seInvalidNamespace.Error())
	require.Equal(t, "host-unknown", seHostUnknown.Error())
	require.Equal(t, "invalid-from", seInvalidFrom.Error())
	require.Equal(t, "policy-violation", sePolicyViolation.Error())
	require.Equal(t, "remote-connection-failed", seRemoteConnectionFailed.Error())
	require.Equal(t, "conflict", seConflict.Error())
	require.Equal(t, "connection-timeout", seConnectionTimeout.Error())
	require.Equal(t, "unsupported-stanza-type", seUnsupportedStanzaType.Error())
	require.Equal(t, "unsupported-version", seUnsupportedVersion.Error())
	require.Equal(t, "not-authorized", seNotAuthorized.Error())
	require.Equal(t, "resource-constraint", seResourceConstraint.Error())
	require.Equal(t, "system-shutdown", seSystemShutdown.Error())
	require.Equal(t, "undefined-condition", seUndefinedCondition.Error())
	require.Equal(t, "internal-server-error", seInternalServerError.Error())
}

func TestStreamError_Element(t *testing.T) {
	// given
	se := E(PolicyViolation)
	se.Lang = "es"
	se.Text = "Límite de conexiones alcanzado"
	se.ApplicationElement = stravaganza.NewBuilder("connection-limit-reached").Build()

	// when
	el := se.Element()
	errEl := el.Child("policy-violation")
	appEl := el.Child("connection-limit-reached")

	// then
	require.NotNil(t, errEl)
	require.NotNil(t, appEl)

	require.Equal(t, "stream:error", el.Name())
	require.Equal(t, "policy-violation", errEl.Name())
	require.Equal(t, "connection-limit-reached", appEl.Name())
}
