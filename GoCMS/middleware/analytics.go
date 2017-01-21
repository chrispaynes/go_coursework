package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Add is a variadic function that sums ints
func Add(nums ...int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// Chain stores a sum
type Chain struct {
	Sum int
}

// AddNext is a chainable sum function
func (c *Chain) AddNext(num int) *Chain {
	c.Sum += num
	return c
}

// Finally represents the function chain terminus and returns the fina l sum
func (c *Chain) Finally(num int) int {
	return c.Sum + num
}

// CreateLogger creates a new logger that writes to a file
// TODO: LOG TO DATABASE
func CreateLogger(filename string) *log.Logger {
	file, err := os.OpenFile(filename+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

// Time executes the next chain function
func Time(logger *log.Logger, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		logger.Println(elapsed)
	})
}
