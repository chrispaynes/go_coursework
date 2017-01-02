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
1. Move into the chicago_transit folder
2. Compile the app with `$ go install`
3. Execute the app by calling the chicago_transit binary
```
$ cd go_workspace/src/chicago_transit/
$ go install
$ chicago_transit

---------------------------------------------------------------------
All BUSES ON ROUTE 22
---------------------------------------------------------------------
ID     Latitude             Longitude             Direction
1766   41.876529693603516   -87.62919616699219    North Bound
4162   41.89501132965088    -87.6296615600586     North Bound
1897   41.948798485522      -87.65769900915758    North West Bound
1919   41.93331241607666    -87.64532852172852    South East Bound
4332   41.973246932029724   -87.6680519580841     North Bound
4337   41.89396667480469    -87.63121032714844    South Bound
1901   42.018278333333335   -87.673035            North Bound
---------------------------------------------------------------------

---------------------------------------------------------------------
BUSES NORTH OF 41.98 LATITUDE
---------------------------------------------------------------------
ID     Latitude             Longitude     Direction
1901   42.018278333333335   -87.673035    North Bound
---------------------------------------------------------------------

---------------------------------------------------------------------
BUSES WITHIN 0.5 MILES!!!
---------------------------------------------------------------------
ID     Latitude             Longitude     Direction
---------------------------------------------------------------------
```

## Running Tests
 <b>[Go's Standard Testing Package](https://golang.org/pkg/testing/)</b>
Provides support for automated testing of Go packages.<br>


   ``` $ cd go_workspace/src/chicago_transit/``` for a basic Pass/Fail output<br>
   ``` $ go test``` for a basic Pass/Fail output<br>
   ``` $ go test -v ``` for a verbose output<br><br>

   To run an individual test, use `func TestMyCustomTest(...){...}` but strip off the "Test" prefix on the func name<br>
   ``` $ go test -v -run MyCustomTest```<br>

   Output basic Test Cover Percentages to command line<br>
   ``` $ go tool cover -func=cover.out ```<br>


