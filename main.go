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
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	// Any given challenge must be 10 characters or less (challenge.length > 0xb)
	SharedChallenge = "NintyWhyPls"
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
	http.HandleFunc("/ecs/services/ECommerceSOAP", commonHandler)
	http.HandleFunc("/ias/services/IdentityAuthenticationSOAP", commonHandler)
	log.Fatal(http.ListenAndServe(CON.Address, nil))

	// From here on out, all special cool things should go into their respective handler function.
}

func commonHandler(w http.ResponseWriter, r *http.Request) {
	// Figure out the action to handle via header.
	service, action := parseAction(r.Header.Get("SOAPAction"))
	if service == "" || action == "" {
		printError(w, "WiiSOAP can't handle this. Try again later or actually use a Wii instead of a computer.")
		return
	}

	// Verify this is a service type we know.
	switch service {
	case "ecs":
	case "ias":
		break
	default:
		printError(w, "Unsupported service type...")
		return
	}

	fmt.Println("[!] Incoming " + strings.ToUpper(service) + " request - handling for " + action)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		printError(w, "Error reading request body...")
		return
	}

	// Tidy up parsed document for easier usage going forward.
	doc, err := normalise(service, action, strings.NewReader(string(body)))
	if err != nil {
		printError(w, "Error interpreting request body: "+err.Error())
		return
	}

	fmt.Println("Received:", string(body))

	// Insert the current action being performed.
	envelope := NewEnvelope(service, action)

	// Extract shared values from this request.
	err = envelope.ObtainCommon(doc)
	if err != nil {
		printError(w, "Error handling request body: "+err.Error())
		return
	}

	var successful bool
	var result string
	if service == "ias" {
		successful, result = iasHandler(envelope, doc)
	} else if service == "ecs" {
		successful, result = ecsHandler(envelope, doc)
	}

	if successful {
		// Write returned with proper Content-Type
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		w.Write([]byte(result))
	} else {
		printError(w, result)
	}

	fmt.Println("[!] End of " + strings.ToUpper(service) + " Request.\n")
}

func printError(w http.ResponseWriter, reason string) {
	http.Error(w, reason, http.StatusInternalServerError)
	fmt.Println("Failed to handle request: " + reason)
}
