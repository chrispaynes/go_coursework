// Based off Todd McLeod's Code Clinic on Lynda.com
// Collects and outputs weather data from Lake Pend Oreille, Idaho.
// Calculates Mean and Median for wind speed, air temperature and barometric pressure
// DATA SOURCE: https://lpo.dt.navy.mil/data/
// https://lpo.dt.navy.mil/data/DM/Environmental_Data_Deep_Moor_2016.txt

package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type weatherRecord struct {
	pk       int
	date     string
	time     string
	temp     float64
	bPress   float64
	dewPt    float64
	relHum   float64
	windDir  float64
	windGust float64
	windSpd  float64
}

func main() {
	// logs start time for application
	start := time.Now()

	dataSource := "https://raw.githubusercontent.com/lyndadotcom/LPO_weatherdata/master/Environmental_Data_Deep_Moor_2015.txt"

	// Creates HTTP Get request from data_source argument
	// Returns response body or an error
	resp, err := http.Get(dataSource)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("Did not receive a HTTP response from %v\n", "\""+dataSource+"\"")
		log.Fatal(err)
	}

	// Creates NewReader from io.Reader interface
	// Reads records from a CSV-encoded file
	// Return value references the *Reader struct {}
	// By default NewReader returns an &Reader instance with "," set for the Rune "Comma" field
	rdr := csv.NewReader(resp.Body)

	// Reassigns Comma field to tab-delimitation
	rdr.Comma = '\t'

	// Sets object to ignore leading white space in a field
	rdr.TrimLeadingSpace = true

	// ReadAll uses *Reader to read the argument's remaining records
	// Reads and slices each record/line from source file until io.EOF
	// Appends each record/line slice to a "records" array
	// Returns the "records" array as record slices containing string literals
	rows, err := rdr.ReadAll()

	// Begins to panic when an error is present
	// Panic stops the ordinary flow control and begins panicking.
	// Panic stops the function execution
	// Executes defer functions and returns the panicked function to its caller
	if err != nil {
		panic(err)
	}

	// Loops through array printing each row
	// "_" signifies a Blank Identifier
	// Blank Identifiers instruct the program to ignore returned index values
	// "range" iterates through all entries within row
	// for _, row := range rows
	// 	fmt.Println(row)
	//	}

	// declares zero-value variables to hold accumulated totals
	// calculations are adding per each loop iteration
	// :var air_temp_total, baro_pres_total, wind_speed_total, counter float64;

	// RAW INPUT DATA
	// HEADER  [ date        time      Air_Temp  Barometric_Press  Dew_Point  Relative_Humidity   Wind_Dir  Wind_Gust  Wind_Speed ]
	// ROW 1:  [ 2015_06_04  01:09:21  57.70     29.95             51.22      79.00               163.40    12.00      10.00      ]
	// ROW 2:  [ 2015_06_04  01:09:21  57.70     29.95             51.22      79.00               163.40    12.00      10.00      ]
	// ROW N: ...

	// FORMATTED OUTPUT
	//           [Air_Temp   Barometric_Press   Wind_Speed ]
	// # 00001 - [57.70      29.95              10.00      ]
	// # 00002 - [57.70      29.95              10.00      ]
	// for i, row := range rows {
	//   fmt.Println("#", i, " - ", row[1], " ", row[2], " ", row[7])
	// }

	createTable(rows)

	// logs end time for application
	end := time.Now()
	delta := end.Sub(start)

	fmt.Println("\nProgram runtime is", delta)
}

// calcMean returns the mean value for a given data column
func calcMean(rows [][]string, col int) float64 {
	var total float64

	for _, row := range rows {
		val, _ := strconv.ParseFloat(row[col], 64)
		total += val
	}

	return total / float64(len(rows))
}

// calcMedian returns the median value for a given data column
func calcMedian(rows [][]string, index int) float64 {
	var vals []float64

	for _, row := range rows {
		val, _ := strconv.ParseFloat(row[index], 64)
		vals = append(vals, val)
	}

	sort.Float64s(vals)

	if len(vals)%2 == 0 {
		return getEvenMedian(vals)
	}

	mid := len(vals) / 2
	return vals[mid]
}

func getEvenMedian(n []float64) float64 {
	sort.Float64s(n)
	mid := len(n) / 2
	midLow := n[mid-1]
	midHigh := n[mid]
	return (midHigh + midLow) / 2
}

func createTable(rows [][]string) {
	fmt.Println("\n============================================================")
	fmt.Println("           LAKE PEND OREILLE, IDAHO WEATHER DATA            ")
	fmt.Printf("                   TOTAL RECORDS: %v\n", len(rows)-1)
	fmt.Println("============================================================")
	fmt.Println("Column\t\t\t Mean\t\t\t Median\n")
	fmt.Printf("%v\t\t %v\t %v\t\n", rows[0][1], calcMean(rows, 1), calcMedian(rows, 1))
	fmt.Printf("%v\t %v\t %v\t\n", rows[0][2], calcMean(rows, 2), calcMedian(rows, 2))
	fmt.Printf("%v\t\t %v\t %v\t\n", rows[0][7], calcMean(rows, 7), calcMedian(rows, 7))
	fmt.Println("============================================================")
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/pend_db_dev?sslmode=disable")

	createDBTable(db, "deep_moor_2015")

	if err != nil {
		log.Fatal(err)
	}

	return db, err
}

func createDBTable(db *sql.DB, table string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS deep_moor_2015(
    pk SERIAL PRIMARY KEY,
    date text,
    time time,
    air_temp numeric,
    barometric_press numeric,
    dew_point numeric,
    relative_humidity numeric,
    wind_dir numeric,
    wind_gust numeric,
    wind_speed numeric);`)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func insertDBRecord(db *sql.DB) {
	db.Exec("INSERT INTO deep_moor_2015 values (0, '2015_06_04', '01:09:21', 57.70, 29.95, 51.22, 79.00, 163.40, 12.00, 10.00)")
}
