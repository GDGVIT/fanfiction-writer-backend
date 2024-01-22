package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
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

func ValidateEvent(v *validator.Validator, event *Event) {
	v.Check(event.Timeline_ID != 0, "timeline_id", "must be provided")
	v.Check(event.EventTime != time.Time{}, "event_time", "must be provided")
	v.Check(event.Title != "", "title", "must be provided")
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

func (m EventModel) Get(event_id int64) (*Event, error) {
	query := `SELECT id, created_at, timeline_id, event_time, title, description, details, version
	FROM events
	WHERE id = $1`

	var event Event

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, event_id).Scan(
		&event.ID,
		&event.CreatedAt,
		&event.Timeline_ID,
		&event.EventTime,
		&event.Title,
		&event.Description,
		&event.Details,
		&event.Version)

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

func (m EventModel) GetForTimeline(timeline_id int64) ([]*Event, error) {
	query := `SELECT id, created_at, timeline_id, event_time, title, description, details, version
	FROM events
	WHERE timeline_id = $1
	ORDER BY event_time ASC`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, timeline_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event

		err := rows.Scan(
			&event.ID,
			&event.CreatedAt,
			&event.Timeline_ID,
			&event.EventTime,
			&event.Title,
			&event.Description,
			&event.Details,
			&event.Version)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (m EventModel) Update(event *Event) error {
	query := `UPDATE events
	SET event_time = $1, timeline_id = $2, title = $3, description = $4, details = $5, version = version + 1
	WHERE id = $6 and version = $7
	RETURNING version`

	args := []interface{}{event.EventTime, event.Timeline_ID, event.Title, event.Description, event.Details, event.ID, event.Version}

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
