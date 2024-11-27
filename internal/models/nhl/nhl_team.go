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

type TeamRoster struct {
	Name   string       `json:"name"`
	Roster []PlayerInfo `json:"roster"`
}

type PlayerInfo struct {
	Number     string `json:"number"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Position   string `json:"position"`
	Hand       string `json:"s/f"`
	Age        string `json:"age"`
	Acquired   string `json:"acquired"`
	Birthplace string `json:"birthplace"`
}
