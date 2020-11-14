package streamerror

import (
	"testing"

	"github.com/jackal-xmpp/stravaganza"

	"github.com/stretchr/testify/require"
)

func TestStreamError_Name(t *testing.T) {
	// given
	seInvalidXML := E(InvalidXML)
	seInvalidNamespace := E(InvalidNamespace)
	seHostUnknown := E(HostUnknown)
	seInvalidFrom := E(InvalidFrom)
	sePolicyViolation := E(PolicyViolation)
	seRemoteConnectionFailed := E(RemoteConnectionFailed)
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
