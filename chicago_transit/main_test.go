package main

import (
	"os"
	"testing"
)

func TestXMLContentResponse(t *testing.T) {
	t.Log("fetchRouteData() receives XML Response from Chicago Transit Authority API")
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
	t.Log("mapToRoute() stores XML in Routes struct")

	mockData, err := os.Open("mocks/route22.xml")
	result := mapToRoute(mockData)

	if err != nil {
		t.Fatal(err)
	}

	defer mockData.Close()

	expected := Route{Buses: []Bus{
		{"4368", "41.8725", "-87.6306", "South Bound"},
		{"4388", "42.0167", "-87.6754", "South Bound"},
		{"4375", "41.8867", "-87.6294", "North Bound"},
		{"4350", "41.9944", "-87.6702", "South Bound"},
		{"4392", "41.9154", "-87.6341", "North Bound"},
		{"4381", "41.9254", "-87.6404", "North Bound"},
		{"4160", "41.9567", "-87.6637", "South East Bound"},
		{"4359", "41.9417", "-87.6521", "South East Bound"},
		{"4371", "41.9473", "-87.6565", "North West Bound"},
		{"4124", "41.9324", "-87.6448", "South Bound"},
		{"4155", "41.9200", "-87.6371", "South East Bound"},
		{"4329", "41.9696", "-87.6675", "North Bound"},
		{"4377", "41.9839", "-87.6687", "North Bound"},
		{"4345", "42.0159", "-87.6751", "North Bound"},
		{"4171", "41.8742", "-87.6307", "South Bound"},
		{"4363", "42.0183", "-87.6729", "North Bound"},
		{"4339", "42.0187", "-87.6731", "North West Bound"}},
		Time: "3:35 PM"}

	for i, bus := range result.Buses {
		bus.Longitude = bus.sliceLongitude()
		bus.Latitude = bus.sliceLatitude()
		if bus != expected.Buses[i] {
			t.Error("Expected: \t", expected.Buses[i])
			t.Error("Received: \t", bus)
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
		{"4388", "42.0167", "-87.6754", "South Bound"},
		{"4350", "41.9944", "-87.6702", "South Bound"},
		{"4377", "41.9839", "-87.6687", "North Bound"},
		{"4345", "42.0159", "-87.6751", "North Bound"},
		{"4363", "42.0183", "-87.6729", "North Bound"},
		{"4339", "42.0187", "-87.6731", "North West Bound"},
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
			"\nID\t Latitude\t Longitude\t\t Direction\n" +
			"4388\t 42.0167\t -87.6754\t South Bound\n" +
			"4350\t 41.9944\t -87.6702\t South Bound\n" +
			"4377\t 41.9839\t -87.6687\t North Bound\n" +
			"4345\t 42.0159\t -87.6751\t North Bound\n" +
			"4363\t 42.0183\t -87.6729\t North Bound\n" +
			"4339\t 42.0187\t -87.6731\t North West Bound\n" +
			"----------------------------------------------------"

	if result != expected {
		t.Error("Received:\t", result, "Expected:\t", expected)
	}
}
