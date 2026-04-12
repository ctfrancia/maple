package parser

// --- Domain models for scraping ---

// Tournament represents tournament metadata from chess-results.com
type Tournament struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Federation  string `json:"federation,omitempty"`
	Organizer   string `json:"organizer,omitempty"`
	Arbiter     string `json:"arbiter,omitempty"`
	Location    string `json:"location,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	Players     int    `json:"players"`
	Rounds      int    `json:"rounds"`
	TimeControl string `json:"time_control,omitempty"`
}

// Player represents a player entry in tournament standings
type Player struct {
	Rank        int       `json:"rank"`
	Title       string    `json:"title,omitempty"`
	Name        string    `json:"name"`
	FideID      string    `json:"fide_id,omitempty"`
	Federation  string    `json:"federation,omitempty"`
	Rating      int       `json:"rating"`
	Points      float64   `json:"points"`
	RatingPerf  int       `json:"rating_perf,omitempty"`
	TBBreakers  []float64 `json:"tb_breakers,omitempty"`
}

// RoundPairing represents a single pairing in a round
type RoundPairing struct {
	Board       int    `json:"board"`
	WhitePlayer string `json:"white_player"`
	WhiteRating int    `json:"white_rating"`
	Result      string `json:"result"`
	BlackPlayer string `json:"black_player"`
	BlackRating int    `json:"black_rating"`
}

// FederationTournamentEntry represents a tournament listing from a federation page
type FederationTournamentEntry struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Date     string `json:"date,omitempty"`
	Location string `json:"location,omitempty"`
}
