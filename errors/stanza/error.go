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
	"strconv"

	"github.com/jackal-xmpp/stravaganza"
)

const xmppStanzaNamespace = "urn:ietf:params:xml:ns:xmpp-stanzas"

// Type is the stanza error type.
type Type uint8

// String returns Type string representation.
func (t Type) String() string { return type2Str[t] }

const (
	// Auth represents 'auth' error type.
	Auth Type = iota

	// Cancel represents 'cancel' error type.
	Cancel

	// Modify represents 'modify' error type.
	Modify

	// Wait represents 'wait' error type.
	Wait
)

// Reason is the stanza error reason.
type Reason uint8

// Code returns Reason associated code.
func (r Reason) Code() int { return reason2Code[r] }

// Type returns Reason associated error type.
func (r Reason) Type() Type { return reason2Type[r] }

// String returns Reason string representation.
func (r Reason) String() string { return reason2Str[r] }

const (
	// BadRequest error is returned when the sender has sent XML that is malformed or that cannot be processed.
	BadRequest Reason = iota

	// Conflict error is returned when access cannot be granted because an existing resource or session exists with
	// the same name or address.
	Conflict

	// FeatureNotImplemented is returned when the feature requested is not implemented by the server and therefore
	// cannot be processed.
	FeatureNotImplemented

	// Forbidden is returned when the requesting entity does not possess the required permissions to perform the action.
	Forbidden

	// Gone is returned when the recipient or server can no longer be contacted at this address.
	Gone

	// InternalServerError is returned when the server could not process the stanza because of a misconfiguration
	// or an otherwise-undefined internal server error.
	InternalServerError

	// ItemNotFound is returned when the addressed JID or item requested cannot be found.
	ItemNotFound

	// JIDMalformed is returned when the sending entity has provided or communicated an XMPP address or aspect thereof
	// that does not adhere to the syntax defined in https://xmpp.org/rfcs/rfc3920.html#addressing.
	JIDMalformed

	// NotAcceptable is returned when the sending entity has provided or communicated an XMPP address or aspect thereof
	// that does not adhere to the syntax defined in https://xmpp.org/rfcs/rfc3920.html#addressing.
	NotAcceptable

	// NotAllowed is returned when the recipient or server does not allow any entity to perform the action.
	NotAllowed

	// NotAuthorized is returned when the sender must provide proper credentials before being allowed to perform
	// the action, or has provided improper credentials.
	NotAuthorized

	// PaymentRequired is returned when the requesting entity is not authorized to access the requested service
	// because payment is required.
	PaymentRequired

	// RecipientUnavailable is returned when the intended recipient is temporarily unavailable.
	RecipientUnavailable

	// Redirect is redirecting requests for this information to another entity, usually temporarily.
	Redirect

	// RegistrationRequired is redirecting requests for this information to another entity, usually temporarily.
	RegistrationRequired

	// RemoteServerNotFound is returned when a remote server or service specified as part or all of the JID
	// of the intended recipient does not exist.
	RemoteServerNotFound

	// RemoteServerTimeout is returned when a remote server or service specified as part or all of the JID of
	// the intended recipient could not be contacted within a reasonable amount of time.
	RemoteServerTimeout

	// ResourceConstraint is returned when the server or recipient lacks the system resources necessary
	// to service the request.
	ResourceConstraint

	// ServiceUnavailable is returned when the server or recipient does not currently provide
	// the requested service.
	ServiceUnavailable

	// SubscriptionRequired is returned when the requesting entity is not authorized to access
	// the requested service because a subscription is required.
	SubscriptionRequired

	// UndefinedCondition is returned when the error condition is not one of those defined by the other conditions in
	// this list.
	UndefinedCondition

	// UnexpectedCondition is returned when the recipient or server understood the request but was not expecting it
	// at this time.
	UnexpectedCondition

	// UnexpectedRequest is returned by the stream when the recipient or server understood the request but was not
	// expecting it at this time.
	UnexpectedRequest
)

var type2Str = map[Type]string{
	Auth:   "auth",
	Cancel: "cancel",
	Modify: "modify",
	Wait:   "wait",
}

var reason2Code = map[Reason]int{
	BadRequest:            400,
	Conflict:              409,
	FeatureNotImplemented: 501,
	Forbidden:             403,
	Gone:                  302,
	InternalServerError:   500,
	ItemNotFound:          404,
	JIDMalformed:          400,
	NotAcceptable:         406,
	NotAllowed:            405,
	NotAuthorized:         405,
	PaymentRequired:       402,
	RecipientUnavailable:  404,
	Redirect:              302,
	RegistrationRequired:  407,
	RemoteServerNotFound:  404,
	RemoteServerTimeout:   504,
	ResourceConstraint:    500,
	ServiceUnavailable:    503,
	SubscriptionRequired:  407,
	UndefinedCondition:    500,
	UnexpectedCondition:   400,
	UnexpectedRequest:     400,
}

