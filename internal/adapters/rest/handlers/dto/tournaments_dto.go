package dto

type Tournament struct {
	Id                      int        `json:"-"`
	TournamentId            int        `json:"tournament_id"`
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	Poster                  string     `json:"poster"` // url to the poster
	Schedules               []Schedule `json:"schedules"`
	StartDate               string     `json:"tournament_start_date"`
	EndDate                 string     `json:"tournament_end_date"`
	Location                string     `json:"location"`
	State                   string     `json:"state"` // open, closed, cancelled, completed
	PrizePool               int        `json:"prize_pool"`
	WinnerId                int        `json:"winner_id"`
	CreatorId               int        `json:"creator_id"`
	PlayerCount             int        `json:"player_count"`
	PlayerCountLimit        int        `json:"player_count_limit"`
	PlayerCountLimitEnabled bool       `json:"player_count_limit_enabled"` // if true, then player_count_limit is enforced
	AvailableSlots          int        `json:"available_slots"`
	RegistrationStartTime   string     `json:"registration_start_time"` // format: time.RFC3339
	RegistrationEndTime     string     `json:"registration_end_time"`   // format: time.RFC3339
	EnableRegistration      bool       `json:"enable_registration"`     // if true, then registration is enabled
	EnableRanking           bool       `json:"enable_ranking"`          // if true, then ranking is enabled
	RankingType             string     `json:"ranking_type"`            // glicko, elo
	CreatedAt               string     `json:"created_at"`
	UpdatedAt               string     `json:"updated_at"`
}

// Schedule represents a schedule for a tournament
// must have a length of atleast 1 for one day
type Schedule struct {
	DayOfWeek int    `json:"day_of_week"` // 0 = sunday, 1 = monday, etc
	StartTime string `json:"start_time"`  // required
	EndTime   string `json:"end_time"`    // if there is no end time then 23:59:59
}
