package main

import (
	"encoding/xml"
	"fmt"
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

func connectToRoute(route int) (*http.Response, error) {
	url := "http://ctabustracker.com/bustime/map/getBusesForRoute.jsp?route=" + strconv.Itoa(route)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Did not receive a HTTP response from %v\n", "\""+url+"\"")
		log.Fatal(err)
	}

	defer resp.Body.Close()

	return resp, err
}

func mapToRoute(body string) Route {
	r := Route{}
	xml.Unmarshal([]byte(body), &r)
	return r
}
