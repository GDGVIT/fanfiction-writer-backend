package data

type Label struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	SubLabels []int64 `json:"sub_labels,omitempty"`
	Blacklist []int64 `json:"blacklist,omitempty"`
	Version   int32   `json:"version"`
}
