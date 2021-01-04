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

// Envelope represents the root element of any response, soapenv:Envelope.
type Envelope struct {
	XMLName string `xml:"soapenv:Envelope"`
	SOAPEnv string `xml:"xmlns:soapenv,attr"`
	XSD     string `xml:"xmlns:xsd,attr"`
	XSI     string `xml:"xmlns:xsi,attr"`

	// Represents a soapenv:Body within.
	Body Body

	// Used for internal state tracking.
	action string
}

// Body represents the nested soapenv:Body element as a child on the root element,
// containing the response intended for the action being handled.
type Body struct {
	XMLName string `xml:"soapenv:Body"`

	// Represents the actual response inside
	Response Response
}

// Response describes the inner response format, along with common fields across requests.
type Response struct {
	XMLName xml.Name
	XMLNS   string `xml:"xmlns,attr"`

	// These common fields are persistent across all requests.
	Version            string `xml:"Version"`
	DeviceId           string `xml:"DeviceId"`
	MessageId          string `xml:"MessageId"`
	TimeStamp          string `xml:"TimeStamp"`
	ErrorCode          int
	ServiceStandbyMode bool `xml:"ServiceStandbyMode"`

	// Allows for <name>[dynamic content]</name> situations.
	CustomFields []interface{}
}

// KVField represents an individual node in form of <XMLName>Contents</XMLName>.
type KVField struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// Balance represents a common XML structure.
type Balance struct {
	XMLName  xml.Name `xml:"Balance"`
	Amount   int      `xml:"Amount"`
	Currency string   `xml"Currency"`
}

// Transactions represents a common XML structure.
type Transactions struct {
	XMLName       xml.Name `xml:"Transactions"`
	TransactionId string   `xml:"TransactionId"`
	Date          string   `xml:"Date"`
	Type          string   `xml:"Type"`
}
