//go:build integration

package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ctfrancia/maple/api/services/scrape/client"
	"github.com/ctfrancia/maple/api/services/scrape/parser"
)

func newClient() *client.Client {
	return client.New(
		client.WithDelay(1*time.Second),
		client.WithLogger(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))),
	)
}

func printJSON(t *testing.T, v any) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	fmt.Println(string(b))
}

// TestLiveFederation fetches the tournament list for a federation.
// Override the federation with TEST_FED env var (default: ESP).
//
//	go test -tags=integration -v -run TestLiveFederation ./api/services/scrape/
func TestLiveFederation(t *testing.T) {
	fed := os.Getenv("TEST_FED")
	if fed == "" {
		fed = "ESP"
	}

	c := newClient()
	html, err := c.Get(context.Background(), client.FederationURL(fed))
	if err != nil {
		t.Fatalf("fetching federation %s: %v", fed, err)
	}

	tournaments, err := parser.ParseFederationTournaments(html)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}

	t.Logf("federation %s: %d tournaments found", fed, len(tournaments))

	n := min(len(tournaments), 10)
	printJSON(t, tournaments[:n])
}

// TestLiveTournamentInfo fetches and parses a single tournament's info page.
// Requires TEST_TNR_ID env var (e.g. TEST_TNR_ID=1234567).
//
//	go test -tags=integration -v -run TestLiveTournamentInfo ./api/services/scrape/
func TestLiveTournamentInfo(t *testing.T) {
	id := os.Getenv("TEST_TNR_ID")
	if id == "" {
		t.Skip("set TEST_TNR_ID to a chess-results.com tournament ID")
	}

	c := newClient()
	html, err := c.Get(context.Background(), client.TournamentURL(id, 0))
	if err != nil {
		t.Fatalf("fetching tournament %s: %v", id, err)
	}

	info, err := parser.ParseTournamentInfo(html, id)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}

	printJSON(t, info)
}

// TestLiveStandings fetches and parses final standings for a tournament.
// Requires TEST_TNR_ID env var.
//
//	go test -tags=integration -v -run TestLiveStandings ./api/services/scrape/
func TestLiveStandings(t *testing.T) {
	id := os.Getenv("TEST_TNR_ID")
	if id == "" {
		t.Skip("set TEST_TNR_ID to a chess-results.com tournament ID")
	}

	c := newClient()
	html, err := c.Get(context.Background(), client.TournamentURL(id, 0))
	if err != nil {
		t.Fatalf("fetching standings for %s: %v", id, err)
	}

	players, err := parser.ParseStandings(html)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}

	t.Logf("tournament %s: %d players in standings", id, len(players))
	n := min(len(players), 20)
	printJSON(t, players[:n])
}

// TestLiveRoundPairings fetches and parses pairings for a specific round.
// Requires TEST_TNR_ID. Optionally set TEST_ROUND (default: 1).
//
//	go test -tags=integration -v -run TestLiveRoundPairings ./api/services/scrape/
func TestLiveRoundPairings(t *testing.T) {
	id := os.Getenv("TEST_TNR_ID")
	if id == "" {
		t.Skip("set TEST_TNR_ID to a chess-results.com tournament ID")
	}

	round := os.Getenv("TEST_ROUND")
	if round == "" {
		round = "1"
	}

	c := newClient()
	// art=4 is round pairings view, rd=N selects the round
	path := fmt.Sprintf("/tnr%s.aspx?lan=1&art=4&rd=%s", id, round)
	html, err := c.Get(context.Background(), path)
	if err != nil {
		t.Fatalf("fetching round %s pairings for %s: %v", round, id, err)
	}

	pairings, err := parser.ParseRoundPairings(html)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}

	t.Logf("tournament %s round %s: %d pairings", id, round, len(pairings))
	n := min(len(pairings), 20)
	printJSON(t, pairings[:n])
}

// TestLiveAll runs info + standings + all rounds for a tournament in one shot.
// Requires TEST_TNR_ID.
//
//	go test -tags=integration -v -run TestLiveAll ./api/services/scrape/
func TestLiveAll(t *testing.T) {
	id := os.Getenv("TEST_TNR_ID")
	if id == "" {
		t.Skip("set TEST_TNR_ID to a chess-results.com tournament ID")
	}

	c := newClient()
	ctx := context.Background()

	html, err := c.Get(ctx, client.TournamentURL(id, 0))
	if err != nil {
		t.Fatalf("fetching info: %v", err)
	}

	info, err := parser.ParseTournamentInfo(html, id)
	if err != nil {
		t.Fatalf("parsing info: %v", err)
	}
	fmt.Println("=== Tournament Info ===")
	printJSON(t, info)

	players, err := parser.ParseStandings(html)
	if err != nil {
		t.Fatalf("parsing standings: %v", err)
	}
	fmt.Printf("\n=== Standings (%d players, showing first 10) ===\n", len(players))
	n := min(len(players), 10)
	printJSON(t, players[:n])

	if info.Rounds > 0 {
		fmt.Printf("\n=== Round 1 Pairings ===\n")
		path := fmt.Sprintf("/tnr%s.aspx?lan=1&art=4&rd=1", id)
		rhtml, err := c.Get(ctx, path)
		if err != nil {
			t.Logf("warning: could not fetch round 1: %v", err)
		} else {
			pairings, err := parser.ParseRoundPairings(rhtml)
			if err != nil {
				t.Logf("warning: could not parse round 1 pairings: %v", err)
			} else {
				n := min(len(pairings), 10)
				printJSON(t, pairings[:n])
			}
		}
	}
}
