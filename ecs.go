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
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func ecsHandler(w http.ResponseWriter, r *http.Request) {
	// Figure out the action to handle via header.
	action := r.Header.Get("SOAPAction")
	action = parseAction(action, "ecs")

	// Get a sexy new timestamp to use.
	timestampNano := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	timestamp := timestampNano + "000"

	fmt.Println("[!] Incoming ECS request.")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body...", http.StatusInternalServerError)
	}

	// The switch converts the HTTP Body of the request into a string. There is no need to convert the cases to byte format.
	switch action {
	// TODO: Make the case functions cleaner. (e.g. Should the response be a variable?)
	// TODO: Update the responses so that they query the SQL Database for the proper information (e.g. Device Code, Token, etc).

	case "CheckDeviceStatus":
		fmt.Println("CDS.")
		CDS := CDS{}
		if err = xml.Unmarshal(body, &CDS); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "You need to POST some SOAP from WSC if you wanna get some, honey. ;3")
			return
		}
		fmt.Println(CDS)
		fmt.Println("The request is valid! Responding...")
		custom := fmt.Sprintf(`<Balance>
				<Amount>2018</Amount>
				<Currency>POINTS</Currency>
			</Balance>
			<ForceSyncTime>0</ForceSyncTime>
			<ExtTicketTime>%s</ExtTicketTime>
			<SyncTime>%s</SyncTime>`, timestamp, timestamp)
		fmt.Fprint(w, formatSuccess("ecs", action, CDS.Version, CDS.DeviceID, CDS.MessageID, custom))

	case "NotifiedETicketsSynced":
		fmt.Println("NETS")
		NETS := NETS{}
		if err = xml.Unmarshal(body, &NETS); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "This is a disgusting request, but 20 dollars is 20 dollars. ;3")
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(NETS)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprint(w, formatSuccess("ecs", action, NETS.Version, NETS.DeviceID, NETS.MessageID, ""))

	case "ListETickets":
		fmt.Println("LET")
		LET := LET{}
		if err = xml.Unmarshal(body, &LET); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "that's all you got for me? ;3")
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(LET)
		fmt.Println("The request is valid! Responding...")
		custom := fmt.Sprintf(`<ForceSyncTime>0</ForceSyncTime>
			<ExtTicketTime>%s</ExtTicketTime>
			<SyncTime>%s</SyncTime>`, timestamp, timestamp)
		fmt.Fprint(w, formatSuccess("ecs", action, LET.Version, LET.DeviceID, LET.MessageID, custom))

	case "PurchaseTitle":
		fmt.Println("PT")
		PT := PT{}
		if err = xml.Unmarshal(body, &PT); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "if you wanna fun time, its gonna cost ya extra sweetie. ;3")
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(PT)
		fmt.Println("The request is valid! Responding...")
		custom := fmt.Sprintf(`<Balance>
				<Amount>2018</Amount>
				<Currency>POINTS</Currency>
			</Balance>
			<Transactions>
				<TransactionId>00000000</TransactionId>
				<Date>%s</Date>
				<Type>PURCHGAME</Type>
			</Transactions>
			<SyncTime>%s</SyncTime>
			<ETickets>00000000</ETickets>
			<Certs>00000000</Certs>
			<Certs>00000000</Certs>
			<TitleId>00000000</TitleId>`, timestamp, timestamp)
		fmt.Fprint(w, formatSuccess("ecs", action, PT.Version, PT.DeviceID, PT.MessageID, custom))

	default:
		fmt.Fprint(w, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer.")
		return
	}

	fmt.Println("Delivered response!")

	// TODO: Add NUS and CAS SOAP to the case list.
	fmt.Println("[!] End of ECS Request.\n")
}
