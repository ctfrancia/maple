package domain

// Tournament represents a tournament
type Tournament struct {
	Id                      int
	TournamentId            int
	Name                    string
	Description             string
	Poster                  string
	Schedules               []Schedule
	StartDate               string
	EndDate                 string
	Location                string
	State                   string // open, closed, cancelled, completed
	PrizePool               int
	WinnerId                int
	CreatorId               int
	PlayerCount             int
	PlayerCountLimit        int
	PlayerCountLimitEnabled bool // if true, then player_count_limit is enforced
	AvailableSlots          int
	RegistrationStartTime   string // format: time.RFC3339
	RegistrationEndTime     string // format: time.RFC3339
	EnableRegistration      bool   // if true, then registration is enabled
	EnableRanking           bool   // if true, then ranking is enabled
	RankingType             string // glicko, elo
	CreatedAt               string
	UpdatedAt               string
}

// Schedule represents a schedule for a tournament
// must have a length of atleast 1 for one day
type Schedule struct {
	DayOfWeek int    // 0 = sunday, 1 = monday, etc
	StartTime string // required
	EndTime   string // if there is no end time then 23:59:59
}
