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
	"fmt"

	"github.com/jackal-xmpp/stravaganza/v2"
)

const xmppStanzaNamespace = "urn:ietf:params:xml:ns:xmpp-stanzas"

// Reason is the class of the stream error.
type Reason uint8

const (
	InvalidXML             Reason = iota // Invalid XML.
	InvalidNamespace                     // Invalid namespace.
	HostUnknown                          // Host unknown.
	InvalidFrom                          // Invalid from address.
	PolicyViolation                      // Policy violation.
	RemoteConnectionFailed               // Remote connection failed.
	Conflict                             // Conflict.
	ConnectionTimeout                    // Connection timeout.
	UnsupportedStanzaType                // Unsupported stanza type.
	UnsupportedVersion                   // Unsupported version.
	NotAuthorized                        // Not authorized.
	ResourceConstraint                   // Resource constraint.
	SystemShutdown                       // System shutdown.
	UndefinedCondition                   // Undefined condition.
	InternalServerError                  // Internal server error.
)

var reason2Str = map[Reason]string{
	InvalidXML:             "invalid-xml",
	InvalidNamespace:       "invalid-namespace",
	HostUnknown:            "host-unknown",
	InvalidFrom:            "invalid-from",
	PolicyViolation:        "policy-violation",
	RemoteConnectionFailed: "remote-connection-failed",
	Conflict:               "conflict",
	ConnectionTimeout:      "connection-timeout",
	UnsupportedStanzaType:  "unsupported-stanza-type",
	UnsupportedVersion:     "unsupported-version",
	NotAuthorized:          "not-authorized",
	ResourceConstraint:     "resource-constraint",
	SystemShutdown:         "system-shutdown",
	UndefinedCondition:     "undefined-condition",
	InternalServerError:    "internal-server-error",
}

// String returns string error reason representation.
func (r Reason) String() string { return reason2Str[r] }

// Error represents a stream "error" element.
type Error struct {
	// Reason is the stream error reason type.
	Reason Reason

	// Err contains the actual error that originated the stream error.
	Err error

	// Lang is the error text lang code.
	// If none assigned 'en' would be considered by default.
	Lang string

	// Text is the error descriptive text.
	Text string

	// ApplicationElement defines the application specific condition element.
	ApplicationElement stravaganza.Element
}

// Element returns stream error XML node.
func (se *Error) Element() stravaganza.Element {
	b := stravaganza.NewBuilder("stream:error")
	b.WithChild(
		stravaganza.NewBuilder(se.Reason.String()).
			WithAttribute(stravaganza.Namespace, xmppStanzaNamespace).
			Build(),
	)
	if se.ApplicationElement != nil {
		b.WithChild(se.ApplicationElement)
	}
	return b.Build()
}

// Error satisfies error interface.
func (se *Error) Error() string {
	if se.Err != nil {
		return fmt.Sprintf("%s: %v", se.Reason.String(), se.Err)
	}
	return se.Reason.String()
}

// E builds an stream error value from its arguments.
func E(reason Reason) *Error {
	return &Error{Reason: reason}
}
