package mlb

type Team struct {
	Abbreviation string `json:"abbr"`
	Name         string `json:"name"`
	Division     string `json:"division"`
	Conference   string `json:"conference"`
}