var reason2Type = map[Reason]Type{
	BadRequest:            Modify,
	Conflict:              Cancel,
	FeatureNotImplemented: Cancel,
	Forbidden:             Auth,
	Gone:                  Modify,
	InternalServerError:   Wait,
	ItemNotFound:          Cancel,
	JIDMalformed:          Modify,
	NotAcceptable:         Modify,
	NotAllowed:            Cancel,
	NotAuthorized:         Auth,
	PaymentRequired:       Auth,
	RecipientUnavailable:  Wait,
	Redirect:              Modify,
	RegistrationRequired:  Auth,
	RemoteServerNotFound:  Cancel,
	RemoteServerTimeout:   Wait,
	ResourceConstraint:    Wait,
	ServiceUnavailable:    Cancel,
	SubscriptionRequired:  Auth,
	UndefinedCondition:    Wait,
	UnexpectedCondition:   Wait,
	UnexpectedRequest:     Cancel,
}

var reason2Str = map[Reason]string{
	BadRequest:            "bad-request",
	Conflict:              "conflict",
	FeatureNotImplemented: "feature-not-implemented",
	Forbidden:             "forbidden",
	Gone:                  "gone",
	InternalServerError:   "internal-server-error",
	ItemNotFound:          "item-not-found",
	JIDMalformed:          "jid-malformed",
	NotAcceptable:         "not-acceptable",
	NotAllowed:            "not-allowed",
	NotAuthorized:         "not-authorized",
	PaymentRequired:       "payment-required",
	RecipientUnavailable:  "recipient-unavailable",
	Redirect:              "redirect",
	RegistrationRequired:  "registration-required",
	RemoteServerNotFound:  "remote-server-not-found",
	RemoteServerTimeout:   "remote-server-timeout",
	ResourceConstraint:    "resource-constraint",
	ServiceUnavailable:    "service-unavailable",
	SubscriptionRequired:  "subscription-required",
	UndefinedCondition:    "undefined-condition",
	UnexpectedCondition:   "unexpected-condition",
	UnexpectedRequest:     "unexpected-request",
}

// Error represents a stanza "error" element.
type Error struct {
	// Reason is the stanza error reason type.
	Reason Reason

	// SentElement is the original XMPP element that originated the stanza error.
	SentElement stravaganza.Element

	// Lang is the error text lang code.
	// If none assigned 'en' would be considered by default.
	Lang string

	// Text is the error descriptive text.
	Text string

	// ApplicationElement defines the application specific condition element.
	ApplicationElement stravaganza.Element
}

// Error method satisfies error interface.
func (se *Error) Error() string {
	return se.Reason.String()
}

// Element returns se XMPP generic element.
func (se *Error) Element() stravaganza.Element {
	return stravaganza.NewBuilderFromElement(se.SentElement).
		WithAttribute(stravaganza.Type, stravaganza.ErrorType).
		WithAttribute(stravaganza.From, se.SentElement.Attribute(stravaganza.To)).
		WithAttribute(stravaganza.To, se.SentElement.Attribute(stravaganza.From)).
		WithChild(se.errSubElement()).
		Build()
}

// Stanza returns se XMPP stanza element.
func (se *Error) Stanza(validateJIDs bool) (stravaganza.Stanza, error) {
	return stravaganza.NewBuilderFromElement(se.SentElement).
		WithAttribute(stravaganza.Type, stravaganza.ErrorType).
		WithAttribute(stravaganza.From, se.SentElement.Attribute(stravaganza.To)).
		WithAttribute(stravaganza.To, se.SentElement.Attribute(stravaganza.From)).
		WithChild(se.errSubElement()).
		WithValidateJIDs(validateJIDs).
		BuildStanza()
}

func (se *Error) errSubElement() stravaganza.Element {
	b := stravaganza.NewBuilder("error")
	b.WithAttribute("code", strconv.Itoa(se.Reason.Code()))
	b.WithAttribute("type", se.Reason.Type().String())
	b.WithChild(
		stravaganza.NewBuilder(se.Reason.String()).
			WithAttribute(stravaganza.Namespace, xmppStanzaNamespace).
			Build(),
	)
	if len(se.Text) > 0 {
		tb := stravaganza.NewBuilder("text").
			WithAttribute(stravaganza.Namespace, xmppStanzaNamespace)
		if len(se.Lang) > 0 {
			tb.WithAttribute(stravaganza.Language, se.Lang)
		} else {
			tb.WithAttribute(stravaganza.Language, "en")
		}
		tb.WithText(se.Text)

		b.WithChild(tb.Build())
	}
	if se.ApplicationElement != nil {
		b.WithChild(se.ApplicationElement)
	}
	return b.Build()
}

// E builds an error value from its arguments.
func E(reason Reason, sentElement stravaganza.Element) *Error {
	return &Error{Reason: reason, SentElement: sentElement}
}
