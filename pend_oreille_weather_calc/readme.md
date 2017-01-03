# pend_oreille_weather_calc

## About
Collects and outputs weather data from Lake Pend Oreille, Idaho.
Calculates Mean and Median for wind speed, air temperature and barometric pressure

*This project was implemented from Todd McLeod's Golang Code Clinic on Lynda.com

## Features:
* CSV Parsing
* HTTP/Networking
* Median and Mean Calculations
* Postgres Connectivty
* Randomizing Mock DB Records
* String Conversion
* Sorting 
* SQL Statements
* Testing

## Data Sources:
* https://lpo.dt.navy.mil/data/
* https://lpo.dt.navy.mil/data/DM/Environmental_Data_Deep_Moor_2016.txt

While recreating this app, lpo.dt.navy.mil HTTPS certificate was not signed properly by a CA. Unfortunately, This makes Go's http.Get() request harder to implement using the source behind https. I recommend hosting the txt file locally or storing the file on github. 

## Installation
1. [Install Postgres](https://wiki.postgresql.org/wiki/Detailed_installation_guides)
2. Create UTF-8 `pend_db_dev` and `pend_db_dev_test` databases
3. Fetch Go dependencies `go get ./`

## Usage
1. Move into the pend_oreille_weather_calc directory
2. Use `pg_ctl` to start the PostgreSQL server
3. Connect to the databases using `psql`
3. Run the app with `$ go run main.go`

```
$ cd go_workspace/src/pend_oreille_weather_calc/
$ go run main.go

============================================================
           LAKE PEND OREILLE, IDAHO WEATHER DATA
                   TOTAL RECORDS: 68923
============================================================
Air_Temp           43.13632667866028     40.96
Barometric_Press   30.149644100748514    30.16
Wind_Speed         6.475129127734979     5
============================================================

```

## Running Tests
 <b>[Go's Standard Testing Package](https://golang.org/pkg/testing/)</b>
Provides support for automated testing of Go packages.<br>


   ``` $ cd go_workspace/src/pend_oreille_weather_calc/```<br>
   ``` $ go test``` for a basic Pass/Fail output<br>
   ``` $ go test -v ``` for a verbose output<br><br>

   To run an individual test, use `func TestMyCustomTest(...){...}` but strip off the "Test" prefix on the func name<br>
   ``` $ go test -v -run MyCustomTest```<br>

   Output basic Test Cover Percentages to command line<br>
   ``` $ go tool cover -func=cover.out ```<br>