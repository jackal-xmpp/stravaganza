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

// copied from https://github.com/mellium/xmpp/blob/master/jid/jid.go

package jid

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/net/idna"
	"golang.org/x/text/secure/precis"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

var (
	errForbiddenLocalpart = errors.New("jid: localpart contains forbidden characters")
	errInvalidUTF8        = errors.New("jid: JID contains invalid UTF-8")
	errLongLocalpart      = errors.New("jid: localpart must be smaller than 1024 bytes")
	errInvalidDomainLen   = errors.New("jid: domainpart must be between 1 and 1023 bytes")
	errLongResourcepart   = errors.New("jid: resourcepart must be smaller than 1024 bytes")
	errInvalidIPv6        = errors.New("jid: domainpart is not a valid IPv6 address")
	errNoLocalpart        = errors.New("jid: localpart must be larger than 0 bytes")
	errNoResourcepart     = errors.New("jid: resourcepart must be larger than 0 bytes")
)

// MatchingOptions represents a matching jid mask.
type MatchingOptions int8

const (
	// MatchesNode indicates that left and right operand has same node value.
	MatchesNode = MatchingOptions(1)

	// MatchesDomain indicates that left and right operand has same domain value.
	MatchesDomain = MatchingOptions(2)

	// MatchesResource indicates that left and right operand has same resource value.
	MatchesResource = MatchingOptions(4)

	// MatchesBare indicates that left and right operand has same node and domain value.
	MatchesBare = MatchesNode | MatchesDomain

	// MatchesFull indicates that left and right operand has same node, domain and resource value.
	MatchesFull = MatchesNode | MatchesDomain | MatchesResource
)

// JID represents an XMPP address (JID).
// A JID is made up of a node (generally a username), a domain, and a resource.
// The node and resource are optional; domain is required.
type JID struct {
	node     string
	domain   string
	resource string
}

// New constructs a JID given a user, domain, and resource.
// This construction allows the caller to specify if stringprep should be applied or not.
func New(node, domain, resource string, skipStringPrep bool) (*JID, error) {
	if skipStringPrep {
		return &JID{
			node:     node,
			domain:   domain,
			resource: resource,
		}, nil
	}
	var j JID
	if err := j.stringPrep(node, domain, resource); err != nil {
		return nil, err
	}
	return &j, nil
}

// NewWithString constructs a JID from it's string representation.
// This construction allows the caller to specify if stringprep should be applied or not.
func NewWithString(str string, skipStringPrep bool) (*JID, error) {
	if len(str) == 0 {
		return &JID{}, nil
	}
	node, domain, resource, err := splitString(str, true)
	if err != nil {
		return nil, err
	}
	return New(node, domain, resource, skipStringPrep)
}

// Node returns the node, or empty string if this JID does not contain node information.
func (j *JID) Node() string {
	return j.node
}

// Domain returns the domain.
func (j *JID) Domain() string {
	return j.domain
}

// Resource returns the resource, or empty string if this JID does not contain resource information.
func (j *JID) Resource() string {
	return j.resource
}

// ToBareJID returns the JID equivalent of the bare JID, which is the JID with resource information removed.
func (j *JID) ToBareJID() *JID {
	if len(j.node) == 0 {
		return &JID{node: "", domain: j.domain, resource: ""}
	}
	return &JID{node: j.node, domain: j.domain, resource: ""}
}

// IsServer returns true if instance is a server JID.
func (j *JID) IsServer() bool {
	return len(j.node) == 0
}

// IsBare returns true if instance is a bare JID.
func (j *JID) IsBare() bool {
	return len(j.node) > 0 && len(j.resource) == 0
}

// IsFull returns true if instance is a full JID.
func (j *JID) IsFull() bool {
	return len(j.resource) > 0
}

// IsFullWithServer returns true if instance is a full server JID.
func (j *JID) IsFullWithServer() bool {
	return len(j.node) == 0 && len(j.resource) > 0
}

// IsFullWithUser returns true if instance is a full client JID.
func (j *JID) IsFullWithUser() bool {
	return len(j.node) > 0 && len(j.resource) > 0
}

// Matches tells whether or not j2 matches j.
func (j *JID) Matches(j2 *JID) bool {
	if j.IsFullWithUser() {
		return j.MatchesWithOptions(j2, MatchesNode|MatchesDomain|MatchesResource)
	} else if j.IsFullWithServer() {
		return j.MatchesWithOptions(j2, MatchesDomain|MatchesResource)
	} else if j.IsBare() {
		return j.MatchesWithOptions(j2, MatchesNode|MatchesDomain)
	}
	return j.MatchesWithOptions(j2, MatchesDomain)
}

// MatchesWithOptions tells whether two jids are equivalent based on matching options.
func (j *JID) MatchesWithOptions(j2 *JID, options MatchingOptions) bool {
	if (options&MatchesNode) > 0 && j.node != j2.node {
		return false
	}
	if (options&MatchesDomain) > 0 && j.domain != j2.domain {
		return false
	}
	if (options&MatchesResource) > 0 && j.resource != j2.resource {
		return false
	}
	return true
}

// String returns a string representation of the JID.
func (j *JID) String() string {
	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	if len(j.node) > 0 {
		buf.WriteString(j.node)
		buf.WriteString("@")
	}
	buf.WriteString(j.domain)
	if len(j.resource) > 0 {
		buf.WriteString("/")
		buf.WriteString(j.resource)
	}
	return buf.String()
}

