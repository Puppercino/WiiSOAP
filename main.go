//	Copyright (C) 2018  CornierKhan1
//
//	WiiSOAP is Open Shop Channel's SOAP Server Software, designed specifically to handle Wii Shop Channel SOAP.
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
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
		log.Fatal("There has been an error while performing this action. Reason: ", e)
	}
}

func main() {

	fmt.Println("WiiSOAP 0.2.2 Kawauso")
	fmt.Println("Reading the Config...")
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
	db, err := sql.Open("mysql",
		CON.SQLUser+":"+CON.SQLPass+"@tcp(127.0.0.1:3306)/hello")
	CheckError(err)
	defer db.Close()

	fmt.Println("Starting HTTP connection (" + CON.Port + ")...")
	fmt.Println("NOTICE: The SOAP Server runs under a port that doesn't work with WSC naturally.")
	fmt.Println("You can either use proxying from Nginx (recommended) or another web server software, or edit this script to use port 80.")
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe(CON.Port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Get a sexy new timestamp to use.
	timestampnano := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	timestamp := timestampnano + "000"

	fmt.Println("-=Incoming request!=-")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body...",
			http.StatusInternalServerError)
	}

	// The requests are in byte format, so please use []byte("") when adding a new case.
	switch body {
	// TODO: Make the case functions cleaner. (e.g. The should the response be a variable?)
	case []byte("CheckDeviceStatus"):

		fmt.Println("CDS.")
		CDS := CDS{}
		err = xml.Unmarshal([]byte(body), &CDS)
		if err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "You need to POST some SOAP from WSC if you wanna get some, honey. ;)")
			fmt.Printf("error: %v", err)
			return
		}
		fmt.Println(CDS)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprintf(w, `
<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
				  xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
				  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<soapenv:Body>
<CheckDeviceStatusResponse xmlns="urn:ecs.wsapi.broadon.com">
	<Version>`+CDS.Version+`</Version>
	<DeviceId>`+CDS.DeviceId+`</DeviceId>
	<MessageId>`+CDS.MessageId+`</MessageId>
	<TimeStamp>`+timestamp+`</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<Balance>
		<Amount>2018</Amount>
		<Currency>POINTS</Currency>
	</Balance>
	<ForceSyncTime>0</ForceSyncTime>
	<ExtTicketTime>`+timestamp+`</ExtTicketTime>
	<SyncTime>`+timestamp+`</SyncTime>
</CheckDeviceStatusResponse>
</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	case []byte("NotifiedETicketsSynced"):

		fmt.Println("NETS")
		NETS := NETS{}
		err = xml.Unmarshal([]byte(body), &NETS)
		if err != nil {
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
	<Version>`+NETS.Version+`</Version>
	<DeviceId>`+NETS.DeviceId+`</DeviceId>
	<MessageId>`+NETS.MessageId+`</MessageId>
	<TimeStamp>`+timestamp+`</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
</NotifyETicketsSyncedResponse>
</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	case []byte("ListETickets"):

		fmt.Println("LET")
		LET := LET{}
		err = xml.Unmarshal([]byte(body), &LET)
		if err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "This is a disgusting request, but 20 dollars is 20 dollars. ;)")
			fmt.Printf("error: %v", err)
			return
		}
		fmt.Println(LET)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope   xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
					xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
					xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<soapenv:Body>
<ListETicketsResponse xmlns="urn:ecs.wsapi.broadon.com">
	<Version>`+LET.Version+`</Version>
	<DeviceId>`+LET.DeviceId+`</DeviceId>
	<MessageId>`+LET.MessageId+`</MessageId>
	<TimeStamp>`+timestamp+`</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<ForceSyncTime>0</ForceSyncTime>
	<ExtTicketTime>`+timestamp+`</ExtTicketTime>
	<SyncTime>`+timestamp+`</SyncTime>
</ListETicketsResponse>
</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	case []byte("PurchaseTitle"):

		fmt.Println("PT")
		PT := PT{}
		err = xml.Unmarshal([]byte(body), &PT)
		if err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "if you wanna fun time, its gonna cost ya extra sweetie. ;)")
			fmt.Printf("error: %v", err)
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
	<Version>`+PT.Version+`</Version>
	<DeviceId>`+PT.DeviceId+`</DeviceId>
	<MessageId>`+PT.MessageId+`</MessageId>
	<TimeStamp>`+timestamp+`</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<Balance>
		<Amount>2018</Amount>
		<Currency>POINTS</Currency>
	</Balance>
	<Transactions>
		<TransactionId>00000000</TransactionId>
		<Date>`+timestamp+`</Date>
		<Type>PURCHGAME</Type>
	</Transactions>
	<SyncTime>`+timestamp+`</SyncTime>
	<ETickets>00000000</ETickets>
	<Certs>00000000</Certs>
	<Certs>00000000</Certs>
	<TitleId>00000000</TitleId>
</PurchaseTitleResponse>
</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	case []byte("CheckRegistration"):

		fmt.Println("CR.")
		CR := CR{}
		err = xml.Unmarshal([]byte(body), &CR)
		if err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "not good enough for me. ;)")
			fmt.Printf("error: %v", err)
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
	<Version>`+CR.Version+`</Version>
	<DeviceId>`+CR.DeviceId+`</DeviceId>
	<MessageId>`+CR.DeviceId+`</MessageId>
	<TimeStamp>`+timestamp+`</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<OriginalSerialNumber>`+CR.SerialNo+`</OriginalSerialNumber>
	<DeviceStatus>R</DeviceStatus>
</CheckRegistrationResponse>
</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	case []byte("GetRegistrationInfo"):

		fmt.Println("GRI.")
		GRI := GRI{}
		err = xml.Unmarshal([]byte(body), &GRI)
		if err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "how dirty. ;)")
			fmt.Printf("error: %v", err)
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
	<Version>`+GRI.Version+`</Version>
	<DeviceId>`+GRI.DeviceId+`</DeviceId>
	<MessageId>`+GRI.MessageId+`</MessageId>
	<TimeStamp>`+timestamp+`</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<AccountId>`+GRI.AccountId+`</AccountId>
	<DeviceToken>00000000</DeviceToken>
	<DeviceTokenExpired>false</DeviceTokenExpired>
	<Country>`+GRI.Country+`</Country>
	<ExtAccountId></ExtAccountId>
	<DeviceCode>0000000000000000</DeviceCode>
	<DeviceStatus>R</DeviceStatus>
	<Currency>POINTS</Currency>
</GetRegistrationInfoResponse>
</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	case []byte("Register"):

		fmt.Println("REG.")
		REG := REG{}
		err = xml.Unmarshal([]byte(body), &REG)
		if err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "disgustingly invalid. ;)")
			fmt.Printf("error: %v", err)
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
			<Version>`+REG.Version+`</Version>
			<DeviceId>`+REG.DeviceId+`</DeviceId>
			<MessageId>`+REG.MessageId+`</MessageId>
			<TimeStamp>`+timestamp+`</TimeStamp>
			<ErrorCode>0</ErrorCode>
			<ServiceStandbyMode>false</ServiceStandbyMode>
			<AccountId>`+REG.AccountId+`</AccountId>
			<DeviceToken>00000000</DeviceToken>
			<Country>`+REG.Country+`</Country>
			<ExtAccountId></ExtAccountId>
			<DeviceCode>00000000</DeviceCode>
		</RegisterResponse>
	</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	case []byte("Unregister"):

		fmt.Println("UNR.")
		UNR := UNR{}
		err = xml.Unmarshal([]byte(body), &UNR)
		if err != nil {
			fmt.Println("...or not. Bad or incomplete request. (End processing.)")
			fmt.Fprint(w, "how abnormal... ;)")
			fmt.Printf("error: %v", err)
			return
		}
		fmt.Println(UNR)
		fmt.Println("The request is valid! Responding...")
		fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soapenv:Body>
		<UnregisterResponse xmlns="urn:ias.wsapi.broadon.com">
			<Version>`+UNR.Version+`</Version>
			<DeviceId>`+UNR.DeviceId+`</DeviceId>
			<MessageId>`+UNR.MessageId+`</MessageId>
			<TimeStamp>`+timestamp+`</TimeStamp>
			<ErrorCode>0</ErrorCode>
		<ServiceStandbyMode>false</ServiceStandbyMode>
		</UnregisterResponse>
	</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	}

	// TODO: Add NUS and CAS SOAP to the case list.
	fmt.Println("-=End of Request.=-" + "\n")
}
