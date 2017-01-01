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

	if err != nil {
		fmt.Printf("Did not receive a HTTP response from %v\n", "\""+url+"\"")
		log.Fatal(err)
	}

	defer resp.Body.Close()
	return resp.Body, err
}

func mapToRoute(resp io.ReadCloser) Route {
	r := Route{}
	body, err := ioutil.ReadAll(resp)

	if err != nil {
		log.Fatal(err)
	}

	xml.Unmarshal(body, &r)
	return r
}