func (j *JID) stringPrep(node, domain, resource string) error {
	// Ensure that parts are valid UTF-8 (and short circuit the rest of the
	// process if they're not). We'll check the domain after performing
	// the IDNA ToUnicode operation.
	if !utf8.ValidString(node) || !utf8.ValidString(resource) {
		return errInvalidUTF8
	}

	// RFC 7622 §3.2.1.  Preparation
	//
	//    An entity that prepares a string for inclusion in an XMPP domain
	//    slot MUST ensure that the string consists only of Unicode code points
	//    that are allowed in NR-LDH labels or U-labels as defined in
	//    [RFC5890].  This implies that the string MUST NOT include A-labels as
	//    defined in [RFC5890]; each A-label MUST be converted to a U-label
	//    during preparation of a string for inclusion in a domain slot.
	var err error
	domain, err = idna.ToUnicode(domain)
	if err != nil {
		return err
	}
	if !utf8.ValidString(domain) {
		return errInvalidUTF8
	}

	// RFC 7622 §3.2.2.  Enforcement
	//
	//   An entity that performs enforcement in XMPP domain slots MUST
	//   prepare a string as described in Section 3.2.1 and MUST also apply
	//   the normalization, case-mapping, and width-mapping rules defined in
	//   [RFC5892].
	//
	var nodeLen int
	data := make([]byte, 0, len(node)+len(domain)+len(resource))

	if node != "" {
		data, err = precis.UsernameCaseMapped.Append(data, []byte(node))
		if err != nil {
			return err
		}
		nodeLen = len(data)
	}
	data = append(data, []byte(domain)...)

	if resource != "" {
		data, err = precis.OpaqueString.Append(data, []byte(resource))
		if err != nil {
			return err
		}
	}
	if err := commonChecks(data[:nodeLen], domain, data[nodeLen+len(domain):]); err != nil {
		return err
	}
	j.node = string(data[:nodeLen])
	j.domain = string(data[nodeLen : nodeLen+len(domain)])
	j.resource = string(data[nodeLen+len(domain):])
	return nil
}

func splitString(s string, safe bool) (localpart, domainpart, resourcepart string, err error) {

	// RFC 7622 §3.1.  Fundamentals:
	//
	//    Implementation Note: When dividing a JID into its component parts,
	//    an implementation needs to match the separator characters '@' and
	//    '/' before applying any transformation algorithms, which might
	//    decompose certain Unicode code points to the separator characters.
	//
	// so let's do that now. First we'll parse the domainpart using the rules
	// defined in §3.2:
	//
	//    The domainpart of a JID is the portion that remains once the
	//    following parsing steps are taken:
	//
	//    1.  Remove any portion from the first '/' character to the end of the
	//        string (if there is a '/' character present).
	sep := strings.Index(s, "/")

	if sep == -1 {
		resourcepart = ""
	} else {
		// If the resource part exists, make sure it isn't empty.
		if safe && sep == len(s)-1 {
			err = errNoResourcepart
			return
		}
		resourcepart = s[sep+1:]
		if strings.Index(resourcepart, "@") != -1 {
			err = errNoResourcepart
			return
		}
		s = s[:sep]
	}

	//    2.  Remove any portion from the beginning of the string to the first
	//        '@' character (if there is an '@' character present).

	sep = strings.Index(s, "@")

	switch {
	case sep == -1:
		// There is no @ sign, and therefore no localpart.
		localpart = ""
		domainpart = s
	case safe && sep == 0:
		// The JID starts with an @ sign (invalid empty localpart)
		err = errNoLocalpart
		return
	default:
		domainpart = s[sep+1:]
		localpart = s[:sep]
	}

	// We'll throw out any trailing dots on domainparts, since they're ignored:
	//
	//    If the domainpart includes a final character considered to be a label
	//    separator (dot) by [RFC1034], this character MUST be stripped from
	//    the domainpart before the JID of which it is a part is used for the
	//    purpose of routing an XML stanza, comparing against another JID, or
	//    constructing an XMPP URI or IRI [RFC5122].  In particular, such a
	//    character MUST be stripped before any other canonicalization steps
	//    are taken.

	domainpart = strings.TrimSuffix(domainpart, ".")

	return
}

func commonChecks(localpart []byte, domainpart string, resourcepart []byte) error {
	err := localChecks(localpart)
	if err != nil {
		return err
	}

	err = resourceChecks(resourcepart)
	if err != nil {
		return err
	}

	return domainChecks(domainpart)
}

func localChecks(localpart []byte) error {
	if len(localpart) > 1023 {
		return errLongLocalpart
	}

	// RFC 7622 §3.3.1 provides a small table of characters which are still not
	// allowed in localpart's even though the IdentifierClass base class and the
	// UsernameCaseMapped profile don't forbid them; disallow them here.
	// We can't add them to the profiles disallowed characters because they get
	// checked before the profile is applied (so some characters may still be
	// normalized to characters in this set).
	if bytes.ContainsAny(localpart, `"&'/:<>@`) {
		return errForbiddenLocalpart
	}

	return nil
}

func resourceChecks(resourcepart []byte) error {
	if len(resourcepart) > 1023 {
		return errLongResourcepart
	}
	return nil
}

func domainChecks(domainpart string) error {
	if l := len(domainpart); l < 1 || l > 1023 {
		return errInvalidDomainLen
	}

	return checkIP6String(domainpart)
}

func checkIP6String(domain string) error {
	// if the domain is a valid IPv6 address (with brackets), short circuit.
	if l := len(domain); l > 2 && strings.HasPrefix(domain, "[") &&
		strings.HasSuffix(domain, "]") {
		if ip := net.ParseIP(domain[1 : l-1]); ip == nil || ip.To4() != nil {
			return errInvalidIPv6
		}
	}
	return nil
}
