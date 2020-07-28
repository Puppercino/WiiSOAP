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

func iasHandler(w http.ResponseWriter, r *http.Request) {
	// Figure out the action to handle via header.
	action := r.Header.Get("SOAPAction")
	action = parseAction(action, "ias")

	// Get a sexy new timestamp to use.
	timestampNano := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	timestamp := timestampNano + "000"

	fmt.Println("[!] Incoming IAS request.")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body...", http.StatusInternalServerError)
	}

	// The switch converts the HTTP Body of the request into a string. There is no need to convert the cases to byte format.
	switch action {
	// TODO: Make the case functions cleaner. (e.g. Should the response be a variable?)
	// TODO: Update the responses so that they query the SQL Database for the proper information (e.g. Device Code, Token, etc).

	case "CheckRegistration":
		fmt.Println("CR.")
		CR := CR{}
		if err = xml.Unmarshal(body, &CR); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "not good enough for me. ;3")
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(CR)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<soapenv:Body>
<CheckRegistrationResponse xmlns="urn:ias.wsapi.broadon.com">
	<Version>%s</Version>
	<DeviceId>%s</DeviceId>
	<MessageId>%s</MessageId>
	<TimeStamp>%s</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<OriginalSerialNumber>%s</OriginalSerialNumber>
	<DeviceStatus>R</DeviceStatus>
</CheckRegistrationResponse>
</soapenv:Body>
</soapenv:Envelope>`, CR.Version, CR.DeviceID, CR.DeviceID, timestamp, CR.SerialNo)

	case "GetRegistrationInfo":
		fmt.Println("GRI.")
		GRI := GRI{}
		if err = xml.Unmarshal(body, &GRI); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "how dirty. ;3")
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(GRI)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope   xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<GetRegistrationInfoResponse xmlns="urn:ias.wsapi.broadon.com">
			<Version>%s</Version>
			<DeviceId>%s</DeviceId>
			<MessageId>%s</MessageId>
			<TimeStamp>%s</TimeStamp>
			<ErrorCode>0</ErrorCode>
			<ServiceStandbyMode>false</ServiceStandbyMode>
			<AccountId>%s</AccountId>
			<DeviceToken>00000000</DeviceToken>
			<DeviceTokenExpired>false</DeviceTokenExpired>
			<Country>%s</Country>
			<ExtAccountId></ExtAccountId>
			<DeviceCode>0000000000000000</DeviceCode>
			<DeviceStatus>R</DeviceStatus>
			<Currency>POINTS</Currency>
		</GetRegistrationInfoResponse>
	</soapenv:Body>
</soapenv:Envelope>`, GRI.Version, GRI.DeviceID, GRI.MessageID, timestamp, GRI.AccountID, GRI.Country)

	case "Register":
		fmt.Println("REG.")
		REG := REG{}
		if err = xml.Unmarshal(body, &REG); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "disgustingly invalid. ;3")
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(REG)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
		xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<RegisterResponse xmlns="urn:ias.wsapi.broadon.com">
			<Version>%s</Version>
			<DeviceId>%s</DeviceId>
			<MessageId>%s</MessageId>
			<TimeStamp>%s</TimeStamp>
			<ErrorCode>0</ErrorCode>
			<ServiceStandbyMode>false</ServiceStandbyMode>
			<AccountId>%s</AccountId>
			<DeviceToken>00000000</DeviceToken>
			<Country>%s</Country>
			<ExtAccountId></ExtAccountId>
			<DeviceCode>00000000</DeviceCode>
		</RegisterResponse>
	</soapenv:Body>
</soapenv:Envelope>`, REG.Version, REG.DeviceID, REG.MessageID, timestamp, REG.AccountID, REG.Country)

	case "Unregister":
		fmt.Println("UNR.")
		UNR := UNR{}
		if err = xml.Unmarshal(body, &UNR); err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "how abnormal... ;3")
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Println(UNR)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<UnregisterResponse xmlns="urn:ias.wsapi.broadon.com">
			<Version>%s</Version>
			<DeviceId>%s</DeviceId>
			<MessageId>%s</MessageId>
			<TimeStamp>%s</TimeStamp>
			<ErrorCode>0</ErrorCode>
		<ServiceStandbyMode>false</ServiceStandbyMode>
		</UnregisterResponse>
	</soapenv:Body>
</soapenv:Envelope>`, UNR.Version, UNR.DeviceID, UNR.MessageID, timestamp)

	default:
		fmt.Fprintf(w, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer.")
		return
	}

	fmt.Println("Delivered response!")

	// TODO: Add NUS and CAS SOAP to the case list.
	fmt.Println("[!] End of IAS Request.\n")
}
