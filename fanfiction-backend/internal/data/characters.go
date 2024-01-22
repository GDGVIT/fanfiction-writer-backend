package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
	"github.com/lib/pq"
)

type Character struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Story_ID    int64     `json:"story_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Version     int       `json:"version"`
}

type CharacterModel struct {
	DB *sql.DB
}

func ValidateCharacter(v *validator.Validator, character *Character) {
	v.Check(character.Name != "", "name", "must be provided")
}

func (m CharacterModel) Insert(character *Character) error {
	query := `INSERT INTO characters(story_id, name, description)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version`

	args := []interface{}{character.Story_ID, character.Name, character.Description}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&character.ID, &character.CreatedAt, &character.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "characters_story_id_name_key"`:
			return ErrDuplicateCharacter
		case err.Error() == `pq: insert or update on table "characters" violates foreign key constraint "characters_story_id_fkey"`:
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m CharacterModel) InsertCharLabels(character_id int64, label_id ...int64) error {
	query := `INSERT INTO characters_labels
	SELECT $1, labels.id FROM labels WHERE labels.id = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()
	
	_, err := m.DB.ExecContext(ctx, query, character_id, pq.Array(label_id))
	return err
}

func (m CharacterModel) Get(story_id, character_id int64) (*Character, error) {
	query := `SELECT id, created_at, story_id, name, description, version
	FROM characters
	WHERE story_id = $1
	AND id = $2`

	var character Character

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, story_id, character_id).Scan(&character.ID, &character.CreatedAt, &character.Story_ID, &character.Name, &character.Description, &character.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &character, nil
}

func (m CharacterModel) GetForStory(story_id int64) ([]*Character, error) {
	query := `SELECT id, created_at, story_id, name, description, version
	FROM characters
	WHERE story_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, story_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	characters := []*Character{}

	for rows.Next() {
		var character Character

		err := rows.Scan(
			&character.ID,
			&character.CreatedAt,
			&character.Story_ID,
			&character.Name,
			&character.Description,
			&character.Version,
		)
		if err != nil {
			return nil, err
		}

		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}


func (m CharacterModel) GetAllForLabel(label_id int64) ([]*Character, error) {
	query := `
	SELECT DISTINCT c.id, c.created_at, c.story_id, c.name, c.description, c.version
        FROM characters c
        INNER JOIN characters_labels cl ON c.id = cl.character_id
        INNER JOIN labels l ON cl.label_id = l.id
        WHERE l.id = $1 OR l.id IN (SELECT sublabel_id
        FROM sublabels
        WHERE label_id = $1)
		ORDER BY c.id;
	`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()
	
	rows, err := m.DB.QueryContext(ctx, query, label_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	characters := []*Character{}

	for rows.Next() {
		var character Character

		err := rows.Scan(
			&character.ID,
			&character.CreatedAt,
			&character.Story_ID,
			&character.Name,
			&character.Description,
			&character.Version,
		)
		if err != nil {
			return nil, err
		}

		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}	


func (m CharacterModel) Update(character *Character) error {
	query := `UPDATE characters
	SET name = $1, description = $2, version = version + 1
	WHERE story_id = $3 AND id = $4 AND version = $5
	RETURNING version`

	args := []interface{}{character.Name, character.Description, character.Story_ID, character.ID, character.Version}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&character.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "characters_story_id_name_key"`:
			return ErrDuplicateStory
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m CharacterModel) Delete(story_id, character_id int64) error {
	query := `DELETE FROM characters
	WHERE story_id = $1
	AND id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, story_id, character_id)
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

func (m CharacterModel) DeleteCharLabels(character_id int64, label_id ...int64) error {
	query := `DELETE FROM characters_labels
	WHERE character_id = $1
	AND label_id = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, character_id, pq.Array(label_id))
	
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