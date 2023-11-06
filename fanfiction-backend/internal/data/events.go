package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Timeline_ID int64     `json:"timeline_id"`
	EventTime   time.Time `json:"event_time"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Details     string    `json:"details,omitempty"`
	Version     int       `json:"-"`
}

type EventModel struct {
	DB *sql.DB
}

func (m EventModel) Insert(event *Event) error {
	query := `INSERT INTO events(timeline_id, event_time, title, description, details)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, version`

	args := []interface{}{event.Timeline_ID, event.EventTime, event.Title, event.Description, event.Details}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&event.ID, &event.CreatedAt, &event.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "events_timeline_id_title_key"`:
			return ErrDuplicateEvent
		default:
			return err
		}
	}

	return nil

}

func (m EventModel) Get(event_id, timeline_id int64) (*Event, error) {
	query := `SELECT id, created_at, timeline_id, event_time, title, description, details, version
	FROM events
	WHERE timeline_id = $1
	AND id = $2`

	var event Event

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, timeline_id, event_id).Scan(&event.ID, &event.CreatedAt, &event.EventTime, &event.Title, &event.Description, &event.Details, &event.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &event, nil
}

func (m EventModel) Update(event *Event) error {
	query := `UPDATE events
	SET event_time = $1, title = $2, description = $3, details = $4, version = version + 1
	WHERE timeline_id = $5 AND id = $6 and version = $7
	RETURNING version`

	args := []interface{}{event.EventTime, event.Title, event.Description, event.Details, event.Timeline_ID, event.ID, event.Version}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&event.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "events_timeline_id_title_key"`:
			return ErrDuplicateStory
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m EventModel) Delete(event_id, timeline_id int64) error {
	query := `DELETE FROM events
	WHERE timeline_id = $1
	AND id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, timeline_id, event_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
