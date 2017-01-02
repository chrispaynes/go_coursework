# pend_oreille_weather_calc

## About
Collects and outputs weather data from Lake Pend Oreille, Idaho.
Calculates Mean and Median for wind speed, air temperature and barometric pressure

*This project was implemented from Todd McLeod's Golang Code Clinic on Lynda.com

## Features:
 * CSV Parsing
 * HTTP/Networking
 * Median and Mean Calculations
 * String Conversion
 * Sorting 

## Data Sources:
 * https://lpo.dt.navy.mil/data/
 * https://lpo.dt.navy.mil/data/DM/Environmental_Data_Deep_Moor_2016.txt

While recreating this app, lpo.dt.navy.mil HTTPS certificate was not signed properly by a CA. Unfortunately, This makes Go's http.Get() request harder to implement using the source behind https. I recommend hosting the txt file locally or storing the file on github. 

## Usage
1. Move into the pend_oreille_weather_calc directory
2. Run the app with `$ go run main.go`

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