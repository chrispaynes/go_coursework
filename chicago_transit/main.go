package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Bus is one of many vehicles that belong to a Route
type Bus struct {
	ID        string `xml:"id"`
	Latitude  string `xml:"lat"`
	Direction string `xml:"d"`
}

// Route is a collection of buses following a similar path at a given time
type Route struct {
	Buses []Bus  `xml:"bus"`
	Time  string `xml:"time"`
}

func main() {

}

func fetchRouteData(route int) (io.ReadCloser, error) {
	url := "http://ctabustracker.com/bustime/map/getBusesForRoute.jsp?route=" + strconv.Itoa(route)
	resp, err := http.Get(url)
	body := resp.Body

	if err != nil {
		fmt.Printf("Did not receive a HTTP response from %v\n", "\""+url+"\"")
		log.Fatal(err)
	}

	return body, err
}

func mapToRoute(resp io.ReadCloser) Route {
	r := Route{}
	body, err := ioutil.ReadAll(resp)
	defer resp.Close()

	if err != nil {
		log.Fatal(err)
	}

	xml.Unmarshal(body, &r)
	return r
}

func createTable(route []Bus) string {
	var body string
	header := "\n----------------------------------------------------" +
		"\nBUSES NORTH OF 41.98 LATITUDE" +
		"\n----------------------------------------------------" +
		"\nID\t Latitude\t\t Direction\n"
	footer := "----------------------------------------------------"

	for _, bus := range route {
		body += fmt.Sprintf("%v\t %v\t %v\n", bus.ID, bus.Latitude, bus.Direction)
	}

	return header + body + footer
}

func filterNorthOfOffice(buses []Bus) []Bus {
	var filtered []Bus
	const officeLatitude = 41.98

	for _, bus := range buses {
		lat, _ := strconv.ParseFloat(bus.Latitude, 64)

		if lat > officeLatitude {
			filtered = append(filtered, bus)
		}
	}

	return filtered
}
