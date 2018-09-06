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
	"time"
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
// The structures may seem repetitive and redundant, but blame WSC's inconsistent SOAP requests.

// CheckDeviceStatus
type CDS struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>CheckDeviceStatus>Version"`
	DeviceId  string   `xml:"Body>CheckDeviceStatus>DeviceId"`
	MessageId string   `xml:"Body>CheckDeviceStatus>MessageId"`
}

// NotifiedETicketsSynced
type NETS struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>NotifiedETicketsSynced>Version"`
	DeviceId  string   `xml:"Body>NotifiedETicketsSynced>DeviceId"`
	MessageId string   `xml:"Body>NotifiedETicketsSynced>MessageId"`
}

// ListETickets
type LET struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>ListETickets>Version"`
	DeviceId  string   `xml:"Body>ListETickets>DeviceId"`
	MessageId string   `xml:"Body>ListETickets>MessageId"`
}

// PurchaseTitle
type PT struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>PurchaseTitle>Version"`
	DeviceId  string   `xml:"Body>PurchaseTitle>DeviceId"`
	MessageId string   `xml:"Body>PurchaseTitle>MessageId"`
}

// CheckRegistration
type CR struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>CheckRegistration>Version"`
	DeviceId  string   `xml:"Body>CheckRegistration>DeviceId"`
	MessageId string   `xml:"Body>CheckRegistration>MessageId"`
}

// GetRegistrationInfo
type GRI struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>GetRegistrationInfo>Version"`
	DeviceId  string   `xml:"Body>GetRegistrationInfo>DeviceId"`
	MessageId string   `xml:"Body>GetRegistrationInfo>MessageId"`
}

// Register
type REG struct {
}

// Unregister
type UNR struct {
}

func main() {
	fmt.Println("Starting HTTP connection (Port 8000)...")
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-=Incoming request!=-")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body...",
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
	<TimeStamp>00000000</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
</NotifyETicketsSyncedResponse>
</soapenv:Body>
</soapenv:Envelope>`)
			fmt.Println("Delivered response!")

		} else {

			// ListETickets
			if bytes.Contains(body, []byte("ListETickets")) {
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
	<TimeStamp>00000000</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<ForceSyncTime>0</ForceSyncTime>
	<ExtTicketTime>0</ExtTicketTime>
	<SyncTime>00000000</SyncTime>
</ListETicketsResponse>
</soapenv:Body>
</soapenv:Envelope>`)
				fmt.Println("Delivered response!")

			} else {
				if bytes.Contains(body, []byte("PurchaseTitle")) {
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
	<TimeStamp>00000000</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<Balance>
		<Amount>2018</Amount>
		<Currency>POINTS</Currency>
	</Balance>
	<Transactions>
		<TransactionId>00000000</TransactionId>
		<Date>`+time.RFC3339+`</Date>
		<Type>PURCHGAME</Type>
	</Transactions>
	<SyncTime>00000000</SyncTime>
	<ETickets>00000000</ETickets>
	<Certs>00000000</Certs>
	<Certs>00000000</Certs>
	<TitleId>00000000</TitleId>
</PurchaseTitleResponse>
</soapenv:Body>
</soapenv:Envelope>`)

				} else {

					// IDENTITY AUTHENTICATION SOAP

					if bytes.Contains(body, []byte("CheckRegistration")) {
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
	<TimeStamp>00000000</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<OriginalSerialNumber></OriginalSerialNumber>
	<DeviceStatus>$DeviceStatus</DeviceStatus>
</CheckRegistrationResponse>
</soapenv:Body>
</soapenv:Envelope>`)

					} else {

						if bytes.Contains(body, []byte("GetRegistrationInfo")) {
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
							fmt.Fprintf(w, ``)

						} else {

							if bytes.Contains(body, []byte("Register")) {
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
								fmt.Fprintf(w, ``)

							} else {

								if bytes.Contains(body, []byte("Unregister")) {
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
									fmt.Fprintf(w, ``)

								} else {
									fmt.Println("Nothing sent?")
									fmt.Fprintf(w, "the bathtub is empty, it needs SOAP. ;)")
								}
							}
						}
					}
				}
			}
		}

	}
	fmt.Println("-=End of Request.=-")
}
