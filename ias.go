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

func iasHandler(common map[string]string, doc *xmlquery.Node) (bool, string) {

	// All actions below are for IAS-related functions.
	switch common["Action"] {
	// TODO: Make the case functions cleaner. (e.g. Should the response be a variable?)
	// TODO: Update the responses so that they query the SQL Database for the proper information (e.g. Device Code, Token, etc).

	case "CheckRegistration":
		serialNo, err := getKey(doc, "SerialNumber")
		if err != nil {
			return formatError(common, 5, "not good enough for me. ;3", err)
		}

		fmt.Println("The request is valid! Responding...")
		custom := fmt.Sprintf(`<OriginalSerialNumber>%s</OriginalSerialNumber>
			<DeviceStatus>R</DeviceStatus>`, serialNo)
		return formatSuccess(common, custom)

	case "GetRegistrationInfo":
		accountId, err := getKey(doc, "AccountId")
		if err != nil {
			return formatError(common, 6, "how dirty. ;3", err)
		}

		country, err := getKey(doc, "Country")
		if err != nil {
			return formatError(common, 6, "how dirty. ;3", err)
		}

		fmt.Println("The request is valid! Responding...")
		custom := fmt.Sprintf(`<AccountId>%s</AccountId>
			<DeviceToken>00000000</DeviceToken>
			<DeviceTokenExpired>false</DeviceTokenExpired>
			<Country>%s</Country>
			<ExtAccountId></ExtAccountId>
			<DeviceCode>0000000000000000</DeviceCode>
			<DeviceStatus>R</DeviceStatus>
			<Currency>POINTS</Currency>`, accountId, country)
		return formatSuccess(common, custom)

	case "Register":
		accountId, err := getKey(doc, "AccountId")
		if err != nil {
			return formatError(common, 7, "disgustingly invalid. ;3", err)
		}

		country, err := getKey(doc, "Country")
		if err != nil {
			return formatError(common, 7, "disgustingly invalid. ;3", err)
		}

		fmt.Println("The request is valid! Responding...")
		custom := fmt.Sprintf(`<AccountId>%s</AccountId>
			<DeviceToken>00000000</DeviceToken>
			<Country>%s</Country>
			<ExtAccountId></ExtAccountId>
			<DeviceCode>00000000</DeviceCode>`, accountId, country)
		return formatSuccess(common, custom)

	case "Unregister":
		// how abnormal... ;3

		fmt.Println("The request is valid! Responding...")
		return formatSuccess(common, "")

	default:
		return false, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer."
	}
}
