package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Timeline struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Story_ID  int64     `json:"story_id"`
	Name      string    `json:"name"`
	Version   int       `json:"version"`
}

type TimelineModel struct {
	DB *sql.DB
}

func (m TimelineModel) Insert(timeline *Timeline) error {
	query := `INSERT INTO timelines(story_id, name)
	VALUES ($1, $2)
	RETURNING id, created_at, version`

	args := []interface{}{timeline.Story_ID, timeline.Name}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&timeline.ID, &timeline.CreatedAt, &timeline.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "timelines_story_id_name_key"`:
			return ErrDuplicateTimeline
		default:
			return err
		}
	}

	return nil
}

func (m TimelineModel) Get(story_id, timeline_id int64) (*Timeline, error) {
	query := `SELECT (id, created_at, story_id, name, version)
	FROM timelines
	WHERE story_id = $1
	AND id = $2`

	var timeline Timeline

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, story_id, timeline_id).Scan(&timeline.ID, &timeline.CreatedAt, &timeline.Story_ID, &timeline.Name, &timeline.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &timeline, nil
}

func (m TimelineModel) Update(timeline *Timeline) error {
	query := `UPDATE timelines
	SET name = $1, version = version + 1
	WHERE story_id = $2 AND id = $3 AND version = $4
	RETURNING version`

	args := []interface{}{timeline.Name, timeline.Story_ID, timeline.ID, timeline.Version}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&timeline.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "timelines_story_id_name_key"`:
			return ErrDuplicateStory
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m TimelineModel) Delete(story_id, timeline_id int64) error {
	query := `DELETE FROM timelines
	WHERE story_id = $1
	AND id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, story_id, timeline_id)
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
