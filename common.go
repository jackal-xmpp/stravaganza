/*
 * Copyright (c) 2020 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

package stravaganza

func getProtoElementAttribute(pbElement *PBElement, name string) string {
	for _, attr := range pbElement.Attributes {
		if attr.Label == name {
			return attr.Value
		}
	}
	return ""
}
