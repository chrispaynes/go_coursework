// Based off Todd McLeod's Code Clinic on Lynda.com
// Collects and outputs weather data from Lake Pend Oreille, Idaho.
// Calculates Mean and Median for wind speed, air temperature and barometric pressure
// DATA SOURCE: https://lpo.dt.navy.mil/data/
// https://lpo.dt.navy.mil/data/DM/Environmental_Data_Deep_Moor_2016.txt

package main

import (
  "fmt"
  "net/http"
  "log"
  "encoding/csv"
	"strconv"
)

func main() {
  /*  The original data source, "https://lpo.dt.navy.mil/data/DM/Environmental_Data_Deep_Moor_2016.txt"
      had issues with an unsigned HTTPS certificate. Had to use a source file from GitHub
  */
	data_source := "https://raw.githubusercontent.com/lyndadotcom/LPO_weatherdata/master/Environmental_Data_Deep_Moor_2015.txt"

	// Creates HTTP Get request from data_source argument
	// Returns response body or an error
  resp, err := http.Get(data_source)

	// Prints fatal errors and immediately calls os.Exit(1)
  if err != nil {
    fmt.Printf("Did not receive a HTTP response from %v\n", "\"" + data_source + "\"")
    log.Fatal(err)
  }

	// Defers until the surrounding function executes its return statement
		// OR defers until the function reaches its function body ending
		// OR defers because the corresponding goroutine is panicking
	// Closes the reader
  defer resp.Body.Close()

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

  // Prints table column headers
  fmt.Println("----------------------------------------------")
  fmt.Println(rows[0][1], rows[0][2], rows[0][7])
  fmt.Println("----------------------------------------------")

  // declares zero-value variables to hold accumulated totals
  // calculations are adding per each loop iteration
  var air_temp_total, baro_pres_total, wind_speed_total, counter float64;

	// Outputs => # 68923 - [2015_06_04 01:09:21 57.70 29.95 51.22 79.00 163.40 12.00 10.00]
	for i, row := range rows {
		fmt.Println("#", i, "-", row)

		// Parses columns/slices from strings to float64 precision
		// Prints first 10 records
			if i != 0 && i < 10 {
				air_temp, _ := strconv.ParseFloat(row[1], 64)
				baro_pres, _ := strconv.ParseFloat(row[2], 64)
				wind_speed, _ := strconv.ParseFloat(row[7], 64)
				fmt.Println(air_temp, "\t", baro_pres, "\t\t", wind_speed)
				counter++

				// calculates total value during each loop iteration
				air_temp_total += air_temp
				baro_pres_total += baro_pres
				wind_speed_total += wind_speed
			}
		}

	fmt.Println("----------------------------------------------")
	fmt.Println("TOTAL RECORDS:", counter)
	fmt.Println("TOTAL RECORDS:", len(rows) - 1)
	fmt.Println("----------------------------------------------")
	fmt.Println("Mean Air Temp:", air_temp_total/counter)
	fmt.Println("Mean Barometric Pressure:", baro_pres_total/counter)
	fmt.Println("Mean Wind Speed:", wind_speed_total/counter)
}
