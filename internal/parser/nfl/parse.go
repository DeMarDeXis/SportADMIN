package nfl

import (
	"AdminAppForDiplom/internal/models/nfl"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"strings"
)

func ParseAbbr(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("tr").Each(func(i int, s *goquery.Selection) {
		// Buscamos la primera fila (encabezado) y la ignoramos
		if i == 0 {
			return
		}

		teamName := s.Find("a").Text()
		if teamName == "" {
			return
		}

		var abbrs []string
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				abbrs = append(abbrs, text)
			}
		})

		team := nfl.Team{
			Name: teamName,
		}

		if len(abbrs) >= 2 {
			team.Abbr = strings.TrimSpace(abbrs[1]) //Official
			if len(abbrs) >= 3 {
				team.ComUsedAbbr = strings.TrimSpace(abbrs[2]) //Commonly
			}
		}

		if team.Name != "" && team.Abbr != "" {
			g.Exports <- team
		}
	})
}
