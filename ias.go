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
	"crypto/md5"
	sha2562 "crypto/sha256"
	"errors"
	"fmt"
	"github.com/RiiConnect24/wiino/golang"
	"github.com/antchfx/xmlquery"
	"github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"strconv"
)

func iasHandler(e Envelope, doc *xmlquery.Node) (bool, string) {
	// All IAS-related functions should contain these keys.
	region, err := getKey(doc, "Region")
	if err != nil {
		return e.ReturnError(5, "not good enough for me. ;3", err)
	}
	country, err := getKey(doc, "Country")
	if err != nil {
		return e.ReturnError(5, "not good enough for me. ;3", err)
	}
	language, err := getKey(doc, "Language")
	if err != nil {
		return e.ReturnError(5, "not good enough for me. ;3", err)
	}

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
		reason := "how dirty. ;3"
		accountId, err := getKey(doc, "AccountId")
		if err != nil {
			return e.ReturnError(7, reason, err)
		}

		deviceCode, err := getKey(doc, "DeviceCode")
		if err != nil {
			return e.ReturnError(7, reason, err)
		}

		fmt.Println("The request is valid! Responding...")
		e.AddKVNode("AccountId", accountId)
		e.AddKVNode("DeviceToken", "00000000")
		e.AddKVNode("DeviceTokenExpired", "false")
		e.AddKVNode("Country", country)
		e.AddKVNode("ExtAccountId", "")
		e.AddKVNode("DeviceCode", deviceCode)
		e.AddKVNode("DeviceStatus", "R")
		// This _must_ be POINTS.
		e.AddKVNode("Currency", "POINTS")
		break

	case "Register":
		reason := "disgustingly invalid. ;3"
		deviceCode, err := getKey(doc, "DeviceCode")
		if err != nil {
			return e.ReturnError(7, reason, err)
		}

		registerRegion, err := getKey(doc, "RegisterRegion")
		if err != nil {
			return e.ReturnError(7, reason, err)
		}
		if registerRegion != region {
			return e.ReturnError(7, reason, errors.New("region does not match registration region"))
		}

		serialNo, err := getKey(doc, "SerialNumber")
		if err != nil {
			return e.ReturnError(7, reason, err)
		}

		// Validate given friend code.
		userId, err := strconv.ParseUint(deviceCode, 10, 64)
		if err != nil {
			return e.ReturnError(7, reason, err)
		}
		if wiino.NWC24CheckUserID(userId) != 0 {
			return e.ReturnError(7, reason, err)
		}

		// Generate a random 9-digit number, padding zeros as necessary.
		accountId := fmt.Sprintf("%9d", rand.Intn(999999999))

		// This is where it gets hairy.
		// Generate a device token, 21 characters...
		deviceToken := RandString(21)
		// ...and then its md5, because the Wii sends this...
		md5DeviceToken := fmt.Sprintf("%x", md5.Sum([]byte(deviceToken)))
		// ...and then the sha256 of that md5.
		// We'll store this in our database, as storing the md5 itself is effectively the token.
		// It would not be good for security to directly store the token either.
		// This is the hash of the md5 represented as a string, not individual byte values.
		doublyHashedDeviceToken := fmt.Sprintf("%x", sha2562.Sum256([]byte(md5DeviceToken)))

		// Insert all of our obtained values to the database..
		stmt, err := db.Prepare(`INSERT INTO wiisoap.userbase (DeviceId, DeviceToken, AccountId, Region, Country, Language, SerialNo, DeviceCode)  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Printf("error preparing statement: %v\n", err)
			return e.ReturnError(7, reason, errors.New("failed to prepare statement"))
		}
		_, err = stmt.Exec(e.DeviceId(), doublyHashedDeviceToken, accountId, region, country, language, serialNo, deviceCode)
		if err != nil {
			// It's okay if this isn't a MySQL error, as perhaps other issues have come in.
			if driverErr, ok := err.(*mysql.MySQLError); ok {
				if driverErr.Number == 1062 {
					return e.ReturnError(7, reason, errors.New("user already exists"))
				}
			}
			log.Printf("error executing statement: %v\n", err)
			return e.ReturnError(7, reason, errors.New("failed to execute db operation"))
		}

		fmt.Println("The request is valid! Responding...")
		e.AddKVNode("AccountId", accountId)
		e.AddKVNode("DeviceToken", deviceToken)
		e.AddKVNode("DeviceTokenExpired", "false")
		e.AddKVNode("Country", country)
		// Optionally, one can send back DeviceCode and ExtAccountId to update on device.
		// We send these back as-is regardless.
		e.AddKVNode("ExtAccountId", "")
		e.AddKVNode("DeviceCode", deviceCode)
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
