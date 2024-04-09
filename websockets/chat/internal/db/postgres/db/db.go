package db

import "github.com/jackc/pgx/v5/pgxpool"

type DB struct {
	db *pgxpool.Pool
}

type DBArgs struct {
	DB *pgxpool.Pool
}

func NewDB(ua DBArgs) *DB {
	return &DB{db: ua.DB}
}
