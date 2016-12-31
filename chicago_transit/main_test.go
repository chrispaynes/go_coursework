package main

import "testing"

func TestXMLContentResponse(t *testing.T) {
	t.Log("connectToRoute() receives XML Response from Chicago Transit Authority API\n")
	resp, err := connectToRoute(22)
	resultContentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/xml;charset=UTF-8"

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resultContentType != expectedContentType {
		t.Error("Received Content-Type:\t", resultContentType, "Expected:\t", expectedContentType)
		t.Fatal("Did not receive an XML response from CTA API")
	}
}

// Test unmarshals XML data to slice
// Test prints northbound buses that currently north the office_latitude (41.98)
// Test osutput buses to a table
