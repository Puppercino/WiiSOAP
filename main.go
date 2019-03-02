//	Copyright (C) 2018-2019  CornierKhan1
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
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// CheckError makes error handling not as ugly and inefficient.
func CheckError(e error) {
	if e != nil {
		log.Fatal("WiiSOAP forgot how to drive and suddenly crashed! Reason: ", e)
	}
}

func main() {
	// Initial Start.
	fmt.Println("WiiSOAP 0.2.5 Kawauso\nReading the Config...")

	// Check the Config.
	configfile, err := os.Open("./config.xml")
	CheckError(err)
	ioconfig, err := ioutil.ReadAll(configfile)
	CheckError(err)
	CON := Config{}
	err = xml.Unmarshal([]byte(ioconfig), &CON)
	fmt.Println(CON)
	CheckError(err)

	fmt.Println("Initializing core...")

	// Start SQL.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s)", CON.SQLUser, CON.SQLPass, CON.SQLPort, CON.SQLDB))
	CheckError(err)

	// Close SQL after everything else is done.
	defer db.Close()
	err = db.Ping()
	CheckError(err)

	// Start the HTTP server.
	fmt.Printf("Starting HTTP connection (%s)...\nNot using the usual port for HTTP? Be sure to use a proxy, otherwise the Wii can't connect!", CON.Address)
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe(CON.Address, nil))

	// From here on out, all special cool things should go into the handler function.
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Get a sexy new timestamp to use.
	timestampnano := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	timestamp := timestampnano + "000"

	fmt.Println("-=Incoming request!=-")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body...", http.StatusInternalServerError)
	}

	// The switch converts the HTTP Body of the request into a string. There is no need to convert the cases to byte format.
	switch string(body) {
	// TODO: Make the case functions cleaner. (e.g. Should the response be a variable?)
	// TODO: Update the responses so that they query the SQL Database for the proper information (e.g. Device Code, Token, etc).

	case "CheckDeviceStatus":
		fmt.Println("CDS.")
		CDS := CDS{}
		if err = xml.Unmarshal(body, &CDS); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "You need to POST some SOAP from WSC if you wanna get some, honey. ;)")
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
		fmt.Println("Delivered response!")

	case "NotifiedETicketsSynced":
		fmt.Println("NETS")
		NETS := NETS{}
		if err = xml.Unmarshal(body, &NETS); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "This is a disgusting request, but 20 dollars is 20 dollars. ;)")
			fmt.Printf("error: %v", err)
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
		fmt.Println("Delivered response!")

	case "ListETickets":
		fmt.Println("LET")
		LET := LET{}
		if err = xml.Unmarshal(body, &LET); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "This is a disgusting request, but 20 dollars is 20 dollars. ;)")
			fmt.Printf("error: %v", err)
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
		fmt.Println("Delivered response!")

	case "PurchaseTitle":
		fmt.Println("PT")
		PT := PT{}
		if err = xml.Unmarshal(body, &PT); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "if you wanna fun time, its gonna cost ya extra sweetie. ;)")
			fmt.Printf("Error: %s", err.Error())
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
		fmt.Println("Delivered response!")

	case "CheckRegistration":
		fmt.Println("CR.")
		CR := CR{}
		if err = xml.Unmarshal(body, &CR); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "not good enough for me. ;)")
			fmt.Printf("Error: %s", err.Error())
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
		fmt.Println("Delivered response!")

	case "GetRegistrationInfo":
		fmt.Println("GRI.")
		GRI := GRI{}
		if err = xml.Unmarshal(body, &GRI); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "how dirty. ;)")
			fmt.Printf("Error: %s", err.Error())
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
		fmt.Println("Delivered response!")

	case "Register":
		fmt.Println("REG.")
		REG := REG{}
		if err = xml.Unmarshal(body, &REG); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "disgustingly invalid. ;)")
			fmt.Printf("Error: %s", err.Error())
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
		fmt.Println("Delivered response!")

	case "Unregister":
		fmt.Println("UNR.")
		UNR := UNR{}
		if err = xml.Unmarshal(body, &UNR); os.IsExist(err) {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "how abnormal... ;)")
			fmt.Printf("Error: %s", err.Error())
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
		fmt.Println("Delivered response!")

	default:
		fmt.Fprintf(w, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer.")
	}

	// TODO: Add NUS and CAS SOAP to the case list.
	fmt.Println("-=End of Request.=-" + "\n")
}
