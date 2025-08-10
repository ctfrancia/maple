// package system handlers within this package are for the system, this means the apiconsumers
// the ones that can consume the api
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto"
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/mappers"
	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/validator"
	"github.com/ctfrancia/maple/internal/adapters/rest/response"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type SystemHealthHandler struct {
	system    ports.SystemServicer
	response  ports.ResponseHelper
	logger    ports.Logger
	validator ports.ValidatorServicer
	service   ports.SystemServicer
}

func NewSystemHandler(log ports.Logger, ss ports.SystemServicer) ports.SystemHandler {
	handler := &SystemHealthHandler{
		system:    ss,
		response:  response.NewResponseWriter(log),
		logger:    log,
		validator: validator.NewValidator(),
	}

	return handler
}

func (h *SystemHealthHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	// start of the business logic
	sysInfo := h.system.ProcessSystemHealthRequest()

	// end of the business logic
	res, err := json.Marshal(sysInfo)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	// write the response
	h.response.WriteJSON(w, http.StatusOK, res, nil)
}

// LoginHandler handles the login request for logging into the system this LoginHandler is for
// logging in as a API consumer, which will have access to the APIs of Maple
func (h *SystemHealthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// get the request body
	var requestBody dto.SystemLoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		h.response.BadRequestResponse(w, r, err)
		return
	}
	// TODO: add validation for the request body here it will also transform to the domain model

	// Pass request to the service layer
	resp, err := h.system.Login(requestBody.Username, requestBody.Password)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	// TODO: this request needs to be sent to the service/domain layer

	h.logger.Info(r.Context(), "login attempt", ports.String("username", requestBody.Username))

	// FIXME: write headers this should be done in the response package
	headers := http.Header{"Content-Type": []string{"application/json"}}

	h.response.WriteJSON(w, http.StatusOK, resp, headers)
}

// NewConsumerHandler handles the request for creating a new API consumer that will be
// able to use the APIs of Maple
func (h *SystemHealthHandler) NewConsumerHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.NewAPIConsumerRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		h.response.BadRequestResponse(w, r, err)
		return
	}

	h.validator.Check(requestBody.Email != "", "email", "must be provided")
	h.validator.Check(requestBody.FirstName != "", "first_name", "must be provided")
	h.validator.Check(requestBody.LastName != "", "last_name", "must be provided")
	h.validator.Check(requestBody.Website != "", "website", "must be provided")
	h.validator.Check(validator.Matches(requestBody.Email, validator.EmailRX), "email", "must be a valid email address")

	if requestBody.ClubAffiliation != "" {
		msg := "not a valid club affiliation" // TODO: create a better message and figure out how to solve
		h.validator.AddError("club_affiliation", msg)
	}
	if !h.validator.Valid() {
		h.response.FailedValidationResponse(w, r, h.validator.ReturnErrors())
		return
	}

	// now that the validation is completed we need to create the domain model
	// and then pass it to the service layer for processing
	consumer := mappers.TransformNewAPIConsumerRequestToDomainModel(requestBody)
	resp, err := h.service.CreateNewConsumer(consumer)
	if err != nil {
		h.response.ServerErrorResponse(w, r, err)
		return
	}

	// TODO: this request needs to be sent to the service/domain layer
	// right now it is just a placeholder
	h.response.WriteJSON(w, http.StatusCreated, resp, nil)

}
