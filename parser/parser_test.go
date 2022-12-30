// Copyright 2022 The jackal Authors
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

package xmppparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_ErrTooLargeStanzaRead(t *testing.T) {
	// given
	docSrc := `<a/><be/>`
	p := New(strings.NewReader(docSrc), SocketStream, 4)

	// when
	a, err0 := p.Parse()
	be, err1 := p.Parse()

	// then
	require.Nil(t, err0)
	require.NotNil(t, a)
	require.Equal(t, "<a/>", a.String())

	require.Nil(t, be)
	require.Equal(t, ErrTooLargeStanza, err1)
}

func TestParser_ParseSeveralElements(t *testing.T) {
	// given
	docSrc := `<?xml version="1.0" encoding="UTF-8"?><a/><b/><c/>`

	r := strings.NewReader(docSrc)
	p := New(r, DefaultMode, 1024)

	// when
	a, err1 := p.Parse()
	b, err2 := p.Parse()
	c, err3 := p.Parse()

	// then
	require.NotNil(t, a)
	require.Nil(t, err1)

	require.NotNil(t, b)
	require.Nil(t, err2)

	require.NotNil(t, c)
	require.Nil(t, err3)
}

func TestParser_DocChildElements(t *testing.T) {
	// given
	docSrc := `<parent><a/><b/><c/></parent>\n`
	p := New(strings.NewReader(docSrc), DefaultMode, 1024)

	// when
	elem, err := p.Parse()
	require.Nil(t, err)
	require.NotNil(t, elem)

	childs := elem.AllChildren()

	// then
	require.Equal(t, 3, len(childs))
	require.Equal(t, "a", childs[0].Name())
	require.Equal(t, "b", childs[1].Name())
	require.Equal(t, "c", childs[2].Name())
}

func TestParser_Stream(t *testing.T) {
	openStreamXML := `<stream:stream xmlns:stream="http://etherx.jabber.org/streams" version="1.0" xmlns="jabber:client" to="localhost" xml:lang="en" xmlns:xml="http://www.w3.org/XML/1998/namespace"> `
	p := New(strings.NewReader(openStreamXML), SocketStream, 1024)
	elem, err := p.Parse()
	require.Nil(t, err)
	require.Equal(t, "stream:stream", elem.Name())

	closeStreamXML := `</stream:stream> `
	p = New(strings.NewReader(closeStreamXML), SocketStream, 1024)

	_, err = p.Parse()

	require.Equal(t, ErrStreamClosedByPeer, err)
}

func BenchmarkParser_Parse(b *testing.B) {
	docSrc := "<iq id='config1' type='result' from='pubsub.shakespeare.lit' to='hamlet@denmark.lit/elsinore'><pubsub xmlns='http://jabber.org/protocol/pubsub#owner'><configure node='princely_musings'><x xmlns='jabber:x:data' type='form'><field type='hidden' var='FORM_TYPE'><value>http://jabber.org/protocol/pubsub#node_config</value></field><field type='text-single' label='The default language of the node' var='pubsub#language'/><field type='text-single' label='A friendly name for the node' var='pubsub#title'/><field type='text-single' label='A description of the node' var='pubsub#description'/><field type='boolean' label='Whether to deliver payloads with event notifications' var='pubsub#deliver_payloads'><value>false</value></field><field type='boolean' label='Whether to deliver event notifications' var='pubsub#deliver_notifications'><value>false</value></field><field type='boolean' label='Whether to notify subscribers when the node configuration changes' var='pubsub#notify_config'><value>false</value></field><field type='boolean' label='Whether to notify subscribers when the node is deleted' var='pubsub#notify_delete'><value>false</value></field><field type='boolean' label='Whether to notify subscribers when items are removed from the node' var='pubsub#notify_retract'><value>false</value></field><field type='boolean' label='Whether to notify owners about new subscribers and unsubscribes' var='pubsub#notify_sub'><value>false</value></field><field type='boolean' label='Whether to persist items to storage' var='pubsub#persist_items'><value>false</value></field><field type='text-single' label='The maximum number of items to persist. `max` for no specific limit other than a server imposed maximum.' var='pubsub#max_items'><value>120</value></field><field type='text-single' label='Number of seconds after which to automatically purge items. `max` for no specific limit other than a server imposed maximum.' var='pubsub#item_expire'/><field type='boolean' label='Whether to allow subscriptions' var='pubsub#subscribe'><value>false</value></field><field type='list-single' label='Who may subscribe and retrieve items' var='pubsub#access_model'><value>open</value></field><field type='list-multi' label='The roster group(s) allowed to subscribe and retrieve items' var='pubsub#roster_groups_allowed'/><field type='list-single' label='The publisher model' var='pubsub#publish_model'><value/></field><field type='boolean' label='Whether to purge all items when the relevant publisher goes offline' var='pubsub#purge_offline'><value>false</value></field><field type='text-single' label='The maximum payload size in bytes' var='pubsub#max_payload_size'><value>65536</value></field><field type='list-single' label='When to send the last published item' var='pubsub#send_last_published_item'><value/></field><field type='boolean' label='Whether to deliver notifications to available users only' var='pubsub#presence_based_delivery'><value>false</value></field><field type='list-single' label='Specify the delivery style for notifications' var='pubsub#notification_type'><value/></field><field type='text-single' label='The semantic type information of data in the node, usually specified by the namespace of the payload (if any)' var='pubsub#type'/><field type='text-single' label='The URL of an XSL transformation which can be applied to payloads in order to generate an appropriate message body element.' var='pubsub#body_xslt'/><field type='text-single' label='The URL of an XSL transformation which can be applied to the payload format in order to generate a valid Data Forms result that the client could display using a generic Data Forms rendering engine' var='pubsub#dataform_xslt'/></x></configure></pubsub></iq>"

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(docSrc)
		p := New(r, SocketStream, 64*1024)

		elem, err := p.Parse()

		require.NoError(b, err)
		require.NotNil(b, elem)
	}
}
