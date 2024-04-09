package messages

import "github.com/jackc/pgx/v5/pgxpool"

type Messages struct {
	db *pgxpool.Pool
}

type MessagesArgs struct {
	DB *pgxpool.Pool
}

func NewMessages(ua MessagesArgs) *Messages {
	return &Messages{db: ua.DB}
}
