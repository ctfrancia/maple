package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	tnrIDRegex    = regexp.MustCompile(`tnr(\d+)\.aspx`)
	whitespaceRe  = regexp.MustCompile(`\s+`)
	playerCountRe = regexp.MustCompile(`(\d+)\s*(?:players?|participants?)`)
	roundCountRe  = regexp.MustCompile(`(\d+)\s*(?:rounds?)`)
)

// ParseTournamentInfo extracts tournament metadata from the main tournament page.
func ParseTournamentInfo(html []byte, tournamentID string) (*Tournament, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	t := &Tournament{
		ID:  tournamentID,
		URL: fmt.Sprintf("https://chess-results.com/tnr%s.aspx?lan=1", tournamentID),
	}

	title := doc.Find("title").Text()
	if idx := strings.Index(title, " - "); idx != -1 {
		t.Name = strings.TrimSpace(title[idx+3:])
	} else {
		t.Name = strings.TrimSpace(title)
	}

	doc.Find("table.CRs1 tr, table.CRs2 tr, div.defaultDialog tr").Each(func(_ int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() < 2 {
			return
		}
		label := normalizeText(cells.First().Text())
		value := normalizeText(cells.Last().Text())

		switch {
		case containsAny(label, "tournament name", "nombre del torneo", "turniername"):
			if t.Name == "" || t.Name == tournamentID {
				t.Name = value
			}
		case containsAny(label, "federation", "federaci"):
			t.Federation = value
		case containsAny(label, "organizer", "organizador"):
			t.Organizer = value
		case containsAny(label, "arbiter", "árbitro"):
			t.Arbiter = value
		case containsAny(label, "location", "lugar"):
			t.Location = value
		case containsAny(label, "date", "fecha"):
			t.StartDate = value
		case containsAny(label, "players", "jugadores", "teilnehmer"):
			t.Players = extractInt(value)
		case containsAny(label, "rounds", "rondas", "runden"):
			t.Rounds = extractInt(value)
		case containsAny(label, "time control", "ritmo de juego"):
			t.TimeControl = value
		}
	})

	pageText := doc.Text()
	if t.Players == 0 {
		if m := playerCountRe.FindStringSubmatch(pageText); len(m) > 1 {
			t.Players = extractInt(m[1])
		}
	}
	if t.Rounds == 0 {
		if m := roundCountRe.FindStringSubmatch(pageText); len(m) > 1 {
			t.Rounds = extractInt(m[1])
		}
	}

	return t, nil
}

// ParseStandings extracts the final standings (art=0) or starting rank (art=1).
func ParseStandings(html []byte) ([]Player, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var players []Player

	var headers []string
	doc.Find("table.CRs1 tr").First().Find("th, td").Each(func(_ int, s *goquery.Selection) {
		headers = append(headers, normalizeText(s.Text()))
	})
	if len(headers) == 0 {
		doc.Find("table.CRs2 tr").First().Find("th, td").Each(func(_ int, s *goquery.Selection) {
			headers = append(headers, normalizeText(s.Text()))
		})
	}

	colIdx := mapColumns(headers)

	doc.Find("table.CRs1 tr, table.CRs2 tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		cells := s.Find("td")
		if cells.Length() < 3 {
			return
		}

		p := Player{}
		cellTexts := make([]string, cells.Length())
		cells.Each(func(j int, cell *goquery.Selection) {
			cellTexts[j] = normalizeText(cell.Text())
		})

		if idx, ok := colIdx["rank"]; ok && idx < len(cellTexts) {
			p.Rank = extractInt(cellTexts[idx])
		}
		if idx, ok := colIdx["title"]; ok && idx < len(cellTexts) {
			p.Title = cellTexts[idx]
		}
		if idx, ok := colIdx["name"]; ok && idx < len(cellTexts) {
			p.Name = cellTexts[idx]
		}
		if idx, ok := colIdx["fideid"]; ok && idx < len(cellTexts) {
			p.FideID = cellTexts[idx]
			if p.FideID == "" {
				cells.Eq(idx).Find("a").Each(func(_ int, a *goquery.Selection) {
					if href, exists := a.Attr("href"); exists && strings.Contains(href, "fide.com") {
						parts := strings.Split(href, "/")
						if len(parts) > 0 {
							p.FideID = parts[len(parts)-1]
						}
					}
				})
			}
		}
		if idx, ok := colIdx["fed"]; ok && idx < len(cellTexts) {
			p.Federation = cellTexts[idx]
		}
		if idx, ok := colIdx["rating"]; ok && idx < len(cellTexts) {
			p.Rating = extractInt(cellTexts[idx])
		}
		if idx, ok := colIdx["points"]; ok && idx < len(cellTexts) {
			p.Points = extractFloat(cellTexts[idx])
		}
		if idx, ok := colIdx["ratingperf"]; ok && idx < len(cellTexts) {
			p.RatingPerf = extractInt(cellTexts[idx])
		}

		for tbIdx := 1; tbIdx <= 5; tbIdx++ {
			key := fmt.Sprintf("tb%d", tbIdx)
			if idx, ok := colIdx[key]; ok && idx < len(cellTexts) {
				p.TBBreakers = append(p.TBBreakers, extractFloat(cellTexts[idx]))
			}
		}

		if p.Name != "" {
			players = append(players, p)
		}
	})

	return players, nil
}

