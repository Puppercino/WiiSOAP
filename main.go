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
	"strconv"
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
	SerialNo  string   `xml:"Body>CheckRegistration>SerialNumber"`
}

// GetRegistrationInfo
type GRI struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>GetRegistrationInfo>Version"`
	DeviceId  string   `xml:"Body>GetRegistrationInfo>DeviceId"`
	MessageId string   `xml:"Body>GetRegistrationInfo>MessageId"`
	AccountId string   `xml:"Body>GetRegistrationInfo>AccountId"`
	Country   string   `xml:"Body>GetRegistrationInfo>Country"`
}

// Register
type REG struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>Register>Version"`
	DeviceId  string   `xml:"Body>Register>DeviceId"`
	MessageId string   `xml:"Body>Register>MessageId"`
	AccountId string   `xml:"Body>Register>AccountId"`
	Country   string   `xml:"Body>Register>Country"`
}

// Unregister
type UNR struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>Unregister>Version"`
	DeviceId  string   `xml:"Body>Unregister>DeviceId"`
	MessageId string   `xml:"Body>Unregister>MessageId"`
}

func main() {
	fmt.Println("Starting HTTP connection (Port 8000)...")
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Get a sexy new timestamp to use.
	timestampnano := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	fmt.Println(timestampnano)
	timestamp := timestampnano + "000"
	fmt.Println(timestamp)

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
	<TimeStamp>`+timestamp+`</TimeStamp>
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
	<TimeStamp>`+timestamp+`</TimeStamp>
	<ErrorCode>0</ErrorCode>
	<ServiceStandbyMode>false</ServiceStandbyMode>
	<OriginalSerialNumber>`+CR.SerialNo+`</OriginalSerialNumber>
	<DeviceStatus>R</DeviceStatus>
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
