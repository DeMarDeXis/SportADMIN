package nba

type Team struct {
	Abbreviation string `json:"abbr"`
	Name         string `json:"name"`
}

type TeamDB struct {
	ID    int    `db:"id"`
	Image string `json:"img_url" db:"img_url"`
	Abbr  string `json:"abbr" db:"abbr"`
	//ComUsedAbbr *string `json:"com_used_abbr" db:"com_used_abbr"`
	Name       string `json:"name" db:"name"`
	Conference string `json:"conference" db:"conference"`
	Divisions  string `json:"divisions" db:"divisions"`
}
