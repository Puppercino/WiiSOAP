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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// The Check struct(ure) will attempt to retrieve all the namespace data.
// Assuming that namespaces that don't exist are given a "nil", the first result that  isn't a "nil" will be used as the template response.
const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

type Check struct {

	// SOAP envelope doesn't matter to OSC. We'll only need the BODY.
	SOAP xml.Name `xml:"SOAP-ENV:Body"`

	// ECommerce Namespaces
	CDS  string `xml:"CheckDeviceStatus>Version"`
	LET  string `xml:"ListETickets>Version"`
	NETS string `xml:"NotifyETicketsSynced>Version"`
	PT   string `xml:"PurchaseTitle>Version"`

	// Identity Authentication Namespaces
	CR  string `ias:"CheckRegistration>Version"`
	GRI string `ias:"GetRegistrationInfo>Version"`
	REG string `ias:"Register>Version"`
	UNR string `ias:"Unregister>Version"`
}

func main() {
	fmt.Println("Starting HTTP connection (Port 8000)...")
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ChRes := Check{
		SOAP: xml.Name{},
		CDS:  "",
		LET:  "",
		NETS: "",
		PT:   "",
		CR:   "",
		GRI:  "",
		REG:  "",
		UNR:  "",
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
	}

	err = xml.Unmarshal([]byte(body), &ChRes)
	if err != nil {
		fmt.Println(ChRes)
		fmt.Fprint(w, "You need to POST some SOAP from WSC if you wanna get some, honey. ;)")
		fmt.Printf("error: %v", err)
		return
	}

}
