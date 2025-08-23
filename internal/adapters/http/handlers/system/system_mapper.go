package systemhandlers

import (
	"github.com/ctfrancia/maple/internal/adapters/http/handlers/dto/system"
	"github.com/ctfrancia/maple/internal/core/domain"
)

func transformNewAPIConsumerRequestToDomainModel(requestBody dto.NewAPIConsumerRequest) domain.NewAPIConsumer {
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
