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
	"github.com/antchfx/xmlquery"
)

func iasHandler(e Envelope, doc *xmlquery.Node) (bool, string) {

	// All actions below are for IAS-related functions.
	switch e.Action() {
	// TODO: Make the case functions cleaner. (e.g. Should the response be a variable?)
	// TODO: Update the responses so that they query the SQL Database for the proper information (e.g. Device Code, Token, etc).

	case "CheckRegistration":
		serialNo, err := getKey(doc, "SerialNumber")
		if err != nil {
			return e.ReturnError(5, "not good enough for me. ;3", err)
		}

		fmt.Println("The request is valid! Responding...")
		e.AddKVNode("OriginalSerialNumber", serialNo)
		e.AddKVNode("DeviceStatus", "R")
		break

	case "GetChallenge":
		fmt.Println("The request is valid! Responding...")
		// The official Wii Shop Channel requests a Challenge from the server, and promptly disregards it.
		// (Sometimes, it may not request a challenge at all.) No attempt is made to validate the response.
		// It then uses another hard-coded value in place of this returned value entirely in any situation.
		// For this reason, we consider it irrelevant.
		e.AddKVNode("Challenge", SharedChallenge)
		break

	case "GetRegistrationInfo":
		accountId, err := getKey(doc, "AccountId")
		if err != nil {
			return e.ReturnError(7, "how dirty. ;3", err)
		}

		country, err := getKey(doc, "Country")
		if err != nil {
			return e.ReturnError(7, "how dirty. ;3", err)
		}

		fmt.Println("The request is valid! Responding...")
		e.AddKVNode("AccountId", accountId)
		e.AddKVNode("DeviceToken", "00000000")
		e.AddKVNode("DeviceTokenExpired", "false")
		e.AddKVNode("Country", country)
		e.AddKVNode("ExtAccountId", "")
		e.AddKVNode("DeviceCode", "0000000000000000")
		e.AddKVNode("DeviceStatus", "R")
		// This _must_ be POINTS.
		e.AddKVNode("Currency", "POINTS")
		break

	case "Register":
		accountId, err := getKey(doc, "AccountId")
		if err != nil {
			return e.ReturnError(8, "disgustingly invalid. ;3", err)
		}

		country, err := getKey(doc, "Country")
		if err != nil {
			return e.ReturnError(8, "disgustingly invalid. ;3", err)
		}

		fmt.Println("The request is valid! Responding...")
		e.AddKVNode("AccountId", accountId)
		e.AddKVNode("DeviceToken", "00000000")
		e.AddKVNode("Country", country)
		e.AddKVNode("ExtAccountId", "")
		e.AddKVNode("DeviceCode", "0000000000000000")
		break

	case "Unregister":
		// how abnormal... ;3
		fmt.Println("The request is valid! Responding...")
		break

	default:
		return false, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer."
	}

	return e.ReturnSuccess()
}
