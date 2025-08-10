package usecases

import (
	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type TournamentUseCase struct {
}

func NewTournamentUseCase() ports.TournamentUseCase {
	return &TournamentUseCase{}
}

func (tuc *TournamentUseCase) ProcessTournamentRequest() ([]domain.Tournament, error) {
	return []domain.Tournament{}, nil
}
