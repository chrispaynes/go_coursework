package main

import (
	// "fmt"
	"testing"
)

func TestReturn200StatusCode(t *testing.T) {
	t.Log("getData receives 200 reponse from Chicago Public Transit Bus Route 22")
	resp, err := connectToRoute(22)

	defer resp.Body.Close()

	if resp.StatusCode != 200 || err != nil {
		t.Errorf("Expected to connect to the site", err)
	}
}

// Test unmarshals XML data to slice
// Test prints northbound buses that currently north the office_latitude (41.98)
// Test osutput buses to a table
