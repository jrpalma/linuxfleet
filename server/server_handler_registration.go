package server

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/jrpalma/linuxfleet/data"
)

type initiateRegistrationRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}
type initiateRegistrationResponse struct {
	Token string
}

func (h *Server) initiateRegistrationHandler(c echo.Context) error {
	sc := h.ServerContext(c)

	var request initiateRegistrationRequest
	if err := sc.BindModel(&request); err != nil {
		return sc.BadRequest(err.Error())
	}

	token, err := uuid.NewRandom()
	if err != nil {
		return sc.InternalError("Failed to generate signup link")
	}

	registrationURL := sc.FormatURL("/registration/complete?token=%v", token)
	templateValues := map[string]any{"URL": registrationURL}

	salt, err := uuid.NewRandom()
	if err != nil {
		return sc.InternalError("Failed to generate signup salt")
	}

	htmlEmailContent, err := sc.ExecuteTemplate("registration-email.tmpl", templateValues)
	if err != nil {
		return sc.InternalError("Failed to execute email template")
	}

	hash := sha256.New()
	hash.Write([]byte(salt.String()))
	hash.Write([]byte(request.Password))
	passwordHash := hex.EncodeToString(hash.Sum(nil))

	registrationObject := data.Object{
		ID:      token.String(),
		Version: 1,
		Attributes: map[string]any{
			"email":    request.Email,
			"password": passwordHash,
			"salt":     salt.String(),
		},
	}

	err = sc.DataInsert("registration", registrationObject)
	if err != nil {
		return sc.InternalError("Failed to store user data")
	}

	to := mail.NewEmail(request.Email, request.Email)
	from := mail.NewEmail("LinuxFleet Support", "support@linuxfleet.com")
	message := mail.NewSingleEmail(from, "LinuxFleet Registration", to, "", htmlEmailContent)

	err = sc.SendEmail(message)
	if err != nil {
		return sc.InternalError("Failed to send registration email")
	}

	return sc.OKJSON(map[string]any{"token": token.String()})
}

type completeRegistrationRequest struct {
	Token string `validate:"required,uuid"`
	Seed  string `validate:"required,min=8"`
}

func (h *Server) completeRegistrationHandler(c echo.Context) error {
	sc := h.ServerContext(c)

	var request completeRegistrationRequest
	if err := sc.BindModel(&request); err != nil {
		return sc.BadRequest(err.Error())
	}

	registrationObject, err := sc.DataGetByID("registration", request.Token)
	if errors.Is(err, sql.ErrNoRows) {
		return sc.NotFound("The registration does not exists")
	} else if err != nil {
		return sc.InternalError(err.Error())
	}

	adminID, err := uuid.NewRandom()
	if err != nil {
		return sc.InternalError("Failed not generate admin ID")
	}

	adminObject := data.Object{Attributes: registrationObject.Attributes, ID: adminID.String(), Version: 1}
	adminObject.Attributes["seed"] = request.Seed

	if err := sc.DataInsert("admin", adminObject); err != nil {
		return sc.InternalError("Failed to save administrator")
	}

	return sc.OK("User registration was completed successfully")
}
