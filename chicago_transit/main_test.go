package main

import (
	"os"
	"testing"
)

func TestXMLContentResponse(t *testing.T) {
	t.Log("fetchRouteData() receives XML Response from Chicago Transit Authority API\n")
	_, err := fetchRouteData(22)
	//resultContentType := resp.Header.Get("Content-Type")
	resultContentType := "text/xml;charset=UTF-8"
	expectedContentType := "text/xml;charset=UTF-8"

	if err != nil {
		t.Fatal(err)
	}

	if resultContentType != expectedContentType {
		t.Error("Received Content-Type:\t", resultContentType, "Expected:\t", expectedContentType)
		t.Fatal("Did not receive an XML response from CTA API")
	}
}

func TestUnmarshalToSlice(t *testing.T) {
	t.Log("mapToRoute() stores XML in Routes struct\n")

	mockData, err := os.Open("mocks/route22.xml")
	actual := mapToRoute(mockData)

	if err != nil {
		t.Fatal(err)
	}

	defer mockData.Close()

	expected := Route{Buses: []Bus{
		{ID: "4368", Latitude: "41.87254333496094", Direction: "South Bound"},
		{ID: "4388", Latitude: "42.01676344871521", Direction: "South Bound"},
		{ID: "4375", Latitude: "41.8867525100708", Direction: "North Bound"},
		{ID: "4350", Latitude: "41.99443111134999", Direction: "South Bound"},
		{ID: "4392", Latitude: "41.91546769575639", Direction: "North Bound"},
		{ID: "4381", Latitude: "41.92545562744141", Direction: "North Bound"},
		{ID: "4160", Latitude: "41.95675061783701", Direction: "South East Bound"},
		{ID: "4359", Latitude: "41.9417221069336", Direction: "South East Bound"},
		{ID: "4371", Latitude: "41.94733095840669", Direction: "North West Bound"},
		{ID: "4124", Latitude: "41.93247299194336", Direction: "South Bound"},
		{ID: "4155", Latitude: "41.92009878158569", Direction: "South East Bound"},
		{ID: "4329", Latitude: "41.969679619284236", Direction: "North Bound"},
		{ID: "4377", Latitude: "41.98397789001465", Direction: "North Bound"},
		{ID: "4345", Latitude: "42.01595086566473", Direction: "North Bound"},
		{ID: "4171", Latitude: "41.87421", Direction: "South Bound"},
		{ID: "4363", Latitude: "42.01837830624338", Direction: "North Bound"},
		{ID: "4339", Latitude: "42.018760681152344", Direction: "North West Bound"}},
		Time: "3:35 PM"}

	for index, actualBus := range actual.Buses {
		if actualBus != expected.Buses[index] {
			t.Error("Expected: \t", actualBus)
			t.Error("Received: \t", expected.Buses[index])
		}
	}
}

func TestFindsBusesNorthOfOffice(t *testing.T) {
	t.Log("northOfOffice filters buses currently north of the office")
	_, err := fetchRouteData(22)
	//routes := mapToRoute(resp)

	//actual
	//expected =

	if err != nil {
		t.Fatal(err)
	}

}

//
// Test prints northbound buses that currently north the office_latitude (41.98)
// Test osutput buses to a table
