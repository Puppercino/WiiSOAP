package main

import (
	"net/http"
	"encoding/xml"
)

const (
	// Header is a generic XML header suitable for use with the output of Marshal.
	// This is not automatically added to any output of this package,
	// it is provided as a convenience.
	Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// The Check struct(ure) will attempt to retrieve all the namespace data.
// Assuming that namespaces that don't exist are given a "nil", the first result that
// isn't a "nil" will be used as the template response.

type Check struct {
	// ECommerce Namespaces
	CDS xml.Name `ecs:"CheckDeviceStatus"`
	LET	xml.Name `ecs:"ListETickets"`
	NETS xml.Name `ecs:"NotifyETicketsSynced"`
	PT xml.Name	`ecs:"PurchaseTitle"`

	// Identity Authentication Namespaces
}
func main() {


	http.ListenAndServe(":80", nil)
}