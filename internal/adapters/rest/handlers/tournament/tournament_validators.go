package tournamenthandlers

import (
	dto "github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto/tournament"
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/validator"
)

// this file contains handler specific logic

// isValidRequest checks if the request is valid, TODO: I don't think that
// I am doing this the best way, but for now it works
func isValidRequest(t dto.CreateTournamentRequest) (bool, map[string]string) {
	v := validator.NewValidator()
	v.Check(t.Name != "", "name", "name is required")
	if !v.Valid() {
		/*
			errMsg := "invalid request: "
			for key, message := range v.ReturnErrors() {
				errMsg += key + ": " + message + ", "
			}
		*/

		return false, v.ReturnErrors()
	}
	return true, nil
}
