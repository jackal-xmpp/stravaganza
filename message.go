/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

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
