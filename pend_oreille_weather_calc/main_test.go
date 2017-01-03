package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func TestCalcMeanCalculation(t *testing.T) {
	t.Log("calcMean() returns the correct mean")
	rows1 := [][]string{[]string{"0"}, []string{"10"}}
	expected1 := 5.0
	result1 := calcMean(rows1, 0)

	rows2 := [][]string{[]string{""}, []string{"23.40"}, []string{"1453.2308"},
		[]string{"2544.994"}, []string{"3985.66"}, []string{"44.4"}}
	expected2 := 1341.9474666666667
	result2 := calcMean(rows2, 0)

	rows3 := [][]string{[]string{"1", "11"}, []string{"2", "22"}, []string{"3", "33"}}
	expected3 := 22.0
	result3 := calcMean(rows3, 1)

	if result1 != expected1 {
		t.Error("Received:\t", result1)
		t.Error("Expected:\t", expected1)
	}

	if result2 != expected2 {
		t.Error("Received:\t", result2)
		t.Error("Expected:\t", expected2)
	}

	if result3 != expected3 {
		t.Error("Received:\t", result3)
		t.Error("Expected:\t", expected3)
	}
}

func TestCalcMedianCalculation(t *testing.T) {
	t.Log("calcMedian() returns the correct median")
	rows1 := [][]string{[]string{"0"}, []string{"10"}}
	expected1 := 5.0
	result1 := calcMedian(rows1, 0)

	rows2 := [][]string{[]string{""}, []string{"23.40"}, []string{"1453.2308"},
		[]string{"2544.994"}, []string{"3985.66"}, []string{"44.4"}}
	expected2 := 748.8154000000001
	result2 := calcMedian(rows2, 0)

	rows3 := [][]string{[]string{"1", "11"}, []string{"2", "22"}, []string{"3", "33"}}
	expected3 := 22.0
	result3 := calcMedian(rows3, 1)

	if result1 != expected1 {
		t.Error("Received:\t", result1)
		t.Error("Expected:\t", expected1)
	}

	if result2 != expected2 {
		t.Error("Received:\t", result2)
		t.Error("Expected:\t", expected2)
	}

	if result3 != expected3 {
		t.Error("Received:\t", result3)
		t.Error("Expected:\t", expected3)
	}
}

func TestCalcEvenMedianCalculation(t *testing.T) {
	t.Log("getEvenMedian() returns the correct median when given an even amount of elements")

	n1 := []float64{10, 0}
	expected1 := 5.0
	result1 := getEvenMedian(n1)

	n2 := []float64{999, 0, 30, 20, 10, -999}
	expected2 := 15.0
	result2 := getEvenMedian(n2)

	n3 := []float64{56.43, 69.56, 99.34, -99.56, 459.3}
	expected3 := 62.995000000000005
	result3 := getEvenMedian(n3)

	if result1 != expected1 {
		t.Error("Received:\t", result1)
		t.Error("Expected:\t", expected1)
	}

	if result2 != expected2 {
		t.Error("Received:\t", result2)
		t.Error("Expected:\t", expected2)
	}

	if result3 != expected3 {
		t.Error("Received:\t", result3)
		t.Error("Expected:\t", expected3)
	}
}

func newTestDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/pend_db_dev_test?sslmode=disable")

	return db, err
}

func newTestDBTableCreation(db *sql.DB) (sql.Result, error) {
	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS deep_moor_2015(
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

	return res, err
}

// newTestDBInsertion inserts randomized weather data into the database
func newTestDBInsertion(db *sql.DB) (sql.Result, error) {
	airTmp := randomizeData(19, 70)
	barPres := randomizeData(0.333, 29.333)
	dewPt := randomizeData(14, 56)
	relHum := randomizeData(36, 98)
	wndDir := randomizeData(360, 1)
	wndGust := randomizeData(36, 1)
	wndSpd := randomizeData(22, 1)

	vals := "'2015_06_04', '01:09:21'," + airTmp + barPres + dewPt + relHum + wndDir + wndGust + wndSpd[0:4]

	res, err := db.Exec(`INSERT INTO deep_moor_2015 (date, time, air_temp,
    barometric_press, dew_point, relative_humidity, wind_dir, wind_gust,
    wind_speed) values(` + vals + `);`)

	return res, err
}

// randomizeData prepares randomized weather data based on mix and max values.
// Converts values to a string and appends ", " so it concatenates into a SQL
// Insertion string. The last randomized value in the SQL Insertion needs to
// slice off ", " to prevent a SQL Insertion Error. Eg: wndGust + wndSpd[0:4]
func randomizeData(min, max float64) string {
	r := (rand.Float64() * min) + (max - min)

	return strconv.FormatFloat(math.Abs(r), 'f', 2, 64) + ", "
}

func TestConnectToDatabase(t *testing.T) {
	t.Log("initDB() connects to a Postgres DB")

	var expected string
	db, err := newTestDB()
	db.Exec("DROP TABLE deep_moor_2015;")
	defer db.Close()

	result := db.Ping()
	_, tblErr := newTestDBTableCreation(db)
	_, insErr := newTestDBInsertion(db)

	for i := 0; i < 9; i++ {
		newTestDBInsertion(db)
	}

	if err != nil {
		t.Fatal(err)
	}

	if result != nil || tblErr != nil || insErr != nil {
		t.Error("Received Ping() response:\t", result)
		t.Error("Expected Ping() response:\t", expected)
		t.Error("Received error when creating DB table:\t", tblErr)
		t.Error("Received error when inserted into table:\t", insErr)
		t.Fatal(err)
	}
}

func TestCreateDBTable(t *testing.T) {
	t.Log("can insert table into database")

	db, err := newTestDB()
	defer db.Close()

	db.Exec("DROP TABLE deep_moor_2015;")

	_, result := newTestDBTableCreation(db)

	if result != nil || err != nil {
		t.Fatal("Received error when connecting to test DB:\t", err)
		t.Error("Received error when creating test DB table:\t", result)
	}
}

func TestTableValueInsertion(t *testing.T) {
	t.Log("can insert values into database table")

	db, err := newTestDB()
	defer db.Close()

	db.Exec("DROP TABLE deep_moor_2015;")
	_, tblErr := newTestDBTableCreation(db)

	_, result := newTestDBInsertion(db)

	for i := 0; i < 9; i++ {
		newTestDBInsertion(db)
	}

	if result != nil || err != nil || tblErr != nil {
		t.Fatal("Received error when connecting to test DB:\t", err)
		t.Fatal("Received error when creating test DB table:\t", tblErr)
		t.Error("Received error when inserting values into DB table:\t", result)
	}
}
