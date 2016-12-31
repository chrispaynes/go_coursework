# pend_oreille_weather_calc

## About
Collects and outputs weather data from Lake Pend Oreille, Idaho.
Calculates Mean and Median for wind speed, air temperature and barometric pressure

*This project was implemented from Todd McLeod's Golang Code Clinic on Lynda.com



#### Data Sources:
 * https://lpo.dt.navy.mil/data/
 * https://lpo.dt.navy.mil/data/DM/Environmental_Data_Deep_Moor_2016.txt

While recreating this app, lpo.dt.navy.mil HTTPS certificate was not signed properly by a CA. Unfortunately, This makes Go's http.Get() request harder to implement using the source behind https. I recommend hosting the txt file locally or storing the file on github. 


#### Featured Concepts:
 * CSV Parsing
 * HTTP/Networking
 * Median and Mean Calculations
 * String Conversion
 * Sorting 
