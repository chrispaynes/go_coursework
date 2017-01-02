package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Bus struct {
	ID int `xml:"id"`
	Point
	Dir string `xml:"d"`
}

// A Route represents a Bus collection traveling a similar path at a given time
type Route struct {
	Buses []Bus  `xml:"bus"`
	Time  string `xml:"time"`
}

// A Point represents a location with latitude and longitude coordinates
type Point struct {
	Lat float64 `xml:"lat"`
	Lon float64 `xml:"lon"`
}

var office = &Point{41.9801433, -87.6683411}

func main() {
	data, err := fetchRouteData(22)

	if err != nil {
		log.Fatal(err)
	}

	route := mapToRoute(data)
	northOfOffice := filterNorthOfOffice(route.Buses)
	closeBy := withinHalfMile(northOfOffice)

	fmt.Println(createTable(route.Buses, "All BUSES ON ROUTE 22"))
	fmt.Println(createTable(northOfOffice, "BUSES NORTH OF 41.98 LATITUDE"))
	fmt.Println(createTable(closeBy, "BUSES WITHIN 0.5 MILES!!!"))

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

func filterNorthOfOffice(buses []Bus) []Bus {
	var filtered []Bus

	for _, bus := range buses {
		if bus.Point.Lat > office.Lat {
			filtered = append(filtered, bus)
		}
	}

	return filtered
}

func createTable(route []Bus, title string) string {
	var body string
	header := "\n----------------------------------------------------" +
		"\n" + title +
		"\n----------------------------------------------------" +
		"\nID\t Latitude\t Longitude\t\t Direction\n"
	footer := "----------------------------------------------------"

	for _, bus := range route {
		body += fmt.Sprintf("%v\t %v\t %v\t %v\n", bus.ID, bus.Point.Lat, bus.Point.Lon, bus.Dir)
	}

	return header + body + footer
}

func withinHalfMile(buses []Bus) []Bus {
	var filtered []Bus

	for _, bus := range buses {
		if findDistance(&bus.Point, office) < 0.6 {
			filtered = append(filtered, bus)
		}
	}

	return filtered
}

// findDistance uses the Haversine formula to calculate the great-circle
// distance between two points on a map.
// source: http://www.movable-type.co.uk/scripts/latlong.html
func findDistance(p1, p2 *Point) float64 {
	const earthRadius = 3961
	const radians = math.Pi / 180.0

	dLat := (p2.Lat - p1.Lat) * radians
	dLon := (p2.Lon - p1.Lon) * radians

	lat1 := p1.Lat * radians
	lat2 := p2.Lat * radians

	a := (math.Sin(dLat/2) * math.Sin(dLat/2)) +
		(math.Sin(dLon/2) * math.Sin(dLon/2) *
			math.Cos(lat1) * math.Cos(lat2))

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return (earthRadius * c)

}
