package middleware

import (
	"fmt"
	"net/http"
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

// Next executes the next chain function
func Next(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before")
		next.ServeHTTP(w, r)
		fmt.Println("After")
	})
}
