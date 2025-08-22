// Package dto is the data transfer object for the REST API
package dto

type CreateTournamentRequest struct {
	Name string `json:"name"`
}
