package main

import (
	"log"
	"net/http"

	"movieexamplekhubaib.com/movie/internal/controller/movie"
	metadatagateway "movieexamplekhubaib.com/movie/internal/gateway/metadata/http"
	ratinggateway "movieexamplekhubaib.com/movie/internal/gateway/rating/http"
	httphandler "movieexamplekhubaib.com/movie/internal/handler/http"
)

func main() {
	log.Println("Starting movie service")
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
