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
	t.Log("filterNorthOfOffice returns buses currently north of 41.98 degrees latitude")

	mockData, err := os.Open("mocks/route22.xml")
	defer mockData.Close()
	route := mapToRoute(mockData)
	result := filterNorthOfOffice(route.Buses)
	expected := []Bus{
		{"4388", "42.01676344871521", "South Bound"},
		{"4350", "41.99443111134999", "South Bound"},
		{"4377", "41.98397789001465", "North Bound"},
		{"4345", "42.01595086566473", "North Bound"},
		{"4363", "42.01837830624338", "North Bound"},
		{"4339", "42.018760681152344", "North West Bound"},
	}

	if err != nil {
		t.Fatal(err)
	}

	for index, bus := range result {
		if bus != expected[index] {
			t.Error("Expected: \t", bus)
			t.Error("Received: \t", expected[index])
		}
	}

}

func TestCreateTable(t *testing.T) {
	t.Log("createTable() creates a data table of buses")
	mockData, err := os.Open("mocks/route22.xml")
	defer mockData.Close()

	if err != nil {
		t.Fatal(err)
	}

	routes := mapToRoute(mockData)
	filtered := filterNorthOfOffice(routes.Buses)
	result := createTable(filtered)
	expected :=
		"\n----------------------------------------------------" +
			"\nBUSES NORTH OF 41.98 LATITUDE" +
			"\n----------------------------------------------------" +
			"\nID\t Latitude\t\t Direction\n" +
			"4388\t 42.01676344871521\t South Bound\n" +
			"4350\t 41.99443111134999\t South Bound\n" +
			"4377\t 41.98397789001465\t North Bound\n" +
			"4345\t 42.01595086566473\t North Bound\n" +
			"4363\t 42.01837830624338\t North Bound\n" +
			"4339\t 42.018760681152344\t North West Bound\n" +
			"----------------------------------------------------"

	if result != expected {
		t.Error("Received:\t", result, "Expected:\t", expected)
	}
}
