package commands

type UpdateTournamentCommand struct {
	Name               string       `json:"name"`        //`json:"name" validate:"required,gte=3,lte=100"` look into this?
	Description        string       `json:"description"` // optional
	Schedule           []Schedule   `json:"schedule,omitempty"`
	AdditionalInfo     string       `json:"additional_info"`      // optional TODO: add this to the DTO
	LocationID         string       `json:"location_id"`          // need to revisit
	MaxPlayers         int          `json:"max_players"`          // optional when creating
	Contact            Contact      `json:"contact"`              // optional
	OpenToPublic       bool         `json:"open_to_public"`       // optional
	OpenToRegistration bool         `json:"open_to_registration"` // optional
	Registration       Registration `json:"registration"`         // optional
}
