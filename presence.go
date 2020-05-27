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
	"errors"
	"fmt"
	"strconv"
)

const capabilitiesNamespace = "http://jabber.org/protocol/caps"

// PresenceName represents 'presence' stanza type name.
const PresenceName = "presence"

const (
	// AvailableType represents an 'available' Presence type.
	AvailableType = ""

	// UnavailableType represents a 'unavailable' Presence type.
	UnavailableType = "unavailable"

	// SubscribeType represents a 'subscribe' Presence type.
	SubscribeType = "subscribe"

	// UnsubscribeType represents a 'unsubscribe' Presence type.
	UnsubscribeType = "unsubscribe"

	// SubscribedType represents a 'subscribed' Presence type.
	SubscribedType = "subscribed"

	// UnsubscribedType represents a 'unsubscribed' Presence type.
	UnsubscribedType = "unsubscribed"

	// ProbeType represents a 'probe' Presence type.
	ProbeType = "probe"
)

// ShowState represents Presence show state.
type ShowState int

const (
	// AvailableShowState represents 'available' Presence show state.
	AvailableShowState ShowState = iota

	// AwayShowState represents 'away' Presence show state.
	AwayShowState

	// ChatShowState represents 'chat' Presence show state.
	ChatShowState

	// DoNotDisturbShowState represents 'dnd' Presence show state.
	DoNotDisturbShowState

	// ExtendedAwaysShowState represents 'xa' Presence show state.
	ExtendedAwaysShowState
)

// Capabilities represents presence entity capabilities
type Capabilities struct {
	Node string
	Hash string
	Ver  string
}

// Presence type represents a <presence> element.
type Presence struct {
	stanza
	showState ShowState
	priority  int8
}

// IsAvailable returns true if this is an 'available' type Presence.
func (p *Presence) IsAvailable() bool {
	return p.Type() == AvailableType
}

// IsUnavailable returns true if this is an 'unavailable' type Presence.
func (p *Presence) IsUnavailable() bool {
	return p.Type() == UnavailableType
}

// IsSubscribe returns true if this is a 'subscribe' type Presence.
func (p *Presence) IsSubscribe() bool {
	return p.Type() == SubscribeType
}

// IsUnsubscribe returns true if this is an 'unsubscribe' type Presence.
func (p *Presence) IsUnsubscribe() bool {
	return p.Type() == UnsubscribeType
}

// IsSubscribed returns true if this is a 'subscribed' type Presence.
func (p *Presence) IsSubscribed() bool {
	return p.Type() == SubscribedType
}

// IsUnsubscribed returns true if this is an 'unsubscribed' type Presence.
func (p *Presence) IsUnsubscribed() bool {
	return p.Type() == UnsubscribedType
}

// IsProbe returns true if this is an 'probe' type Presence.
func (p *Presence) IsProbe() bool {
	return p.Type() == ProbeType
}

// Status returns presence stanza default status.
func (p *Presence) Status() string {
	if st := p.Child("status"); st != nil {
		return st.Text()
	}
	return ""
}

// ShowState returns presence stanza show state.
func (p *Presence) ShowState() ShowState {
	return p.showState
}

// Priority returns presence stanza priority value.
func (p *Presence) Priority() int8 {
	return p.priority
}

// Capabilities returns presence stanza capabilities element
func (p *Presence) Capabilities() *Capabilities {
	c := p.ChildNamespace("c", capabilitiesNamespace)
	if c == nil {
		return nil
	}
	return &Capabilities{
		Node: c.Attribute("node"),
		Hash: c.Attribute("hash"),
		Ver:  c.Attribute("ver"),
	}
}

func (p *Presence) validateStatus() error {
	sts := p.Children("status")
	for _, st := range sts {
		switch st.AttributeCount() {
		case 0:
			break
		case 1:
			attrs := st.AllAttributes()
			if attrs[0].Label == "xml:lang" {
				break
			}
			fallthrough
		default:
			return errors.New("stravaganza: presence <status/> element MUST NOT possess any attributes, with the exception of the 'xml:lang' attribute")
		}
	}
	return nil
}

func (p *Presence) setShow() error {
	shs := p.Children("show")
	switch len(shs) {
	case 0:
		p.showState = AvailableShowState
	case 1:
		if shs[0].AttributeCount() > 0 {
			return errors.New("stravaganza: presence <show/> element MUST NOT possess any attributes")
		}
		switch shs[0].Text() {
		case "away":
			p.showState = AwayShowState
		case "chat":
			p.showState = ChatShowState
		case "dnd":
			p.showState = DoNotDisturbShowState
		case "xa":
			p.showState = ExtendedAwaysShowState
		default:
			return fmt.Errorf("stravaganza: invalid presence show state: %s", shs[0].Text())
		}

	default:
		return errors.New("stravaganza: presence stanza MUST NOT contain more than one <show/> element")
	}
	return nil
}

func (p *Presence) setPriority() error {
	ps := p.Children("priority")
	switch len(ps) {
	case 0:
		break
	case 1:
		pr, err := strconv.Atoi(ps[0].Text())
		if err != nil {
			return err
		}
		if pr < -128 || pr > 127 {
			return errors.New("stravaganza: presence priority value MUST be an integer between -128 and +127")
		}
		p.priority = int8(pr)

	default:
		return errors.New("stravaganza: a presence stanza MUST NOT contain more than one <priority/> element")
	}
	return nil
}
