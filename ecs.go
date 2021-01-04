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

func ecsHandler(e Envelope, doc *xmlquery.Node) (bool, string) {
	// All actions below are for ECS-related functions.
	switch e.Action() {
	// TODO: Make the case functions cleaner. (e.g. Should the response be a variable?)
	// TODO: Update the responses so that they query the SQL Database for the proper information (e.g. Device Code, Token, etc).

	case "CheckDeviceStatus":
		//You need to POST some SOAP from WSC if you wanna get some, honey. ;3

		fmt.Println("The request is valid! Responding...")
		e.AddCustomType(Balance{
			Amount:   2018,
			Currency: "POINTS",
		})
		e.AddKVNode("ForceSyncTime", "0")
		e.AddKVNode("ExtTicketTime", e.Timestamp())
		e.AddKVNode("SyncTime", e.Timestamp())
		break

	case "NotifyETicketsSynced":
		// This is a disgusting request, but 20 dollars is 20 dollars. ;3

		fmt.Println("The request is valid! Responding...")
		break

	case "ListETickets":
		// that's all you've got for me? ;3

		fmt.Println("The request is valid! Responding...")
		e.AddKVNode("ForceSyncTime", "0")
		e.AddKVNode("ExtTicketTime", e.Timestamp())
		e.AddKVNode("SyncTime", e.Timestamp())
		break

	case "PurchaseTitle":
		// If you wanna fun time, it's gonna cost ya extra sweetie... ;3

		fmt.Println("The request is valid! Responding...")
		e.AddCustomType(Balance{
			Amount:   2018,
			Currency: "POINTS",
		})
		e.AddCustomType(Transactions{
			TransactionId: "00000000",
			Date:          e.Timestamp(),
			Type:          "PURCHGAME",
		})
		e.AddKVNode("SyncTime", e.Timestamp())
		e.AddKVNode("Certs", "00000000")
		e.AddKVNode("TitleId", "00000000")
		e.AddKVNode("ETickets", "00000000")
		break

	default:
		return false, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer."
	}

	return e.ReturnSuccess()
}
