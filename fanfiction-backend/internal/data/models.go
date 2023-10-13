package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("Record not found")
)

var (
	TimeoutDuration = 3 * time.Second
)

type Models struct {
	Labels LabelModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Labels: LabelModel{DB: db},
	}
}
