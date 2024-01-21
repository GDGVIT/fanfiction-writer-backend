package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrDuplicateLabel     = errors.New("duplicate label")
	ErrDuplicateStory     = errors.New("duplicate story")
	ErrDuplicateTimeline  = errors.New("duplicate timeline")
	ErrDuplicateEvent     = errors.New("duplicate event")
	ErrDuplicateCharacter = errors.New("duplicate character")
)

// The amount of time given for a database command to run
var (
	TimeoutDuration = 3 * time.Second
)

type Models struct {
	Labels    LabelModel
	Users     UserModel
	Tokens    TokenModel
	Stories   StoryModel
	Timelines TimelineModel
	Events    EventModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Labels:    LabelModel{DB: db},
		Users:     UserModel{DB: db},
		Tokens:    TokenModel{DB: db},
		Stories:   StoryModel{DB: db},
		Timelines: TimelineModel{DB: db},
		Events:    EventModel{DB: db},
	}
}
