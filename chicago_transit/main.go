package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Bus struct {
	XMLName   xml.Name `xml:"bus"`
	ID        string   `xml:"id"`
	Latitude  string   `xml:"lat"`
	Direction string   `xml:"d"`
}

type Route struct {
	XMLName xml.Name `xml:"buses"`
	Buses   []Bus    `xml:"bus"`
	Time    string   `xml:"time"`
}

func main() {

}

Unmarshal parses the XML-encoded data and stores the result in the value pointed to by v, which must be an arbitrary struct, slice, or string. Well-formed data that does not fit into v is discarded.



func connectToRoute(route int) (*http.Response, error) {
	url := "http://ctabustracker.com/bustime/map/getBusesForRoute.jsp?route=" + strconv.Itoa(route)
	resp, err := http.Get(url)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
  
	xml_data := Route{}

	routes := xml.Unmarshal([]byte(body), &xml_data)
	fmt.Println(routes)
	fmt.Println(&xml_data)

	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Did not receive a HTTP response from %v\n", "\""+url+"\"")
		log.Fatal(err)
	}

	return resp, err
}
