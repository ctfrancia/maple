package tournamenthandlers

import (
	dto "github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto/tournament"
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/validator"
)

// this file contains handler specific logic

func isValidRequest(t dto.CreateTournamentRequest) (bool, error) {
	v := validator.NewValidator()
	v.Check(t.Name != "", "name", "name is required")
	return true, nil
}
