package mlb

import (
	"AdminAppForDiplom/internal/models/mlb"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"strings"
)

func NameParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {

		if i == 0 {
			return
		}

		conference := "American League"
		if i < 15 {
			conference = "National League"
		}

		cells := s.Find("td")

		if cells.Length() >= 3 {
			teamName := strings.TrimSpace(cells.Eq(0).Text())
			abbreviation := strings.TrimSpace(cells.Eq(1).Text())
			division := strings.TrimSpace(cells.Eq(2).Text())
			fmt.Println(teamName, abbreviation, division, conference)

			if teamName != "" && abbreviation != "" {
				team := mlb.TeamParse{
					Name:         teamName,
					Abbreviation: abbreviation,
					Division:     division,
					Conference:   conference,
				}
				g.Exports <- team
			}
		}
	})
}
