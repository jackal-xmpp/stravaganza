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

// MessageName represents 'message' stanza type name.
const MessageName = "message"

const (
	// NormalType represents a 'normal' message type.
	NormalType = "normal"

	// HeadlineType represents a 'headline' message type.
	HeadlineType = "headline"

	// ChatType represents a 'chat' message type.
	ChatType = "chat"

	// GroupChatType represents a 'groupchat' message type.
	GroupChatType = "groupchat"
)

// Message type represents a <message> element.
type Message struct {
	stanza
}

// IsNormal returns true if this is a 'normal' type Message.
func (m *Message) IsNormal() bool {
	return m.Type() == NormalType || m.Type() == ""
}

// IsHeadline returns true if this is a 'headline' type Message.
func (m *Message) IsHeadline() bool {
	return m.Type() == HeadlineType
}

// IsChat returns true if this is a 'chat' type Message.
func (m *Message) IsChat() bool {
	return m.Type() == ChatType
}

// IsGroupChat returns true if this is a 'groupchat' type Message.
func (m *Message) IsGroupChat() bool {
	return m.Type() == GroupChatType
}

// IsMessageWithBody returns true if the message has a body sub element.
func (m *Message) IsMessageWithBody() bool {
	return m.Child("body") != nil
}
