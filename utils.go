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

import "strings"

func namespaceForType(service string) string {
	return "urn:" + service + ".wsapi.broadon.com"
}

// Expected contents are along the lines of "urn:ecs.wsapi.broadon.com/CheckDeviceStatus"
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
