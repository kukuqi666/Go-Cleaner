package cleaner

type Rule struct {
	Name        string   `json:"name"`
	Paths       []string `json:"paths"`
	Extensions  []string `json:"extensions"`
	MinSizeMB   int      `json:"min_size_mb"`
	MaxAgeDays  int      `json:"max_age_days"`
	Description string   `json:"description"`
}
