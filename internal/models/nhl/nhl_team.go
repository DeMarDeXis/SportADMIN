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
