//	Copyright (C) 2018-2020 CornierKhan1
//
//	WiiSOAP is SOAP Server Software, designed specifically to handle Wii Shop Channel SOAP.
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as published
//    by the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see http://www.gnu.org/licenses/.

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io"
	"regexp"
	"time"
)

var namespaceParse = regexp.MustCompile(`^urn:(.{3})\.wsapi\.broadon\.com/(.*)$`)

// parseAction interprets contents along the lines of "urn:ecs.wsapi.broadon.com/CheckDeviceStatus",
//where "CheckDeviceStatus" is the action to be performed.
func parseAction(original string) (string, string) {
	// Intended to return the original string, the service's name and the name of the action.
	matches := namespaceParse.FindStringSubmatch(original)
	if len(matches) != 3 {
		// It seems like the passed action was not matched properly.
		return "", ""
	}

	service := matches[1]
	action := matches[2]
	return service, action
}

// NewEnvelope returns a new Envelope with proper attributes initialized.
func NewEnvelope(service string, action string) Envelope {
	// Get a sexy new timestamp to use.
	timestampNano := fmt.Sprint(time.Now().UTC().UnixNano())[0:13]

	return Envelope{
		SOAPEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		XSD:     "http://www.w3.org/2001/XMLSchema",
		XSI:     "http://www.w3.org/2001/XMLSchema-instance",
		Body: Body{
			Response: Response{
				XMLName: xml.Name{Local: action + "Response"},
				XMLNS:   "urn:" + service + ".wsapi.broadon.com",

				TimeStamp: timestampNano,
			},
		},
		action: action,
	}
}

// Action returns the action for this service.
func (e *Envelope) Action() string {
	return e.action
}

// Timestamp returns a shared timestamp for this request.
func (e *Envelope) Timestamp() string {
	return e.Body.Response.TimeStamp
}

// obtainCommon interprets a given node, and updates the envelope with common key values.
func (e *Envelope) ObtainCommon(doc *xmlquery.Node) error {
	var err error

	// These fields are common across all requests.
	e.Body.Response.Version, err = getKey(doc, "Version")
	if err != nil {
		return err
	}
	e.Body.Response.DeviceId, err = getKey(doc, "DeviceId")
	if err != nil {
		return err
	}
	e.Body.Response.MessageId, err = getKey(doc, "MessageId")
	if err != nil {
		return err
	}

	return nil
}

// AddKVNode adds a given key by name to a specified value, such as <key>value</key>.
func (e *Envelope) AddKVNode(key string, value string) {
	e.Body.Response.CustomFields = append(e.Body.Response.CustomFields, KVField{
		XMLName: xml.Name{Local: key},
		Value:   value,
	})
}

// AddCustomType adds a given key by name to a specified structure.
func (e *Envelope) AddCustomType(customType interface{}) {
	e.Body.Response.CustomFields = append(e.Body.Response.CustomFields, customType)
}

// becomeXML marshals the Envelope object, returning the intended boolean state on success.
// ..there has to be a better way to do this, TODO.
func (e *Envelope) becomeXML(intendedStatus bool) (bool, string) {
	contents, err := xml.Marshal(e)
	if err != nil {
		return false, "an error occurred marshalling XML: " + err.Error()
	} else {
		// Add XML header on top of existing contents.
		result := xml.Header + string(contents)
		return intendedStatus, result
	}
}

// ReturnSuccess returns a standard SOAP response with a positive error code.
func (e *Envelope) ReturnSuccess() (bool, string) {
	// Ensure the error code is 0.
	e.Body.Response.ErrorCode = 0

	return e.becomeXML(true)
}

// formatError returns a standard SOAP response with an error code.
func (e *Envelope) ReturnError(errorCode int, reason string, err error) (bool, string) {
	e.Body.Response.ErrorCode = errorCode

	// Ensure all additional fields are empty to avoid conflict.
	e.Body.Response.CustomFields = nil

	e.AddKVNode("UserReason", reason)
	e.AddKVNode("ServerReason", err.Error())

	return e.becomeXML(false)
}

// normalise parses a document, returning a document with only the request type's child nodes, stripped of prefix.
func normalise(service string, action string, reader io.Reader) (*xmlquery.Node, error) {
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		return nil, err
	}

	// Find the keys for this element named after the action.
	result := doc.SelectElement("//" + service + ":" + action)
	if result == nil {
		return nil, errors.New("missing root node")
	}
	stripNamespace(result)

	return result, nil
}

// stripNamespace removes a prefix from nodes, changing a key from "ias:Version" to "Version".
// It is based off of https://github.com/antchfx/xmlquery/issues/15#issuecomment-567575075.
func stripNamespace(node *xmlquery.Node) {
	node.Prefix = ""

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		stripNamespace(child)
	}
}

// getKey returns the value for a child key from a node, if documented.
func getKey(doc *xmlquery.Node, key string) (string, error) {
	node := xmlquery.FindOne(doc, "//"+key)

	if node == nil {
		return "", errors.New("missing mandatory key named " + key)
	} else {
		return node.InnerText(), nil
	}
}
