// Package mappers is used to map the DTOs to the domain models
package mappers

import (
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto"
	"github.com/ctfrancia/maple/internal/core/domain"
)

func TransformNewAPIConsumerRequestToDomainModel(requestBody dto.NewAPIConsumerRequest) domain.NewAPIConsumer {
	return domain.NewAPIConsumer{
		FirstName:       requestBody.FirstName,
		LastName:        requestBody.LastName,
		Username:        requestBody.Username,
		Email:           requestBody.Email,
		Password:        requestBody.Password,
		Website:         requestBody.Website,
		ClubAffiliation: requestBody.ClubAffiliation,
	}
}
