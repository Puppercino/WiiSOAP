//	Copyright (C) 2018  CornierKhan1
//
//	SOAP-GO-OSC is Open Shop Channel's SOAP Server Software, designed specifically to handle Wii Shop Channel SOAP.
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
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

/////////////////////
// SOAP STRUCTURES //
/////////////////////

type CDS struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>CheckDeviceStatus>Version"`
	DeviceId  string   `xml:"Body>CheckDeviceStatus>DeviceId"`
	MessageId string   `xml:"Body>CheckDeviceStatus>MessageId"`
}
type NETS struct {
}

func main() {
	fmt.Println("Starting HTTP connection (Port 8000)...")
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request! (Processing...)")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}

	// The following block of code here is a cluster of disgust.

	// CheckDeviceStatus
	if bytes.Contains(body, []byte("CheckDeviceStatus")) {
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
	<Version>2.0</Version>
	<DeviceId>`+CDS.DeviceId+`</DeviceId>
	<MessageId>`+CDS.MessageId+`</MessageId>
	<TimeStamp>00000000</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<Balance>
		<Amount>2018</Amount>
		<Currency>POINTS</Currency>
	</Balance>
	<ForceSyncTime>0</ForceSyncTime>
	<ExtTicketTime>00000000</ExtTicketTime>
	<SyncTime>00000000</SyncTime>
</CheckDeviceStatusResponse>
</soapenv:Body>
</soapenv:Envelope>`)
		fmt.Println("Delivered response!")

	} else {

		// NotifyETicketsSynced.
		if bytes.Contains(body, []byte("NotifyETicketsSynced")) {
			fmt.Println("NETS")
			NETS := NETS{}
			err = xml.Unmarshal([]byte(body), &NETS)
			if err != nil {
				fmt.Println("...or not. Bad or incomplete request. (End processing.)")
				fmt.Fprint(w, "You need to POST some SOAP from WSC if you wanna get some, honey. ;)")
				fmt.Printf("error: %v", err)
				return
			}
		}

	}
	fmt.Println("End of Request.")
}
