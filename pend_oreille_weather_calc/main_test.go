package main

import (
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
