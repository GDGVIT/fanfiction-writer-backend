package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

var (
	TimeoutDuration = 3 * time.Second
)

type Models struct {
	Labels LabelModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Labels: LabelModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
