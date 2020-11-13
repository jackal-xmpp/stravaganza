package streamerror

import (
	"fmt"

	"github.com/jackal-xmpp/stravaganza"
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
	r Reason

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
		stravaganza.NewBuilder(se.r.String()).
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
		return fmt.Sprintf("%s: %v", se.r.String(), se.Err)
	}
	return se.r.String()
}

// E builds an stream error value from its arguments.
func E(reason Reason) *Error {
	return &Error{r: reason}
}
