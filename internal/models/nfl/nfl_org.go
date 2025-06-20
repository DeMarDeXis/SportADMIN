package nfl

type Team struct {
	Name        string `json:"name"`
	Abbr        string `json:"abbr"`
	ComUsedAbbr string `json:"com_used_abbr"`
}

type TeamDB struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Abbr        string `json:"abbr" db:"abbr"`
	ComUsedAbbr string `json:"com_used_abbr" db:"com_used_abbr"`
	Conference  string `json:"conference" db:"conference"`
	Divisions   string `json:"divisions" db:"divisions"`
	ImgUrl      string `json:"img_url" db:"img_url"`
}
