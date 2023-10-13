package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
)

type Label struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	SubLabels []int64   `json:"sub_labels,omitempty"`
	Blacklist []int64   `json:"blacklist,omitempty"`
	Version   int32     `json:"version"`
}

// ValidateLabel is a helper function to validate a label
func ValidateLabel(v *validator.Validator, label *Label) {
	v.Check(label.Name != "", "name", "cannot be empty")
	v.Check(len(label.Name) <= 100, "name", "must not be more than 100 bytes long")

	v.Check(validator.Unique(label.SubLabels), "sublabels", "must be unique")
	v.Check(validator.Unique(label.Blacklist), "blacklist", "must be unique")
}

type LabelModel struct {
	DB *sql.DB
}

// Create a label entry in the database
func (m LabelModel) Create(label *Label) error {
	query := `INSERT INTO labels (name)
	VALUES $1
	RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, label.Name).Scan(&label.ID, &label.CreatedAt, &label.Version)
}

// Retrieve a specific label based on label_id
func (m LabelModel) Get(id int64) (*Label, error) {
	query := `SELECT id, created_at, name, version 
	FROM labels
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	var label Label

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&label.ID,
		&label.CreatedAt,
		&label.Name,
		&label.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	sublabels, err := m.GetAllSublabel(id)
	if err != nil {
		return nil, err
	}
	label.SubLabels = sublabels

	blacklist, err := m.GetAllBlacklistLabel(id)
	if err != nil {
		return nil, err
	}
	label.Blacklist = blacklist

	return &label, nil
}

// Retrieve all sublabels based on label_id
func (m LabelModel) GetAllSublabel(id int64) ([]int64, error) {
	query := `SELECT sublabel_id
	FROM sublabels
	WHERE label_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	sublabels := []int64{}

	for rows.Next() {
		var subid int64

		err := rows.Scan(&subid)
		if err != nil {
			return nil, err
		}

		sublabels = append(sublabels, subid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sublabels, nil
}

// Retrieve the label blacklist based on label_id
func (m LabelModel) GetAllBlacklistLabel(id int64) ([]int64, error) {
	query := `SELECT blacklist_id
	FROM blacklist_labels
	WHERE label_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	blacklist := []int64{}

	for rows.Next(){
		var blacklist_id int64

		err := rows.Scan(&blacklist_id)
		if err != nil {
			return nil, err
		}

		blacklist = append(blacklist, blacklist_id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blacklist, nil
}

func (m LabelModel) Update(name string) error {
	return nil
}

func (m LabelModel) Delete(id int64) error {
	return nil
}
