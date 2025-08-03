// package system handlers within this package are for the system, this means the apiconsumers
// the ones that can consume the api
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/maple/internal/adapters/rest/handlers/dto"
	"github.com/ctfrancia/maple/internal/adapters/rest/response"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type SystemHealthHandler struct {
	system   ports.SystemServicer
	response ports.SystemResponder
	logger   ports.Logger
}

func NewSystemHandler(ss ports.SystemServicer, log ports.Logger) *SystemHealthHandler {
	handler := &SystemHealthHandler{
		system:   ss,
		response: response.NewHelper(log),
		logger:   log,
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
		// badRequestResponse(w, r, err)
		return
	}
	/*
		// mashal the request body into a struct
		var requestBody model.NewAPIConsumerRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// Validate the request body that the required fields are present
		v := validator.New()
		v.Check(requestBody.Email != "", "email", "must be provided")
		v.Check(requestBody.FirstName != "", "first_name", "must be provided")
		v.Check(requestBody.LastName != "", "last_name", "must be provided")
		v.Check(requestBody.Website != "", "website", "must be provided")
		v.Check(validator.Matches(requestBody.Email, validator.EmailRX), "email", "must be a valid email address")
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		uuid := uuid.New().String()
		// Create a new auth model
		authModelUser := &repository.Auth{
			UUID:      uuid,
			Email:     requestBody.Email,
			FirstName: requestBody.FirstName,
			LastName:  requestBody.LastName,
			Website:   requestBody.Website,
		}

		// check if user is in DB and password if the user is in the DB then return a 409
		err := app.repository.Auth.SelectByEmail(authModelUser)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			app.conflictResponse(w, r)
			return
		}

		generatedPW, err := auth.CreateSecretKey(auth.PasswordGeneratorDefaultLength)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Hash the password
		encodedHash, err := auth.Hash(generatedPW)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Assign the argon2 hash to the user password
		authModelUser.Password = encodedHash

		// Create the user in DB
		err = app.repository.Auth.Create(authModelUser)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		/// Return the user with the generated password
		authModelUser.Password = generatedPW
		err = app.writeJSON(w, http.StatusCreated, envelope{"consumer": authModelUser}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	*/
}
