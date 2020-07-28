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
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// checkError makes error handling not as ugly and inefficient.
func checkError(err error) {
	if err != nil {
		log.Fatalf("WiiSOAP forgot how to drive and suddenly crashed! Reason: %v\n", err)
	}
}

func main() {
	// Initial Start.
	fmt.Println("WiiSOAP 0.2.6 Kawauso\n[i] Reading the Config...")

	// Check the Config.
	ioconfig, err := ioutil.ReadFile("./config.xml")
	checkError(err)
	CON := Config{}
	err = xml.Unmarshal(ioconfig, &CON)
	checkError(err)

	fmt.Println("[i] Initializing core...")

	// Start SQL.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", CON.SQLUser, CON.SQLPass, CON.SQLAddress, CON.SQLDB))
	checkError(err)

	// Close SQL after everything else is done.
	defer db.Close()
	err = db.Ping()
	checkError(err)

	// Start the HTTP server.
	fmt.Printf("Starting HTTP connection (%s)...\nNot using the usual port for HTTP?\nBe sure to use a proxy, otherwise the Wii can't connect!\n", CON.Address)

	// These following endpoints don't have to match what the official WSC have.
	// However, semantically, it feels proper.
	http.HandleFunc("/ecs/services/ECommerceSOAP", ecsHandler)              // For ECS operations
	http.HandleFunc("/ias/services/IdentityAuthenticationSOAP", iasHandler) // For IAS operations
	log.Fatal(http.ListenAndServe(CON.Address, nil))

	// From here on out, all special cool things should go into their respective handler function.
}
