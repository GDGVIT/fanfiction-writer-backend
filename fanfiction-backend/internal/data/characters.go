package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Character struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Story_ID    int64     `json:"story_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Index       int       `json:"index"`
	Version     int       `json:"version"`
}

type CharacterModel struct {
	DB *sql.DB
}

func ValidateCharacter(v *validator.Validator, character *Character) {
	v.Check(character.Name != "", "name", "must be provided")
}

func (m CharacterModel) Insert(character *Character) error {
	index, err := m.getLastIndex(character.Story_ID)
	index += 1
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			index = 1
		default:
			return err
		}
	}
	if character.Index == 0 || character.Index > index {
		character.Index = index
	}

	err = m.increaseIndex(character.Index, character.Story_ID)
	if err != nil {
		return err
	}

	query := `INSERT INTO characters(story_id, name, description, index)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version`

	args := []interface{}{character.Story_ID, character.Name, character.Description}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&character.ID, &character.CreatedAt, &character.Version)
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

func (m CharacterModel) InsertCharLabels(character_id uuid.UUID, label_id ...int64) error {
	query := `INSERT INTO characters_labels
	SELECT $1, labels.id FROM labels WHERE labels.id = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, character_id, pq.Array(label_id))
	fmt.Println(err)
	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "characters_labels" violates foreign key constraint "characters_labels_character_id_fkey""`:
			return ErrRecordNotFound
		case err.Error() == `pq: duplicate key value violates unique constraint "characters_labels_pkey"`:
			return ErrDuplicateLabel
		default:
			return err
		}
	}

	return nil
}

func (m CharacterModel) Get(character_id uuid.UUID) (*Character, error) {
	query := `SELECT id, created_at, story_id, name, description, index, version
	FROM characters
	WHERE id = $1`

	var character Character

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, character_id).Scan(&character.ID, &character.CreatedAt, &character.Story_ID, &character.Name, &character.Description, &character.Index, &character.Version)
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
	query := `SELECT id, created_at, story_id, name, description, index, version
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
			&character.Index,
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

func (m CharacterModel) GetIndexForCharacter(id uuid.UUID) (int, error) {
	query := `SELECT index
	FROM characters
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	var index int

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&index)
	if err != nil {
		return -1, err
	}

	return index, nil
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

func (m CharacterModel) Update(character *Character, oldIndex int) error {
	index, err := m.getLastIndex(character.Story_ID)
	index += 1
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			index = 1
		default:
			return err
		}
	}
	if character.Index == 0 || character.Index > index {
		character.Index = index
	}

	err = m.increaseIndex(character.Index, character.Story_ID)
	if err != nil {
		return err
	}

	query := `UPDATE characters
	SET name = $1, description = $2, index = $3, version = version + 1
	WHERE id = $4 AND version = $5
	RETURNING version`

	args := []interface{}{character.Name, character.Description, character.Index, character.ID, character.Version}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&character.Version)
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

	err = m.decreaseIndex(oldIndex, character.Story_ID)
	if err != nil {
		return err
	}

	return nil
}

func (m CharacterModel) Delete(story_id int64, character_id uuid.UUID) error {
	index, err := m.GetIndexForCharacter(character_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	query := `DELETE FROM characters
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, character_id)
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

	err = m.decreaseIndex(index, story_id)
	if err != nil {
		return err
	}

	return nil
}

func (m CharacterModel) DeleteCharLabels(character_id uuid.UUID, label_id ...int64) error {
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

func (m CharacterModel) getLastIndex(story_id int64) (int, error) {
	query := `SELECT index 
	FROM characters 
	WHERE story_id = $1
	ORDER BY index DESC
	fetch first 1 row only;
	`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	var index int

	err := m.DB.QueryRowContext(ctx, query, story_id).Scan(&index)
	if err != nil {
		return -1, err
	}

	return index, err

}

func (m CharacterModel) increaseIndex(index int, story_ID int64) error {
	query := `UPDATE characters
	SET index = index + 1
	WHERE index >= $1 AND
	story_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, index, story_ID)
	if err != nil {
		return err
	}

	return nil
}

func (m CharacterModel) decreaseIndex(index int, story_ID int64) error {
	query := `UPDATE characters
	SET index = index - 1
	WHERE index > $1 AND
	story_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, index, story_ID)
	if err != nil {
		return err
	}

	return nil
}
