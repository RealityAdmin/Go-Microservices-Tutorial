package main

import (
	"log"
	"net/http"

	"movieexamplekhubaib.com/rating/internal/controller/rating"
	httphandler "movieexamplekhubaib.com/rating/internal/handler/http"
	"movieexamplekhubaib.com/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
