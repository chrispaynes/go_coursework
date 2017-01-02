# Chicago Transit - The Traveling Suitcase

Adapted from David Beazley's "Learn Python Through Public Data Hacking," [The Traveling Suitcase Coding Challenge](https://www.youtube.com/watch?v=RrPZza_vZ3w), presented at PyCon'13, March 13, 2013 in Santa Clara, California. Illustrates web scraping and analyzing real world data.

##The Challenge
Travis traveled to Chicago and left his suitcase on the Clark Street #22 bus on his way to Dave's office! Find a way to track down the current location of the suitcase. Travis doesn't know the Bus ID he was riding. Use the Chicago Transit Authority API to find likely candidates traveling northbound of Dave's office. Monitor the identified buses and report their current distance from Dave's office. When the bus gets closer than 0.5 miles, create an HTML alert showing the bus location on a Google Map.


##Todo
* Display a web alert when a bus gets within 0.5 miles of a bus location
* Display A Google Map of the Bus with the Suitcase

## Features
* Building Type Structs
* Chicago Transit Authority API
* Data Scraping
* Haversine Formula
* HTTP GET Request
* Google Maps API
* Testing Package
* XML Parsing

## Usage
Todo

## Running Tests
 * <b>[Go's Standard Testing Package](https://golang.org/pkg/testing/)</b>
Provides support for automated testing of Go packages.<br> 
   ``` $ go test -v ```
   To run an individual test, use `func TestMyCustomTest(...){...}` but strip of the "Test" prefix on the func name
   ``` $ go test -v -run MyCustomTest ./```
   Output basic Test Cover Percentages to command line
   ``` $ go tool cover -func=cover.out ```
