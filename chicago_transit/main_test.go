package main

import (
	"fmt"
	"testing"
)

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

func TestUnmarshalToSlice(t *testing.T) {
	t.Log("mapToRoute() stores XML in Routes struct\n")

	expected := "main.Route"
	xmlData := `
					<buses rt="20" rtRtpiFeedName="" rtdd="20">
						<!--  @11 @16  -->
						<time>8:59 AM</time>
						<bus>
							<id>1819</id>
							<consist/>
							<cars/>
							<rtpiFeedName/>
							<m>1</m>
							<rt>20</rt>
							<rtRtpiFeedName/>
							<rtdd>20</rtdd>
							<d>West Bound</d>
							<dd>Westbound</dd>
							<dn>W</dn>
							<lat>41.88179016113281</lat>
							<lon>-87.64892578125</lon>
							<pid>954</pid>
							<pd>Westbound</pd>
							<pdRtpiFeedName/>
							<run>5008</run>
							<fs>Austin</fs>
							<op>44687</op>
							<dip>6543</dip>
							<bid>1951</bid>
							<wid1>05</wid1>
							<wid2>008</wid2>
						</bus>
						<bus>
							<id>1838</id>
							<consist/>
							<cars/>
							<rtpiFeedName/>
							<m>1</m>
							<rt>20</rt>
							<rtRtpiFeedName/>
							<rtdd>20</rtdd>
							<d>West Bound</d>
							<dd>Westbound</dd>
							<dn>W</dn>
							<lat>41.880330579034215</lat>
							<lon>-87.76507226352034</lon>
							<pid>949</pid>
							<pd>Westbound</pd>
							<pdRtpiFeedName/>
							<run>5009</run>
							<fs>Austin</fs>
							<op>45143</op>
							<dip>10582</dip>
							<bid>1950</bid>
							<wid1>05</wid1>
							<wid2>009</wid2>
						</bus>
					<buses>
	`
	routes := mapToRoute(xmlData)
	actual := string(fmt.Sprintf("%T", routes))
	expectedRoute := Route{Buses: []Bus{
		{ID: "1819", Latitude: "41.88179016113281", Direction: "West Bound"},
		{ID: "1838", Latitude: "41.880330579034215", Direction: "West Bound"}},
		Time: "8:59 AM"}

	if actual != expected {
		t.Error("Expected type: \t", expected)
		t.Error("Received type: \t", actual)
	}

	for index, actualBus := range routes.Buses {
		if actualBus != expectedRoute.Buses[index] {
			t.Error("Expected: \t", actualBus)
			t.Error("Received: \t", expectedRoute.Buses[index])
		}
	}
}

// Test prints northbound buses that currently north the office_latitude (41.98)
// Test osutput buses to a table
