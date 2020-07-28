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

import "encoding/xml"

/////////////////////
// SOAP STRUCTURES //
/////////////////////
// The structures may seem repetitive and redundant, but blame WSC's inconsistent SOAP requests.

// Config - WiiSOAP Configuration data.
type Config struct {
	XMLName xml.Name `xml:"Config"`

	Address string `xml:"Address"`

	SQLAddress string `xml:"SQLAddress"`
	SQLUser    string `xml:"SQLUser"`
	SQLPass    string `xml:"SQLPass"`
	SQLDB      string `xml:"SQLDB"`
}

// CDS - CheckDeviceStatus
type CDS struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>CheckDeviceStatus>Version"`
	DeviceID  string   `xml:"Body>CheckDeviceStatus>DeviceId"`
	MessageID string   `xml:"Body>CheckDeviceStatus>MessageId"`
}

// NETS - NotifiedETicketsSynced
type NETS struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>NotifiedETicketsSynced>Version"`
	DeviceID  string   `xml:"Body>NotifiedETicketsSynced>DeviceId"`
	MessageID string   `xml:"Body>NotifiedETicketsSynced>MessageId"`
}

// LET - ListETickets
type LET struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>ListETickets>Version"`
	DeviceID  string   `xml:"Body>ListETickets>DeviceId"`
	MessageID string   `xml:"Body>ListETickets>MessageId"`
}

// PT - PurchaseTitle
type PT struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>PurchaseTitle>Version"`
	DeviceID  string   `xml:"Body>PurchaseTitle>DeviceId"`
	MessageID string   `xml:"Body>PurchaseTitle>MessageId"`
}

// CR - CheckRegistration
type CR struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>CheckRegistration>Version"`
	DeviceID  string   `xml:"Body>CheckRegistration>DeviceId"`
	MessageID string   `xml:"Body>CheckRegistration>MessageId"`
	SerialNo  string   `xml:"Body>CheckRegistration>SerialNumber"`
}

// GRI - GetRegistrationInfo
type GRI struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>GetRegistrationInfo>Version"`
	DeviceID  string   `xml:"Body>GetRegistrationInfo>DeviceId"`
	MessageID string   `xml:"Body>GetRegistrationInfo>MessageId"`
	AccountID string   `xml:"Body>GetRegistrationInfo>AccountId"`
	Country   string   `xml:"Body>GetRegistrationInfo>Country"`
}

// REG - Register
type REG struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>Register>Version"`
	DeviceID  string   `xml:"Body>Register>DeviceId"`
	MessageID string   `xml:"Body>Register>MessageId"`
	AccountID string   `xml:"Body>Register>AccountId"`
	Country   string   `xml:"Body>Register>Country"`
}

// UNR - Unregister
type UNR struct {
	XMLName   xml.Name `xml:"Envelope"`
	Version   string   `xml:"Body>Unregister>Version"`
	DeviceID  string   `xml:"Body>Unregister>DeviceId"`
	MessageID string   `xml:"Body>Unregister>MessageId"`
}
