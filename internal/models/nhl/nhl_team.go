package nhl

import "encoding/json"

type Team struct {
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}

type TeamDB struct {
	ID           int    `db:"id"`
	Name         string `db:"name" json:"name"`
	Abbreviation string `db:"abbreviation" json:"abbr"`
	ImgURL       string `db:"img_url"`
}

func (m Team) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string `json:"name"`
		Abbr string `json:"abbr"`
	}{
		Name: m.Name,
		Abbr: m.Abbr,
	})
}

type TeamRoster struct {
	Name   string       `json:"name"`
	Roster []PlayerInfo `json:"roster"`
}

type PlayerInfo struct {
	Number     string `json:"number"`     // [0]
	Name       string `json:"name"`       // [1]
	Surname    string `json:"surname"`    // [2]
	Position   string `json:"position"`   // [3]
	Hand       string `json:"s/f"`        // [4]
	Age        string `json:"age"`        // [5]
	Acquired   string `json:"acquired"`   // [6]
	Birthplace string `json:"birthplace"` // [7]
	Role       string `json:"role"`
	Injured    bool   `json:"injured"` // find: img
}
