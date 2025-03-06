package nhl

import (
	"AdminAppForDiplom/internal/models/nhl"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"strconv"
	"strings"
)

func NameParse(g *geziyor.Geziyor, r *client.Response) {
	count := 0
	r.HTMLDoc.Find("li").Each(func(i int, s *goquery.Selection) {
		if count >= 32 {
			return
		}

		abbr := strings.TrimSpace(strings.Split(s.Text(), "–")[0])
		if len(abbr) <= 3 {
			teamName := s.Find("a").Text()

			teams := nhl.Team{
				Name: teamName,
				Abbr: abbr,
			}

			g.Exports <- teams
			count++
		}
	})
}

// TODO: add and fix roster parse
func RosterParse(g *geziyor.Geziyor, r *client.Response) {
	count := 0
	r.HTMLDoc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if count >= 29 && count <= 53 {

			//text := strings.TrimSpace(strings.Split(s.Text(), "–")[0])
			text := strings.TrimSpace(s.Text())
			parts := strings.Fields(text)

			fmt.Println(count)
			fmt.Println()
			fmt.Println(text)
			//fmt.Println(name)

			if len(parts) >= 8 {
				player := nhl.PlayerInfo{
					Number:     parts[0],
					Name:       parts[1],
					Surname:    parts[2],
					Position:   parts[3],
					Hand:       parts[4],
					Age:        parts[5],
					Acquired:   parts[6],
					Birthplace: strings.Join(parts[7:], " "), // Join remaining parts for birthplace
				}

				g.Exports <- player
			}
		}
		count++
	})
}

func AllRosterParse(g *geziyor.Geziyor, r *client.Response) {
	processedTeams := make(map[string]bool)

	r.HTMLDoc.Find(".navbox-list").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(j int, link *goquery.Selection) {
			teamName := strings.TrimSpace(link.Text())

			// Skip if already processed or not a valid team name
			if processedTeams[teamName] ||
				strings.Contains(teamName, "Conference") ||
				strings.Contains(teamName, "Division") {
				return
			}

			currentTeam := &nhl.TeamRoster{
				Name:        teamName,
				Roster:      make([]nhl.PlayerInfo, 0),
				PlayerCount: 0,
			}
			processedTeams[teamName] = true
			//g.Exports <- currentTeam

			// Process team roster
			if href, ok := link.Attr("href"); ok {
				teamURL := r.JoinURL(href)
				g.Get(teamURL, func(g *geziyor.Geziyor, r *client.Response) {
					unoRosterParse(g, r, currentTeam)
				})
			}
		})
	})
}

func unoRosterParse(g *geziyor.Geziyor, r *client.Response, team *nhl.TeamRoster) {
	count := 0
	r.HTMLDoc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if count >= 29 && count <= 53 {
			text := strings.TrimSpace(s.Text())
			parts := strings.Fields(text)
			fmt.Println(len(parts))

			if len(parts) >= 8 && parts[0] != " " {
				// Validation checks
				isInj := isInjured(s)

				if isValidPlayerData(parts) && isCaptainOrAssistant(parts) == false {
					player := nhl.PlayerInfo{
						Number:     parts[0],
						Name:       parts[1],
						Surname:    parts[2],
						Position:   parts[3],
						Hand:       parts[4],
						Age:        parts[5],
						Acquired:   parts[6],
						Birthplace: strings.Join(parts[7:], " "),
						Injured:    isInj,
						Role:       " ",
					}
					team.Roster = append(team.Roster, player)
				} else if isCaptainOrAssistant(parts) {
					player := nhl.PlayerInfo{
						Number:     parts[0],
						Name:       parts[1],
						Surname:    parts[2],
						Role:       parts[3],
						Position:   parts[4],
						Hand:       parts[5],
						Age:        parts[6],
						Acquired:   parts[7],
						Birthplace: strings.Join(parts[8:], " "),
					}
					team.Roster = append(team.Roster, player)
				}
			}
		}
		count++
	})

	// Only export team if roster is not empty
	if len(team.Roster) > 0 {
		team.PlayerCount = len(team.Roster)
		g.Exports <- team
	}
}

func isValidPlayerData(parts []string) bool {
	// Number check - must be numeric
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return false
	}

	// Position check - max 2 characters
	if len(parts[3]) > 2 {
		return false
	}

	// Hand (s/f) check - max 1 character
	if len(parts[4]) > 1 {
		return false
	}

	// Age check - must be numeric
	if _, err := strconv.Atoi(parts[5]); err != nil {
		return false
	}

	// Additional validation can be added here
	return true
}

func isCaptainOrAssistant(parts []string) bool {
	if parts[3] == "(C)" || parts[3] == "(A)" {
		return true
	}

	return false
}

func isInjured(s *goquery.Selection) bool {
	if s.Find("img[alt='Injured Reserve']").Length() > 0 {
		return true
	}

	return false
}
