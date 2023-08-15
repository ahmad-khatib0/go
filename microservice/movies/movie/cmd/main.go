package main

import (
	"log"
	"net/http"

	"github.com/ahmad-khatib0/go/microservice/movies/movie/internal/controller/movie"

	metadataGateway "github.com/ahmad-khatib0/go/microservice/movies/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/ahmad-khatib0/go/microservice/movies/movie/internal/gateway/rating/http"
	httphandler "github.com/ahmad-khatib0/go/microservice/movies/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")

	metadataGateway := metadataGateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
