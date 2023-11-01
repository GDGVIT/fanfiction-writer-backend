package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
	"github.com/lib/pq"
)

type Label struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	SubLabels []int64   `json:"sublabels,omitempty"`
	Blacklist []int64   `json:"blacklist,omitempty"`
	Version   int32     `json:"version"`
}

// ValidateLabel is a helper function to validate a label
func ValidateLabel(v *validator.Validator, label *Label) {
	v.Check(label.Name != "", "name", "cannot be empty")
	v.Check(len(label.Name) <= 100, "name", "must not be more than 100 bytes long")

	v.Check(validator.Unique(label.SubLabels), "sublabels", "must be unique")
	// v.Check(validator.In(label.Name, label.SubLabels), "sublabels", "cannot contain itself")

	v.Check(validator.Unique(label.Blacklist), "blacklist", "must be unique")
	// v.Check(validator.In(label.Name, label.Blacklist), "blacklist", "cannot contain itself")
}

type LabelModel struct {
	DB *sql.DB
}

// Create a label entry in the database
func (m LabelModel) Create(label *Label) error {
	query := `INSERT INTO labels (name)
	VALUES ($1)
	RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, label.Name).Scan(&label.ID, &label.CreatedAt, &label.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "labels_name_key"`:
			return ErrDuplicateLabel
		default:
			return err
		}
	}

	err = m.CreateSublabel(label.ID, label.SubLabels...)
	if err != nil {
		return err
	}

	err = m.CreateBlacklist(label.ID, label.Blacklist...)
	if err != nil {
		return err
	}

	return nil
}

// Create the sublabels of a label
func (m LabelModel) CreateSublabel(label_id int64, sublabel_ids ...int64) error {
	query := `INSERT INTO sublabels 
	SELECT $1, labels.id FROM labels where labels.id = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, label_id, pq.Array(sublabel_ids))
	return err
}

// Create the blacklist of a label
func (m LabelModel) CreateBlacklist(label_id int64, blacklist ...int64) error {
	query := `INSERT INTO blacklist_labels
	SELECT $1, labels.id FROM labels where labels.id = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, label_id, pq.Array(blacklist))
	return err
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

	for rows.Next() {
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

func (m LabelModel) Update(label *Label) error {
	query := `UPDATE labels
	SET name=$1, version = version + 1
	WHERE id=$2 AND version = $3
	RETURNING version`

	args := []interface{}{
		label.Name,
		label.ID,
		label.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&label.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "labels_name_key"`:
			return ErrDuplicateLabel
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete a label based on label_id. Due to CASCADE in database, all related records are also deleted
func (m LabelModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM labels
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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

// DeleteSubLabel deletes the sublabels associated with a given label. If no sublabels are given, all are deleted.
func (m LabelModel) DeleteSublabel(label_id int64, sublabel_ids ...int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	var (
		result sql.Result
		err    error
	)

	if len(sublabel_ids) == 0 {
		query := `DELETE FROM sublabels 
		WHERE label_id = $1`

		result, err = m.DB.ExecContext(ctx, query, label_id)
	} else {
		query := `DELETE FROM sublabels 
		WHERE label_id = $1 AND  sublabel_id = ANY($2)`

		result, err = m.DB.ExecContext(ctx, query, label_id, pq.Array(sublabel_ids))
	}

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

// DeleteSubLabel deletes the sublabels associated with a given label. If no sublabels are given, all are deleted.
func (m LabelModel) DeleteBlacklist(label_id int64, blacklist ...int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	var (
		result sql.Result
		err    error
	)

	if len(blacklist) == 0 {
		query := `DELETE FROM blacklist_labels 
		WHERE label_id = $1`

		result, err = m.DB.ExecContext(ctx, query, label_id)
	} else {
		query := `DELETE FROM blacklist_labels 
		WHERE label_id = $1 AND  blacklist_id = ANY($2)`

		result, err = m.DB.ExecContext(ctx, query, label_id, pq.Array(blacklist))
	}

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
