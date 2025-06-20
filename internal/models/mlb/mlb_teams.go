package mlb

type TeamParse struct {
	Abbreviation string `json:"abbr"`
	Name         string `json:"name"`
	Division     string `json:"division"`
	Conference   string `json:"conference"`
}

type TeamDB struct {
	ID         int    `db:"id"`
	Image      string `json:"img_url" db:"img_url"`
	Abbr       string `json:"abbr" db:"abbreviation"`
	Name       string `json:"name" db:"name"`
	Division   string `json:"division" db:"division"`
	Conference string `json:"conference" db:"conference"`
}
