package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type point struct {
	x, y int
}

var points = &[]point{}

func main() {
	var err error
	f, err := generateRandomDataPoints("points.txt", time.Now().UnixNano())

	if err != nil {
		log.Fatal("could not create data file:", err)
	}

	p, err := scanDataPoints(f.Name(), points)

	if err != nil {
		log.Fatal("could not scan data points:", err)
	}

	fmt.Print((*p)[0].y)
}

func generateRandomDataPoints(file string, seed int64) (*os.File, error) {
	f, err := os.Create(file)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	rand.Seed(seed)

	linebreak := "\n"

	for i := 0; i < 100; i++ {
		x := rand.Intn(200)
		y := rand.Intn(200)

		if i == 99 {
			linebreak = ""
		}
		_, err := f.Write([]byte(fmt.Sprintf("%d,%d%s", x, y, linebreak)))

		if err != nil {
			os.Remove(file)
			return nil, err
		}
	}

	return f, nil
}

func scanDataPoints(file string, p *[]point) (*[]point, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		var x, y int
		_, err := fmt.Sscanf(s.Text(), "%d,%d", &x, &y)

		if err != nil {
			return nil, err
		}

		*p = append(*p, point{x, y})
	}

	return p, nil
}
