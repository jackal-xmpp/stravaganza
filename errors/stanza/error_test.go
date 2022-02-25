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

package stanzaerror

import (
	"testing"

	"github.com/jackal-xmpp/stravaganza"
	"github.com/stretchr/testify/require"
)

func TestStanzaError_Error(t *testing.T) {
	// given
	msg := testMessageStanza()

	var seBadRequest interface{} = E(BadRequest, msg)

	// when
	err, ok := seBadRequest.(error)

	// then
	require.NotNil(t, err)
	require.True(t, ok)
}

func TestStanzaError_Reason(t *testing.T) {
	// given
	msg := testMessageStanza()

	seBadRequest := E(BadRequest, msg)
	seConflict := E(Conflict, msg)
	seFeatureNotImplemented := E(FeatureNotImplemented, msg)
	seForbidden := E(Forbidden, msg)
	seGone := E(Gone, msg)
	seInternalServerError := E(InternalServerError, msg)
	seItemNotFound := E(ItemNotFound, msg)
	seJIDMalformed := E(JIDMalformed, msg)
	seNotAcceptable := E(NotAcceptable, msg)
	seNotAllowed := E(NotAllowed, msg)
	seNotAuthorized := E(NotAuthorized, msg)
	sePaymentRequired := E(PaymentRequired, msg)
	seRecipientUnavailabled := E(RecipientUnavailable, msg)
	seRedirect := E(Redirect, msg)
	seRegistrationRequired := E(RegistrationRequired, msg)
	seRemoteServerNotFound := E(RemoteServerNotFound, msg)
	seRemoteServerTimeout := E(RemoteServerTimeout, msg)
	seResourceConstraint := E(ResourceConstraint, msg)
	seServiceUnavailable := E(ServiceUnavailable, msg)
	seSubscriptionRequired := E(SubscriptionRequired, msg)
	seUndefinedCondition := E(UndefinedCondition, msg)
	seUnexpectedCondition := E(UnexpectedCondition, msg)
	seUnexpectedRequest := E(UnexpectedRequest, msg)

	// then
	require.Equal(t, 400, seBadRequest.Reason.Code())
	require.Equal(t, Modify, seBadRequest.Reason.Type())

	require.Equal(t, 409, seConflict.Reason.Code())
	require.Equal(t, Cancel, seConflict.Reason.Type())

	require.Equal(t, 501, seFeatureNotImplemented.Reason.Code())
	require.Equal(t, Cancel, seFeatureNotImplemented.Reason.Type())

	require.Equal(t, 403, seForbidden.Reason.Code())
	require.Equal(t, Auth, seForbidden.Reason.Type())

	require.Equal(t, 302, seGone.Reason.Code())
	require.Equal(t, Modify, seGone.Reason.Type())

	require.Equal(t, 500, seInternalServerError.Reason.Code())
	require.Equal(t, Wait, seInternalServerError.Reason.Type())

	require.Equal(t, 404, seItemNotFound.Reason.Code())
	require.Equal(t, Cancel, seItemNotFound.Reason.Type())

	require.Equal(t, 400, seJIDMalformed.Reason.Code())
	require.Equal(t, Modify, seJIDMalformed.Reason.Type())

	require.Equal(t, 406, seNotAcceptable.Reason.Code())
	require.Equal(t, Modify, seNotAcceptable.Reason.Type())

	require.Equal(t, 405, seNotAllowed.Reason.Code())
	require.Equal(t, Cancel, seNotAllowed.Reason.Type())

	require.Equal(t, 405, seNotAuthorized.Reason.Code())
	require.Equal(t, Auth, seNotAuthorized.Reason.Type())

	require.Equal(t, 402, sePaymentRequired.Reason.Code())
	require.Equal(t, Auth, sePaymentRequired.Reason.Type())

	require.Equal(t, 404, seRecipientUnavailabled.Reason.Code())
	require.Equal(t, Wait, seRecipientUnavailabled.Reason.Type())

	require.Equal(t, 302, seRedirect.Reason.Code())
	require.Equal(t, Modify, seRedirect.Reason.Type())

	require.Equal(t, 407, seRegistrationRequired.Reason.Code())
	require.Equal(t, Auth, seRegistrationRequired.Reason.Type())

	require.Equal(t, 404, seRemoteServerNotFound.Reason.Code())
	require.Equal(t, Cancel, seRemoteServerNotFound.Reason.Type())

	require.Equal(t, 504, seRemoteServerTimeout.Reason.Code())
	require.Equal(t, Wait, seRemoteServerTimeout.Reason.Type())

	require.Equal(t, 500, seResourceConstraint.Reason.Code())
	require.Equal(t, Wait, seResourceConstraint.Reason.Type())

	require.Equal(t, 503, seServiceUnavailable.Reason.Code())
	require.Equal(t, Cancel, seServiceUnavailable.Reason.Type())

	require.Equal(t, 407, seSubscriptionRequired.Reason.Code())
	require.Equal(t, Auth, seSubscriptionRequired.Reason.Type())

	require.Equal(t, 500, seUndefinedCondition.Reason.Code())
	require.Equal(t, Wait, seUndefinedCondition.Reason.Type())

	require.Equal(t, 400, seUnexpectedCondition.Reason.Code())
	require.Equal(t, Wait, seUnexpectedCondition.Reason.Type())

	require.Equal(t, 400, seUnexpectedRequest.Reason.Code())
	require.Equal(t, Cancel, seUnexpectedRequest.Reason.Type())
}

func TestStanzaError_Element(t *testing.T) {
	// given
	msg := testMessageStanza()

	// when
	se := E(InternalServerError, msg)
	se.Lang = "es"
	se.Text = "Error interno de servidor"
	se.ApplicationElement = stravaganza.NewBuilder("app-specific").
		WithAttribute(stravaganza.Namespace, "app-ns").
		Build()

	seElem := se.Element()

	// then
	expectedOutput := "<message from='ortuman@jackal.im/balcony' to='noelia@jackal.im/yard' type='error'><body>Hi everyone!</body><error code='500' type='wait'><internal-server-error xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/><text xmlns='urn:ietf:params:xml:ns:xmpp-stanzas' xml:lang='es'>Error interno de servidor</text><app-specific xmlns='app-ns'/></error></message>"
	require.Equal(t, expectedOutput, seElem.String())
}

func TestStanzaError_Stanza(t *testing.T) {
	// given
	msg := testMessageStanza()

	// when
	se := E(InternalServerError, msg)
	se.Lang = "es"
	se.Text = "Error interno de servidor"
	se.ApplicationElement = stravaganza.NewBuilder("app-specific").
		WithAttribute(stravaganza.Namespace, "app-ns").
		Build()

	seStanza, _ := se.Stanza(true)

	// then
	expectedOutput := "<message from='ortuman@jackal.im/balcony' to='noelia@jackal.im/yard' type='error'><body>Hi everyone!</body><error code='500' type='wait'><internal-server-error xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/><text xmlns='urn:ietf:params:xml:ns:xmpp-stanzas' xml:lang='es'>Error interno de servidor</text><app-specific xmlns='app-ns'/></error></message>"
	require.Equal(t, expectedOutput, seStanza.String())
}

func testMessageStanza() *stravaganza.Message {
	b := stravaganza.NewMessageBuilder()
	b.WithValidateJIDs(true)
	b.WithAttribute("from", "noelia@jackal.im/yard")
	b.WithAttribute("to", "ortuman@jackal.im/balcony")
	b.WithChild(
		stravaganza.NewBuilder("body").
			WithText("Hi everyone!").
			Build(),
	)
	msg, _ := b.BuildMessage()
	return msg
}
