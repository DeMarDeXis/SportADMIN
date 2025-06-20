package nba

import (
	"AdminAppForDiplom/internal/models/nba"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"strings"
)

func NameParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("tr").Each(func(i int, s *goquery.Selection) {
		teamName := s.Find("td a").Text()
		abbr := strings.TrimSpace(s.Find("td").First().Text())

		if teamName != "" && abbr != "" {
			teams := nba.Team{
				Name:         teamName,
				Abbreviation: abbr,
			}
			g.Exports <- teams
		}
	})
}
