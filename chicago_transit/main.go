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

// A Bus represent a vehicle belonging to a Route
type Bus struct {
	ID  string `xml:"id"`
	Lat string `xml:"lat"`
	Lon string `xml:"lon"`
	Dir string `xml:"d"`
}

// A Route represents a Bus collection traveling a similar path at a given time
type Route struct {
	Buses []Bus  `xml:"bus"`
	Time  string `xml:"time"`
}

type Point struct {
	lat float64
	lon float64
}

func main() {

}

func (b *Bus) sliceLat() string {
	return b.Lat[0:7]
}

func (b *Bus) sliceLon() string {
	return b.Lon[0:8]
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
		"\nID\t Latitude\t Longitude\t\t Direction\n"
	footer := "----------------------------------------------------"

	for _, bus := range route {
		body += fmt.Sprintf("%v\t %v\t %v\t %v\n", bus.ID, bus.Lat, bus.Lon, bus.Dir)
	}

	return header + body + footer
}

func filterNorthOfOffice(buses []Bus) []Bus {
	var filtered []Bus
	const officeLat = 41.98

	for _, bus := range buses {
		bus.Lon = bus.sliceLon()
		bus.Lat = bus.sliceLat()
		lat, _ := strconv.ParseFloat(bus.Lat, 64)

		if lat > officeLat {
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

	dLat := (p2.lat - p1.lat) * radians
	dLon := (p2.lon - p1.lon) * radians

	lat1 := p1.lat * radians
	lat2 := p2.lat * radians

	a := (math.Sin(dLat/2) * math.Sin(dLat/2)) +
		(math.Sin(dLon/2) * math.Sin(dLon/2) *
			math.Cos(lat1) * math.Cos(lat2))

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return (earthRadius * c)

}

func withinHalfMile(buses []Bus) []Bus {
	return []Bus{{"4377", "41.9839", "-87.6687", "North Bound"}}
}
