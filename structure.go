package main

import (
	"encoding/xml"
)

//	Copyright (C) 2018-2019  CornierKhan1
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

/////////////////////
// SOAP STRUCTURES //
/////////////////////
// The structures may seem repetitive and redundant, but blame WSC's inconsistent SOAP requests.

type Config struct {
	XMLName xml.Name `xml:"Config"`
	SQLUser string   `xml:"SQLUser"`
	SQLPass string   `xml:"SQLPass"`
	SQLPort string   `xml:"SQLPort"`
	Port    string   `xml:"Port"`
	SQLDB   string   `xml:"SQLDB"`
}

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
