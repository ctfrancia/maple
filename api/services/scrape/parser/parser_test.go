package parser

import (
	"testing"
)

const testStandingsHTML = `
<html>
<head><title>Chess-Results Server - Test Tournament 2024</title></head>
<body>
<table class="CRs1">
<tr>
<td>Rk.</td><td></td><td></td><td>Name</td><td>FideID</td><td>FED</td><td>Rtg</td><td>Pts.</td><td>TB1</td><td>TB2</td>
</tr>
<tr>
<td>1</td><td></td><td>GM</td><td>Adams, Michael</td><td>400041</td><td>ENG</td><td>2666</td><td>7.5</td><td>42.5</td><td>36.0</td>
</tr>
<tr>
<td>2</td><td></td><td>GM</td><td>Pert, Nicholas</td><td>403989</td><td>ENG</td><td>2536</td><td>7.0</td><td>40.0</td><td>34.5</td>
</tr>
<tr>
<td>3</td><td></td><td>IM</td><td>Fernandez, Daniel H</td><td>5801605</td><td>ENG</td><td>2518</td><td>6.5</td><td>38.5</td><td>33.0</td>
</tr>
</table>
</body>
</html>
`

const testTournamentInfoHTML = `
<html>
<head><title>Chess-Results Server - Copa Catalana 2024</title></head>
<body>
<table class="CRs1">
<tr><td>Tournament name</td><td>Copa Catalana 2024</td></tr>
<tr><td>Federation</td><td>ESP</td></tr>
<tr><td>Location</td><td>Barcelona, Spain</td></tr>
<tr><td>Date</td><td>15/02/2024 - 22/02/2024</td></tr>
<tr><td>Players</td><td>128</td></tr>
<tr><td>Rounds</td><td>9</td></tr>
<tr><td>Time control</td><td>90min + 30sec</td></tr>
<tr><td>Organizer</td><td>Federació Catalana d'Escacs</td></tr>
<tr><td>Arbiter</td><td>IA Garcia Lopez, Juan</td></tr>
</table>
</body>
</html>
`

const testFederationHTML = `
<html><body>
<table class="CRs1">
<tr><td>Name</td><td>Location</td><td>Date</td></tr>
<tr>
<td><a href="tnr123456.aspx?lan=1">Copa Catalana 2024</a></td>
<td>Barcelona, Spain</td><td>15/02/2024</td>
</tr>
<tr>
<td><a href="tnr789012.aspx?lan=1">Open Sants 2024</a></td>
<td>Barcelona, Spain</td><td>22/08/2024</td>
</tr>
</table>
</body></html>
`

func TestParseStandings(t *testing.T) {
	players, err := ParseStandings([]byte(testStandingsHTML))
	if err != nil {
		t.Fatalf("ParseStandings: %v", err)
	}
	if len(players) != 3 {
		t.Fatalf("expected 3 players, got %d", len(players))
	}

	p := players[0]
	assertEqual(t, "rank", p.Rank, 1)
	assertEqual(t, "name", p.Name, "Adams, Michael")
	assertEqual(t, "fide_id", p.FideID, "400041")
	assertEqual(t, "fed", p.Federation, "ENG")
	assertEqual(t, "rating", p.Rating, 2666)
	assertEqualFloat(t, "points", p.Points, 7.5)

	if len(p.TBBreakers) < 2 {
		t.Fatalf("expected >=2 tiebreakers, got %d", len(p.TBBreakers))
	}
	assertEqualFloat(t, "tb1", p.TBBreakers[0], 42.5)
}

func TestParseTournamentInfo(t *testing.T) {
	info, err := ParseTournamentInfo([]byte(testTournamentInfoHTML), "999999")
	if err != nil {
		t.Fatalf("ParseTournamentInfo: %v", err)
	}
	assertEqual(t, "name", info.Name, "Copa Catalana 2024")
	assertEqual(t, "federation", info.Federation, "ESP")
	assertEqual(t, "location", info.Location, "Barcelona, Spain")
	assertEqual(t, "players", info.Players, 128)
	assertEqual(t, "rounds", info.Rounds, 9)
	assertEqual(t, "time_control", info.TimeControl, "90min + 30sec")
}

func TestParseFederationTournaments(t *testing.T) {
	tournaments, err := ParseFederationTournaments([]byte(testFederationHTML))
	if err != nil {
		t.Fatalf("ParseFederationTournaments: %v", err)
	}
	if len(tournaments) != 2 {
		t.Fatalf("expected 2 tournaments, got %d", len(tournaments))
	}
	assertEqual(t, "id[0]", tournaments[0].ID, "123456")
	assertEqual(t, "name[0]", tournaments[0].Name, "Copa Catalana 2024")
	assertEqual(t, "id[1]", tournaments[1].ID, "789012")
}

func TestExtractInt(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"2666", 2666}, {"128 players", 128}, {" 42 ", 42}, {"", 0},
	}
	for _, c := range cases {
		if got := extractInt(c.in); got != c.want {
			t.Errorf("extractInt(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestExtractFloat(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"7.5", 7.5}, {"7½", 7.5}, {"42,5", 42.5}, {"", 0},
	}
	for _, c := range cases {
		if got := extractFloat(c.in); got != c.want {
			t.Errorf("extractFloat(%q) = %f, want %f", c.in, got, c.want)
		}
	}
}

func assertEqual[T comparable](t *testing.T, field string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", field, got, want)
	}
}

func assertEqualFloat(t *testing.T, field string, got, want float64) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %f, want %f", field, got, want)
	}
}
