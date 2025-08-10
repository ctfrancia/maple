package repository

import (
	"context"

	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type TournamentRepository struct {
	db any // TODO: add the database connection it will be a gorm.DB
}

func NewTournamentRepository() ports.TournamentRepository {
	return &TournamentRepository{}
}

func (tr *TournamentRepository) CreateTournament(tournament domain.Tournament) error {
	return nil
}

func (tr *TournamentRepository) GetTournaments(ctx context.Context, page int, size int) ([]domain.Tournament, error) {
	return nil, nil
}

func (tr *TournamentRepository) GetTournament(id int) (domain.Tournament, error) {
	return domain.Tournament{}, nil
}

func (tr *TournamentRepository) UpdateTournament(tournament domain.Tournament) error {
	return nil
}

func (tr *TournamentRepository) DeleteTournament(id int) error {
	return nil
}
