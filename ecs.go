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
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
	  	xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
	  	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<CheckDeviceStatusResponse xmlns="urn:ecs.wsapi.broadon.com">
			<Version>%s</Version>
			<DeviceId>%s</DeviceId>
			<MessageId>%s</MessageId>
			<TimeStamp>%s</TimeStamp>
			<ErrorCode>0</ErrorCode>
			<ServiceStandbyMode>false</ServiceStandbyMode>
			<Balance>
				<Amount>2018</Amount>
				<Currency>POINTS</Currency>
			</Balance>
			<ForceSyncTime>0</ForceSyncTime>
			<ExtTicketTime>%s</ExtTicketTime>
			<SyncTime>%s</SyncTime>
		</CheckDeviceStatusResponse>
	</soapenv:Body>
</soapenv:Envelope>`, CDS.Version, CDS.DeviceID, CDS.MessageID, timestamp, timestamp, timestamp)

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
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
			xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
			xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<soapenv:Body>
			<NotifyETicketsSyncedResponse xmlns="urn:ecs.wsapi.broadon.com">
				<Version>%s</Version>
				<DeviceId>%s</DeviceId>
				<MessageId>%s</MessageId>
				<TimeStamp>%s</TimeStamp>
				<ErrorCode>0</ErrorCode>
				<ServiceStandbyMode>false</ServiceStandbyMode>
			</NotifyETicketsSyncedResponse>
	</soapenv:Body>
</soapenv:Envelope>`, NETS.Version, NETS.DeviceID, NETS.MessageID, timestamp)

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
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<ListETicketsResponse xmlns="urn:ecs.wsapi.broadon.com">
			<Version>%s</Version>
			<DeviceId>%s</DeviceId>
			<MessageId>%s</MessageId>
			<TimeStamp>%s</TimeStamp>
			<ErrorCode>0</ErrorCode>
			<ServiceStandbyMode>false</ServiceStandbyMode>
			<ForceSyncTime>0</ForceSyncTime>
			<ExtTicketTime>%s</ExtTicketTime>
			<SyncTime>%s</SyncTime>
		</ListETicketsResponse>
	</soapenv:Body>
</soapenv:Envelope>`, LET.Version, LET.DeviceID, LET.MessageID, timestamp, timestamp, timestamp)

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
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope	xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<PurchaseTitleResponse xmlns="urn:ecs.wsapi.broadon.com">
			<Version>%s</Version>
			<DeviceId>%s</DeviceId>
			<MessageId>%s</MessageId>
			<TimeStamp>%s</TimeStamp>
			<ErrorCode>0</ErrorCode>
			<ServiceStandbyMode>false</ServiceStandbyMode>
			<Balance>
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
			<TitleId>00000000</TitleId>
		</PurchaseTitleResponse>
	</soapenv:Body>
</soapenv:Envelope>`, PT.Version, PT.DeviceID, PT.MessageID, timestamp, timestamp, timestamp)

	default:
		fmt.Fprintf(w, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer.")
		return
	}

	fmt.Println("Delivered response!")

	// TODO: Add NUS and CAS SOAP to the case list.
	fmt.Println("[!] End of ECS Request.\n")
}
