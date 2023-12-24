package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
)

type Story struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	User_ID     int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Version     int       `json:"version"`
}

type StoryModel struct {
	DB *sql.DB
}

func ValidateStory(v *validator.Validator, story *Story) {
	v.Check(story.Title != "", "title", "must be provided")
}

func (m StoryModel) Insert(story *Story) error {
	query := `INSERT INTO stories(user_id, title, description)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version`

	args := []interface{}{story.User_ID, story.Title, story.Description}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&story.ID, &story.CreatedAt, &story.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "stories_user_id_title_key"`:
			return ErrDuplicateStory
		default:
			return err
		}
	}

	return nil
}

func (m StoryModel) GetForUser(user_id int64) ([]*Story, error) {
	query := `SELECT id, created_at, user_id, title, description, version
	FROM stories
	WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stories := []*Story{}

	for rows.Next() {
		var story Story

		err := rows.Scan(
			&story.ID,
			&story.CreatedAt,
			&story.User_ID,
			&story.Title,
			&story.Description,
			&story.Version,
		)
		if err != nil {
			return nil, err
		}

		stories = append(stories, &story)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stories, nil
}

func (m StoryModel) Get(user_id, story_id int64) (*Story, error) {
	query := `SELECT id, created_at, user_id, title, description, version
	FROM stories
	WHERE user_id = $1
	AND id = $2`

	var story Story

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, user_id, story_id).Scan(&story.ID, &story.CreatedAt, &story.User_ID, &story.Title, &story.Description, &story.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &story, nil
}

func (m StoryModel) Update(story *Story) error {
	query := `UPDATE stories
	SET title = $1, description = $2, version = version + 1
	WHERE user_id = $3 and id = $4 and version = $5
	RETURNING version`

	args := []interface{}{story.Title, story.Description, story.User_ID, story.ID, story.Version}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&story.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "stories_user_id_title_key"`:
			return ErrDuplicateStory
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m StoryModel) Delete(user_id, story_id int64) error {
	query := `DELETE FROM stories
	WHERE user_id = $1 
	AND id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, user_id, story_id)
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
