package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"gonum.org/v1/plot/plotter"

	"gonum.org/v1/plot"
)

type point struct {
	x, y int
}

var points = &[]point{}

func main() {
	f, err := generateRandomDataPoints("points.txt", time.Now().UnixNano())

	if err != nil {
		log.Fatalf("could not create data file: %v", err)
	}

	pts, err := scanDataPoints(f.Name(), points)
	_ = pts

	if err != nil {
		log.Fatalf("could not scan data points %v", err)
	}

	if err = plotData("plot.png", "Plot Title", "x", "y"); err != nil {
		log.Fatalf("could not plot data: %v", err)
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

func plotData(file, title, xlabel, ylabel string) error {
	p, err := plot.New()

	if err != nil {
		return err
	}

	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	data := make(plotter.XYs, 100)

	for i, a := range *points {
		data[i].X = float64(a.x)
		data[i].Y = float64(a.y)
	}

	s, err := plotter.NewScatter(data)

	if err != nil {
		return fmt.Errorf("could not create new scatter: %v", err)
	}

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	p.Add(s)

	return p.Save(600, 600, file)
}
