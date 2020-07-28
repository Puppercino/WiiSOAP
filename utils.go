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
	"fmt"
	"strconv"
	"strings"
	"time"
)

// namespaceForType returns the expected XML namespace format for a service.
func namespaceForType(service string) string {
	return "urn:" + service + ".wsapi.broadon.com"
}

// parseAction interprets contents along the lines of "urn:ecs.wsapi.broadon.com/CheckDeviceStatus".
func parseAction(original string, service string) string {
	prefix := namespaceForType(service) + "/"
	stripped := strings.Replace(original, prefix, "", 1)

	if stripped == original {
		// This doesn't appear valid.
		return ""
	} else {
		return stripped
	}
}

// formatHeader formats a response type and the proper service.
func formatHeader(responseType string, service string) string {
	return fmt.Sprintf(Header, responseType, namespaceForType(service))
}

// formatTemplate inserts common, cross-requests values into every request.
func formatTemplate(version string, deviceId string, messageId string, errorCode int) string {
	// Get a sexy new timestamp to use.
	timestampNano := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	timestamp := timestampNano + "000"

	return fmt.Sprintf(Template, version, deviceId, messageId, timestamp, errorCode)
}

// formatFooter formats the closing tags of any SOAP request per previous response type.
func formatFooter(responseType string) string {
	return fmt.Sprintf(Footer, responseType)
}

// formatForNamespace mangles together several variables throughout a SOAP request.
func formatForNamespace(service string, responseType string, version string, deviceId string, messageId string, errorCode int, extraContents string) string {
	return fmt.Sprintf("%s%s%s%s",
		formatHeader(responseType, service),
		formatTemplate(version, deviceId, messageId, errorCode),
		"\t\t"+extraContents,
		formatFooter(responseType),
	)
}

// formatSuccess returns a standard SOAP response with a positive error code, and additional contents.
func formatSuccess(service string, responseType string, version string, deviceId string, messageId string, extraContents string) string {
	return formatForNamespace(service, responseType, version, deviceId, messageId, 0, extraContents)
}

// formatError returns a standard SOAP response with an error code.
func formatError(service string, responseType string, version string, deviceId string, messageId string, errorCode int) string {
	return formatForNamespace(service, responseType, version, deviceId, messageId, errorCode, "")
}
