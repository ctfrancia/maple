package tournamenthandlers

import (
	dto "github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto/tournament"
	"github.com/ctfrancia/maple/internal/core/domain"
)

// mapDomainToTournament converts dto.CreateTournamentRequest to domain.Tournament
func mapTournamentToDomain(t dto.CreateTournamentRequest) domain.Tournament {
	return domain.Tournament{
		Name: t.Name,
	}
}