// ParseRoundPairings extracts pairings from a round view (art=4&rd=N).
func ParseRoundPairings(html []byte) ([]RoundPairing, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var pairings []RoundPairing
	doc.Find("table.CRs1 tr, table.CRs2 tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		cells := s.Find("td")
		if cells.Length() < 5 {
			return
		}
		cellTexts := make([]string, cells.Length())
		cells.Each(func(j int, cell *goquery.Selection) {
			cellTexts[j] = normalizeText(cell.Text())
		})

		p := RoundPairing{}
		if len(cellTexts) >= 6 {
			p.Board = extractInt(cellTexts[0])
			p.WhitePlayer = cellTexts[1]
			p.WhiteRating = extractInt(cellTexts[2])
			p.Result = normalizeResult(cellTexts[3])
			p.BlackPlayer = cellTexts[4]
			p.BlackRating = extractInt(cellTexts[5])
		}
		if p.WhitePlayer != "" || p.BlackPlayer != "" {
			pairings = append(pairings, p)
		}
	})

	return pairings, nil
}

// ParseFederationTournaments extracts tournament listings from a federation page.
func ParseFederationTournaments(html []byte) ([]FederationTournamentEntry, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var tournaments []FederationTournamentEntry
	doc.Find("table.CRs1 tr, table.CRs2 tr, table tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		s.Find("a[href]").Each(func(_ int, a *goquery.Selection) {
			href, exists := a.Attr("href")
			if !exists {
				return
			}
			matches := tnrIDRegex.FindStringSubmatch(href)
			if len(matches) < 2 {
				return
			}
			entry := FederationTournamentEntry{
				ID:   matches[1],
				Name: normalizeText(a.Text()),
				URL:  fmt.Sprintf("https://chess-results.com/tnr%s.aspx?lan=1", matches[1]),
			}
			cells := s.Find("td")
			cellTexts := make([]string, cells.Length())
			cells.Each(func(j int, cell *goquery.Selection) {
				cellTexts[j] = normalizeText(cell.Text())
			})
			for _, ct := range cellTexts {
				if looksLikeDate(ct) && entry.Date == "" {
					entry.Date = ct
				} else if looksLikeLocation(ct) && entry.Location == "" && ct != entry.Name {
					entry.Location = ct
				}
			}
			if entry.Name != "" {
				tournaments = append(tournaments, entry)
			}
		})
	})

	return tournaments, nil
}

func mapColumns(headers []string) map[string]int {
	m := make(map[string]int)
	tbCount := 0
	for i, h := range headers {
		h = strings.ToLower(h)
		switch {
		case containsAny(h, "no.", "rk.", "rank", "nr"):
			m["rank"] = i
		case h == "title" || h == "ti." || h == "tit":
			m["title"] = i
		case containsAny(h, "name", "nombre", "spieler"):
			m["name"] = i
		case containsAny(h, "fideid", "fide id", "id"):
			m["fideid"] = i
		case containsAny(h, "fed", "pais"):
			m["fed"] = i
		case containsAny(h, "rtg", "elo", "rating") && !strings.Contains(h, "perf"):
			m["rating"] = i
		case containsAny(h, "pts.", "pts", "points", "puntos"):
			m["points"] = i
		case containsAny(h, "rp", "perf"):
			m["ratingperf"] = i
		case containsAny(h, "tb", "bu.", "buc", "sb"):
			tbCount++
			m[fmt.Sprintf("tb%d", tbCount)] = i
		}
	}
	return m
}

func normalizeText(s string) string {
	return whitespaceRe.ReplaceAllString(strings.TrimSpace(s), " ")
}

func normalizeResult(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "½", "1/2")
}

func extractInt(s string) int {
	s = strings.TrimSpace(s)
	cleaned := ""
	for i, c := range s {
		if c >= '0' && c <= '9' {
			cleaned += string(c)
		} else if c == '-' && i == 0 {
			cleaned += string(c)
		}
	}
	n, _ := strconv.Atoi(cleaned)
	return n
}

func extractFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "½", ".5")
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func containsAny(s string, substrs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range substrs {
		if strings.Contains(lower, sub) {
			return true
		}
	}
	return false
}

func looksLikeDate(s string) bool {
	return regexp.MustCompile(`\d{2}[\./]\d{2}[\./]\d{4}`).MatchString(s) ||
		regexp.MustCompile(`\d{4}-\d{2}-\d{2}`).MatchString(s) ||
		containsAny(s, "jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec")
}

func looksLikeLocation(s string) bool {
	return strings.Contains(s, ",") || len(s) > 3
}
