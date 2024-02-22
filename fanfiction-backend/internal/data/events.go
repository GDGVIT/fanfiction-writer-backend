package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Character_ID uuid.UUID `json:"character_id"`
	// EventTime    time.Time `json:"event_time"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details,omitempty"`
	Index       int    `json:"index"`
	Version     int    `json:"-"`
}

type Story_Event struct {
	Character_ID    uuid.UUID `json:"character_id"`
	Character_Index int       `json:"character_index"`
	Character_Name  string    `json:"character_name"`
	Character_Desc  string    `json:"character_description"`
	Events          []*Event  `json:"events"`
}

func ValidateEvent(v *validator.Validator, event *Event) {
	v.Check(event.Character_ID != uuid.Nil, "character_id", "must be provided")
	// v.Check(event.EventTime != time.Time{}, "event_time", "must be provided")
	v.Check(event.Title != "", "title", "must be provided")
}

type EventModel struct {
	DB *sql.DB
}

func (m EventModel) Insert(event *Event) error {
	index, err := m.getLastIndex(event.Character_ID)
	index += 1
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			index = 1
		default:
			return err
		}
	}
	if event.Index == 0 || event.Index > index {
		event.Index = index
	}

	err = m.increaseIndex(event.Index, event.Character_ID)
	if err != nil {
		return err
	}

	query := `INSERT INTO events(character_id, title, description, details, index)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, version`

	args := []interface{}{event.Character_ID, event.Title, event.Description, event.Details, event.Index}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&event.ID, &event.CreatedAt, &event.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates primary key constraint "events_pkey"`:
			return ErrDuplicateEvent
		default:
			return err
		}
	}

	return nil

}

func (m EventModel) Get(event_id uuid.UUID) (*Event, error) {
	query := `SELECT id, created_at, character_id, title, description, details, index, version
	FROM events
	WHERE id = $1`

	var event Event

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, event_id).Scan(
		&event.ID,
		&event.CreatedAt,
		&event.Character_ID,
		// &event.EventTime,
		&event.Title,
		&event.Description,
		&event.Details,
		&event.Index,
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

func (m EventModel) GetForCharacter(character_id uuid.UUID) ([]*Event, error) {
	query := `SELECT id, created_at, character_id, title, description, details, index, version
	FROM events
	WHERE character_id = $1
	ORDER BY index ASC`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, character_id)
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
			&event.Character_ID,
			// &event.EventTime,
			&event.Title,
			&event.Description,
			&event.Details,
			&event.Index,
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

func (m EventModel) GetForStory(id int64) ([]*Story_Event, error) {
	query := `SELECT c.id, c.index, c.name, c.description, e.id, e.created_at, e.character_id, e.title, e.description, e.details, e.index, e.version
	FROM events e
	LEFT JOIN characters c
	ON e.character_id = c.id
	WHERE c.story_id = $1
	ORDER BY c.index ASC, e.index ASC`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	story_events := []*Story_Event{}
	for rows.Next() {
		var story_event Story_Event

		var event Event

		err := rows.Scan(
			&story_event.Character_ID,
			&story_event.Character_Index,
			&story_event.Character_Name,
			&story_event.Character_Desc,
			&event.ID,
			&event.CreatedAt,
			&event.Character_ID,
			// &event.EventTime,
			&event.Title,
			&event.Description,
			&event.Details,
			&event.Index,
			&event.Version)
		if err != nil {
			return nil, err
		}

		story_event.Events = append(story_event.Events, &event)
		story_events = append(story_events, &story_event)
	}

	var concat_story_events []*Story_Event
	var char_id uuid.UUID
	j := 0
	for _, event := range story_events {
		if char_id == uuid.Nil {
			char_id = event.Character_ID
			concat_story_events = append(concat_story_events, event)
		} else if event.Character_ID == char_id {
			concat_story_events[j].Events = append(concat_story_events[j].Events, event.Events...)
		} else {
			j += 1
			char_id = event.Character_ID
			concat_story_events = append(concat_story_events, event)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return concat_story_events, nil

}

func (m *EventModel) GetAllForStory(story_id int64) ([]*Story_Event, error) {
	query := `SELECT id, name, description, index
	FROM characters
	WHERE story_id = $1
	ORDER BY index`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, story_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	story_events := []*Story_Event{}

	for rows.Next() {
		var story_event Story_Event

		err := rows.Scan(
			&story_event.Character_ID,
			&story_event.Character_Name,
			&story_event.Character_Desc,
			&story_event.Character_Index,
		)
		if err != nil {
			return nil, err
		}

		story_events = append(story_events, &story_event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, v := range story_events {
		events, err := m.GetForCharacter(v.Character_ID)
		if err != nil {
			return nil, err
		}

		v.Events = events
	}

	return story_events, nil
}

func (m EventModel) GetIndexForEvent(id uuid.UUID) (int, error) {
	query := `SELECT index
	FROM events
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

func (m EventModel) Update(event *Event, oldindex int, oldCharacterID uuid.UUID) error {
	index, err := m.getLastIndex(event.Character_ID)
	index += 1
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			index = 1
		default:
			return err
		}
	}
	if event.Index == 0 || event.Index > index {
		event.Index = index
	}

	err = m.increaseIndex(event.Index, event.Character_ID)
	if err != nil {
		return err
	}

	query := `UPDATE events
	SET character_id = $1, title = $2, description = $3, details = $4, index = $5, version = version + 1
	WHERE id = $6 and version = $7
	RETURNING version`

	args := []interface{}{event.Character_ID, event.Title, event.Description, event.Details, event.Index, event.ID, event.Version}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err = m.DB.QueryRowContext(ctx, query, args...).Scan(&event.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates primary key constraint "events_pkey"`:
			return ErrDuplicateEvent
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	err = m.decreaseIndex(oldindex, oldCharacterID)
	if err != nil {
		return err
	}

	return nil
}

func (m EventModel) Delete(event_id, character_id uuid.UUID) error {
	index, err := m.GetIndexForEvent(event_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	query := `DELETE FROM events
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, event_id)
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

	err = m.decreaseIndex(index, character_id)
	if err != nil {
		return err
	}

	return nil
}

func (m EventModel) getLastIndex(character_id uuid.UUID) (int, error) {
	query := `SELECT index 
	FROM events 
	WHERE character_id = $1
	ORDER BY index DESC
	fetch first 1 row only;
	`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	var index int

	err := m.DB.QueryRowContext(ctx, query, character_id).Scan(&index)
	if err != nil {
		return -1, err
	}

	return index, err

}

func (m EventModel) increaseIndex(index int, character_ID uuid.UUID) error {
	query := `UPDATE events
	SET index = index + 1
	WHERE index >= $1 AND
	character_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, index, character_ID)
	if err != nil {
		return err
	}

	return nil
}

func (m EventModel) decreaseIndex(index int, character_ID uuid.UUID) error {
	query := `UPDATE events
	SET index = index - 1
	WHERE index > $1 AND
	character_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, index, character_ID)
	if err != nil {
		return err
	}

	return nil
}
