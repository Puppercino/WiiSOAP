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
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io"
	"regexp"
	"strconv"
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

// formatHeader formats a response type and the proper service.
func formatHeader(responseType string, service string) string {
	return fmt.Sprintf(Header, responseType, "urn:"+service+".wsapi.broadon.com")
}

// formatTemplate inserts common, cross-requests values into every request.
func formatTemplate(common map[string]string, errorCode int) string {
	return fmt.Sprintf(Template, common["Version"], common["DeviceID"], common["MessageId"], common["Timestamp"], errorCode)
}

// formatFooter formats the closing tags of any SOAP request per previous response type.
func formatFooter(responseType string) string {
	return fmt.Sprintf(Footer, responseType)
}

// formatForNamespace mangles together several variables throughout a SOAP request.
func formatForNamespace(common map[string]string, errorCode int, extraContents string) string {
	return fmt.Sprintf("%s%s%s%s",
		formatHeader(common["Action"], common["Service"]),
		formatTemplate(common, errorCode),
		"\t"+extraContents+"\n",
		formatFooter(common["Action"]),
	)
}

// formatSuccess returns a standard SOAP response with a positive error code, and additional contents.
func formatSuccess(common map[string]string, extraContents string) (bool, string) {
	return true, formatForNamespace(common, 0, extraContents)
}

// formatError returns a standard SOAP response with an error code.
func formatError(common map[string]string, errorCode int, reason string, err error) (bool, string) {
	extra := "<UserReason>" + reason + "</UserReason>\n" +
		"\t<ServerReason>" + err.Error() + "</ServerReason>\n"
	return false, formatForNamespace(common, errorCode, extra)
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

// obtainCommon interprets a given node, and from its children finds common keys and respective values across all requests.
func obtainCommon(doc *xmlquery.Node) (map[string]string, error) {
	info := make(map[string]string)

	// Get a sexy new timestamp to use.
	timestampNano := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	info["Timestamp"] = timestampNano + "000"

	// These fields are common across all requests.
	// Looping through all...
	shared := []string{"Version", "DeviceId", "MessageId"}
	for _, key := range shared {
		// select their node by name...
		value, err := getKey(doc, key)
		if err != nil {
			return nil, err
		}

		// and insert their value to our map.
		info[key] = value
	}

	return info, nil
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
