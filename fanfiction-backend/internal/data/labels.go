package data

import "github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"

type Label struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	SubLabels []int64 `json:"sub_labels,omitempty"`
	Blacklist []int64 `json:"blacklist,omitempty"`
	Version   int32   `json:"version"`
}

// ValidateLabel is a helper function to validate a label
func ValidateLabel(v *validator.Validator, label *Label){
	v.Check(label.Name != "", "name", "cannot be empty")
	v.Check(len(label.Name) <= 100, "name", "must not be more than 100 bytes long")

	v.Check(validator.Unique(label.SubLabels), "sublabels", "must be unique")
	v.Check(validator.Unique(label.Blacklist), "blacklist", "must be unique")
}