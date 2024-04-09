package search

import "github.com/jackc/pgx/v5/pgxpool"

type Search struct {
	db *pgxpool.Pool
}

type SearchArgs struct {
	DB *pgxpool.Pool
}

func NewSearch(ua SearchArgs) *Search {
	return &Search{db: ua.DB}
}
