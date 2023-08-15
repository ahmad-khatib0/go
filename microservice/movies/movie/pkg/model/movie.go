package model

import (
	model "github.com/ahmad-khatib0/go/microservice/movies/metadata/pkg/model"
)

// MovieDetails includes movie metadata its aggregated rating.
type MovieDetails struct {
	Rating   *float64       `json:"rating,omitEmpty"`
	Metadata model.Metadata `json:"metadata"`
}
