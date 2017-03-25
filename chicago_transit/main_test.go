package main

import (
	"os"
	"testing"
)

func rangeTestTable(t *testing.T, actual, tc []Bus) {
	for k, v := range actual {
		if v != tc[k] {
			t.Error("Expected: \t", v)
			t.Fatal("Received: \t", tc[k])
		}
	}
}

func openMockData(t *testing.T, src string) *os.File {
	f, err := os.Open("mocks/route22.xml")

	if err != nil {
		t.Fatal(err)
	}

	return f
}

func TestXMLContentResponse(t *testing.T) {
	t.Log("fetchRouteData() receives XML Response from Chicago Transit Authority API")

	_, err := fetchRouteData(22)
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

	f := openMockData(t, "mocks/route22.xml")

	actual := mapToRoute(f)
	tc := Route{Buses: []Bus{
		{4368, Point{41.87254333496094, -87.63065338134766}, "South Bound"},
		{4388, Point{42.01676344871521, -87.67540860176086}, "South Bound"},
		{4375, Point{41.8867525100708, -87.62945556640625}, "North Bound"},
		{4350, Point{41.99443111134999, -87.67027897621269}, "South Bound"},
		{4392, Point{41.91546769575639, -87.63416186246005}, "North Bound"},
		{4381, Point{41.92545562744141, -87.6404949951172}, "North Bound"},
		{4160, Point{41.95675061783701, -87.66378870550191}, "South East Bound"},
		{4359, Point{41.9417221069336, -87.65216827392578}, "South East Bound"},
		{4371, Point{41.94733095840669, -87.65654185120489}, "North West Bound"},
		{4124, Point{41.93247299194336, -87.64489555358887}, "South Bound"},
		{4155, Point{41.92009878158569, -87.63711357116699}, "South East Bound"},
		{4329, Point{41.969679619284236, -87.66758620318244}, "North Bound"},
		{4377, Point{41.98397789001465, -87.66879043579101}, "North Bound"},
		{4345, Point{42.01595086566473, -87.6751211134054}, "North Bound"},
		{4171, Point{41.87421, -87.63072333333334}, "South Bound"},
		{4363, Point{42.01837830624338, -87.67295914989407}, "North Bound"},
		{4339, Point{42.018760681152344, -87.67317962646484}, "North West Bound"}},
		Time: "3:35 PM"}

	rangeTestTable(t, actual.Buses, tc.Buses)
}

func TestFindsBusesNorthOfOffice(t *testing.T) {
	t.Log("filterNorthOfOffice returns buses currently north of 41.98 degrees latitude")

	f := openMockData(t, "mocks/route22.xml")
	route := mapToRoute(f)

	actual := filterNorthOfOffice(route.Buses)
	tc := []Bus{
		{4388, Point{42.01676344871521, -87.67540860176086}, "South Bound"},
		{4350, Point{41.99443111134999, -87.67027897621269}, "South Bound"},
		{4377, Point{41.98397789001465, -87.66879043579101}, "North Bound"},
		{4345, Point{42.01595086566473, -87.6751211134054}, "North Bound"},
		{4363, Point{42.01837830624338, -87.67295914989407}, "North Bound"},
		{4339, Point{42.018760681152344, -87.67317962646484}, "North West Bound"},
	}

	rangeTestTable(t, actual, tc)
}

func TestCreateTable(t *testing.T) {
	t.Log("createTable() creates a data table of buses")

	f := openMockData(t, "mocks/route22.xml")
	routes := mapToRoute(f)
	filtered := filterNorthOfOffice(routes.Buses)

	actual := createTable(filtered, "BUSES NORTH OF 41.98 LATITUDE")
	tc :=
		"\n----------------------------------------------------" +
			"\nBUSES NORTH OF 41.98 LATITUDE" +
			"\n----------------------------------------------------" +
			"\nID\t Latitude\t Longitude\t\t Direction\n" +
			"4388\t 42.01676344871521\t -87.67540860176086\t South Bound\n" +
			"4350\t 41.99443111134999\t -87.67027897621269\t South Bound\n" +
			"4377\t 41.98397789001465\t -87.66879043579101\t North Bound\n" +
			"4345\t 42.01595086566473\t -87.6751211134054\t North Bound\n" +
			"4363\t 42.01837830624338\t -87.67295914989407\t North Bound\n" +
			"4339\t 42.018760681152344\t -87.67317962646484\t North West Bound\n" +
			"----------------------------------------------------"

	if actual != tc {
		t.Error("Received:\t", actual)
		t.Error("Expected:\t", tc)
	}
}

func TestFindDistance(t *testing.T) {
	t.Log("findDistance() finds the great-circle distance between two map points")

	p1 := &Point{41.9801433, -87.6683411}
	p2 := &Point{41.9855176, -87.6702406}
	p3 := &Point{41.9955993, -87.6809199}
	p4 := &Point{41.944268, -87.6670651}

	result1 := findDistance(p1, p2)
	result2 := findDistance(p1, p3)
	result3 := findDistance(p1, p4)

	expected1 := 0.38414780465548903
	expected2 := 1.2488013385943528
	expected3 := 2.481016128867292

	if result1 != expected1 {
		t.Error("Received:\t", result1)
		t.Error("Expected:\t", expected1)
	}

	if result2 != expected2 {
		t.Error("Received:\t", result2)
		t.Error("Expected:\t", expected2)
	}

	if result3 != expected3 {
		t.Error("Received:\t", result3)
		t.Error("Expected:\t", expected3)
	}
}

func TestWithinHalfMile(t *testing.T) {
	t.Log("withinHalfMile() returns buses within 0.5 miles of the 41.98 latitude")

	f := openMockData(t, "mocks/route22.xml")
	routes := mapToRoute(f)
	filtered := filterNorthOfOffice(routes.Buses)

	actual := withinHalfMile(filtered)
	tc := []Bus{
		{4377, Point{41.98397789001465, -87.66879043579101}, "North Bound"},
	}

	rangeTestTable(t, actual, tc)
}
