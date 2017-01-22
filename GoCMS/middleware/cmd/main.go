package main

import (
	"GoCMS/middleware"
	"context"
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Executing...")
	w.Write([]byte("Hello"))
}

func panicker(w http.ResponseWriter, r *http.Request) {
	panic(middleware.ErrInvalidEmail)
}

func withContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	bar := ctx.Value("foo")
	if str, ok := bar.(string); ok {
		w.Write([]byte(str))
	}
}

func main() {
	//sum := middleware.Add(1, 2, 3)
	//fmt.Println(sum)

	//chain := &middleware.Chain{0}
	//sum2 := chain.AddNext(1).AddNext(2).AddNext(3).Finally(6)
	//fmt.Println(sum2)

	logger := middleware.CreateLogger("middleware")
	http.Handle("/panic", middleware.Recover(panicker))
	http.Handle("/", middleware.Time(logger, hello))
	http.Handle("/context", middleware.PassContext(withContext))
	http.ListenAndServe(":3000", nil)
}
