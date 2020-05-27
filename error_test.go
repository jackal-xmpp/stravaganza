/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorMessage(t *testing.T) {
	require.Equal(t, badRequestErrorReason, ErrBadRequest.Error())
	require.Equal(t, conflictErrorReason, ErrConflict.Error())
	require.Equal(t, featureNotImplementedErrorReason, ErrFeatureNotImplemented.Error())
	require.Equal(t, forbiddenErrorReason, ErrForbidden.Error())
	require.Equal(t, goneErrorReason, ErrGone.Error())
	require.Equal(t, internalServerErrorErrorReason, ErrInternalServerError.Error())
	require.Equal(t, itemNotFoundErrorReason, ErrItemNotFound.Error())
	require.Equal(t, jidMalformedErrorReason, ErrJidMalformed.Error())
	require.Equal(t, notAcceptableErrorReason, ErrNotAcceptable.Error())
	require.Equal(t, notAllowedErrorReason, ErrNotAllowed.Error())
	require.Equal(t, notAuthroizedErrorReason, ErrNotAuthorized.Error())
	require.Equal(t, paymentRequiredErrorReason, ErrPaymentRequired.Error())
	require.Equal(t, recipientUnavailableErrorReason, ErrRecipientUnavailable.Error())
	require.Equal(t, redirectErrorReason, ErrRedirect.Error())
	require.Equal(t, registrationRequiredErrorReason, ErrRegistrationRequired.Error())
	require.Equal(t, remoteServerNotFoundErrorReason, ErrRemoteServerNotFound.Error())
	require.Equal(t, remoteServerTimeoutErrorReason, ErrRemoteServerTimeout.Error())
	require.Equal(t, resourceConstraintErrorReason, ErrResourceConstraint.Error())
	require.Equal(t, serviceUnavailableErrorReason, ErrServiceUnavailable.Error())
	require.Equal(t, subscriptionRequiredErrorReason, ErrSubscriptionRequired.Error())
	require.Equal(t, undefinedConditionErrorReason, ErrUndefinedCondition.Error())
	require.Equal(t, unexpectedConditionErrorReason, ErrUnexpectedCondition.Error())
	require.Equal(t, unexpectedRequestErrorReason, ErrUnexpectedRequest.Error())
}

func TestErrorElement(t *testing.T) {
	s, _ := NewBuilderFromElement(nil).
		WithName("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "ortuman@jackal.im").
		BuildStanza(false)

	require.NotNil(t, BadRequestError(s).Error().Child(badRequestErrorReason))
	require.NotNil(t, ConflictError(s).Error().Child(conflictErrorReason))
	require.NotNil(t, FeatureNotImplementedError(s).Error().Child(featureNotImplementedErrorReason))
	require.NotNil(t, ForbiddenError(s).Error().Child(forbiddenErrorReason))
	require.NotNil(t, GoneError(s).Error().Child(goneErrorReason))
	require.NotNil(t, InternalServerError(s).Error().Child(internalServerErrorErrorReason))
	require.NotNil(t, ItemNotFoundError(s).Error().Child(itemNotFoundErrorReason))
	require.NotNil(t, JidMalformedError(s).Error().Child(jidMalformedErrorReason))
	require.NotNil(t, NotAcceptableError(s).Error().Child(notAcceptableErrorReason))
	require.NotNil(t, NotAllowedError(s).Error().Child(notAllowedErrorReason))
	require.NotNil(t, NotAuthorizedError(s).Error().Child(notAuthroizedErrorReason))
	require.NotNil(t, PaymentRequiredError(s).Error().Child(paymentRequiredErrorReason))
	require.NotNil(t, RecipientUnavailableError(s).Error().Child(recipientUnavailableErrorReason))
	require.NotNil(t, RedirectError(s).Error().Child(redirectErrorReason))
	require.NotNil(t, RegistrationRequiredError(s).Error().Child(registrationRequiredErrorReason))
	require.NotNil(t, RemoteServerNotFoundError(s).Error().Child(remoteServerNotFoundErrorReason))
	require.NotNil(t, RemoteServerTimeoutError(s).Error().Child(remoteServerTimeoutErrorReason))
	require.NotNil(t, ResourceConstraintError(s).Error().Child(resourceConstraintErrorReason))
	require.NotNil(t, ServiceUnavailableError(s).Error().Child(serviceUnavailableErrorReason))
	require.NotNil(t, SubscriptionRequiredError(s).Error().Child(subscriptionRequiredErrorReason))
	require.NotNil(t, UndefinedConditionError(s).Error().Child(undefinedConditionErrorReason))
	require.NotNil(t, UnexpectedConditionError(s).Error().Child(unexpectedConditionErrorReason))
	require.NotNil(t, UnexpectedRequestError(s).Error().Child(unexpectedRequestErrorReason))
}

func TestErrorJIDAttributes(t *testing.T) {
	s, _ := NewBuilderFromElement(nil).
		WithName("message").
		WithAttribute("from", "ortuman@jackal.im/yard").
		WithAttribute("to", "ortuman@jackal.im").
		BuildStanza(false)

	sErr := BadRequestError(s)

	require.Equal(t, "ortuman@jackal.im/yard", sErr.Attribute("to").Value)
	require.Equal(t, "ortuman@jackal.im", sErr.Attribute("from").Value)
}
