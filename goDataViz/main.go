package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	_, err := generateRandomDataPoints("points.txt", time.Now().UnixNano())

	if err != nil {
		log.Fatal("could not create data file:", err)
	}

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
