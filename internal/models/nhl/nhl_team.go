package nhl

import "encoding/json"

type Team struct {
	Name string `json:"name"`
	Abbr string `json:"abbr"`
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

type TeamDB struct {
	ID           int    `db:"id"`
	ImgURL       string `db:"img_url" json:"img_url"`
	Name         string `db:"name" json:"name"`
	Abbreviation string `db:"abbreviation" json:"abbr"`
	Conference   string `db:"conference" json:"conference"`
	Division     string `db:"division" json:"division"`
}

type TeamRoster struct {
	Name        string       `json:"name"`
	Roster      []PlayerInfo `json:"roster"`
	PlayerCount int          `json:"player_count"`
}

type PlayerInfo struct {
	Number     string `json:"number" db:"number"`         // [0]
	Name       string `json:"name" db:"name"`             // [1]
	Surname    string `json:"surname" db:"surname"`       // [2]
	Position   string `json:"position" db:"position"`     // [3]
	Hand       string `json:"s/f" db:"s/f"`               // [4]
	Age        string `json:"age" db:"age"`               // [5]
	Acquired   string `json:"acquired" db:"acquired"`     // [6]
	Birthplace string `json:"birthplace" db:"birthplace"` // [7]
	Role       string `json:"role" db:"role"`             // [8]
	Injured    bool   `json:"injured" db:"injured"`       // find: img
}
type TeamRosterDB struct {
	TeamName string       `json:"name"`
	Players  []PlayerInfo `json:"roster"`
}
