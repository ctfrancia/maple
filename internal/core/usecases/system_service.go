package usecases

import (
	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type SystemHealthUseCase struct {
	sAdapter ports.SystemAdapter
	repo     ports.SystemRepository
	security ports.SecurityAdapter
}

func NewSystemHealthUseCase(sa ports.SystemAdapter, sr ports.SystemRepository, sec ports.SecurityAdapter) ports.SystemServicer {
	return &SystemHealthUseCase{
		sAdapter: sa,
		repo:     sr,
		security: sec,
	}
}

func (shs *SystemHealthUseCase) ProcessSystemHealthRequest() domain.System {
	return shs.sAdapter.GetSystemInfo()
}

func (shs *SystemHealthUseCase) Login(username, password string) (any, error) {
	return nil, nil
}

func (shs *SystemHealthUseCase) NewAPIConsumer(consumer domain.NewAPIConsumer) (domain.NewAPIConsumer, error) {
	err := shs.repo.SelectByEmail(consumer)
	if err != nil {
		return domain.NewAPIConsumer{}, err
	}

	// Generate password
	//	generatedPassword, err := shs.security.CreateSecretKey(security.PasswordGeneratorDefaultLength)
	generatedPassword, err := shs.security.CreateSecretKey(domain.PasswordGeneratorDefaultLength)
	if err != nil {
		return domain.NewAPIConsumer{}, err
	}

	// Hash the password
	hashedPassword, err := shs.security.Hash(generatedPassword)
	if err != nil {
		return domain.NewAPIConsumer{}, err
	}

	// Create consumer for DB with hashed password
	consumerForDB := consumer
	consumerForDB.Password = hashedPassword

	err = shs.repo.CreateNewConsumer(consumerForDB)
	if err != nil {
		return domain.NewAPIConsumer{}, err
	}

	// Return consumer with plain password for client
	consumer.Password = generatedPassword
	return consumer, nil
}

func (shs *SystemHealthUseCase) CreateNewConsumer(consumer domain.NewAPIConsumer) (domain.NewAPIConsumer, error) {
	return domain.NewAPIConsumer{}, nil
}
