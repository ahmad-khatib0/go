package models

import (
	"context"
	"database/sql"
	"time"
)

// DB model is the type of db connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModel returns a model type with database connection pool
func NewModel(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widget is the type of all widets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cnacel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cnacel()

	var widget Widget
	row := m.DB.QueryRowContext(ctx, "select id, name from widgets where id = ?", id)

	err := row.Scan(&widget.ID, &widget.Name)
	if err != nil {
		return widget, err
	}

	return widget, nil
}
