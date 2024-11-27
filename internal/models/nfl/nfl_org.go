package nfl

type Team struct {
	Name        string `json:"name"`
	Abbr        string `json:"abbr"`
	ComUsedAbbr string `json:"com_used_abbr"`
}
