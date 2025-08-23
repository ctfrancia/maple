package tournamenthandlers

import (
	dto "github.com/ctfrancia/maple/internal/adapters/http/handlers/dto/tournament"
	commands "github.com/ctfrancia/maple/internal/application/commands/tournament"
	"github.com/ctfrancia/maple/internal/core/domain"
)

type TournamentMapper struct{}

func NewTournamentMapper() TournamentMapper {
	return TournamentMapper{}
}

func (m TournamentMapper) MapToCommand(dto dto.CreateTournamentRequest) commands.CreateTournamentCommand {
	return commands.CreateTournamentCommand{
		Name:     dto.Name,
		Schedule: mapScheduleToCommand(dto.Schedule),
	}
}

func mapScheduleToCommand(sch []dto.Schedule) []commands.Schedule {
	xSch := make([]commands.Schedule, len(sch))
	for i, s := range xSch {
		xSch[i] = commands.Schedule{
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
		}
	}
	return xSch
}

// mapDomainToTournament converts dto.CreateTournamentRequest to domain.Tournament
func mapTournamentToDomain(t dto.CreateTournamentRequest) domain.Tournament {
	return domain.Tournament{
		Name: t.Name,
	}
}

func mapTournamentToDto(t domain.Tournament) dto.TournamentResponse {
	return dto.TournamentResponse{
		ID:                 t.PublicID.String(),
		Name:               t.Name,
		Description:        t.Description,
		Location:           mapLocationToDto(t.Location),
		OpenToPublic:       t.OpenToPublic,
		OpenToSpectators:   t.OpenToSpectators,
		OpenToRegistration: t.OpenToRegistration,
		Registration:       mapRegistrationToDto(t.Registration),
		Arbitrator:         t.Arbitrator,
		Matches:            nil,
		Players:            nil,
		NumberOfPlayers:    t.NumberOfPlayers,
		Schedule:           mapScheduleToDto(t.Schedule),
		Results:            nil,
		Status:             dto.TournamentStatus(t.Status),
	}
}

func mapLocationToDto(l domain.Location) dto.Location {
	return dto.Location{
		Address:    l.Address,
		PostalCode: l.PostalCode,
		City:       l.City,
		State:      l.State,
		Country:    l.Country,
		// Name:       l.Name,
		// County:     l.County,
		// Province:   l.Province,
		// Latitude:   l.Latitude,
		// Longitude:  l.Longitude,
		// Timezone: domain.TimezoneUTC,
	}
}

func mapResultsToDto(r []domain.Result) []dto.Result {
	xResults := make([]dto.Result, len(r))
	for i, s := range r {
		xResults[i] = dto.Result{
			Player: s.Player.PublicID.String(),
			Prize:  s.Prize,
		}
	}
	return xResults
}

func mapScheduleToDto(sch []domain.Schedule) []dto.Schedule {
	xSch := make([]dto.Schedule, len(sch))
	for i, s := range xSch {
		xSch[i] = dto.Schedule{
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
		}
	}
	return xSch
}

func mapRegistrationToDto(r domain.Registration) dto.Registration {
	return dto.Registration{
		Status:    dto.RegistrationStatus(r.Status),
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
		Fee:       r.Fee,
		PrizePool: r.PrizePool,
		Payment:   mapRegistrationPayoutToDto(r.Payment),
	}
}

func mapRegistrationPayoutToDto(p []domain.Payment) []dto.Payment {
	xPayout := make([]dto.Payment, len(p))
	for i, s := range xPayout {
		xPayout[i] = dto.Payment{
			Place:  s.Place,
			Amount: s.Amount,
			Other:  s.Other,
		}
	}
	return xPayout
}
