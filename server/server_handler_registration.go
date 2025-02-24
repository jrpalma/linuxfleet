package server

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/jrpalma/linuxfleet/data"
)

type startRegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *Server) startRegistrationHandler(c echo.Context) error {
	sc := h.ServerContext(c)
	validate := validator.New(validator.WithRequiredStructEnabled())

	var user startRegistrationRequest
	if err := c.Bind(&user); err != nil {
		return sc.BadRequest("Invalid request payload")
	}

	if err := validate.Struct(&user); err != nil {
		return sc.BadRequest(err.Error())
	}

	registrationToken, err := uuid.NewRandom()
	if err != nil {
		return sc.InternalError("Failed to generate signup link")
	}

	registrationURL := filepath.Join(sc.GetEnv("BASE_URL"),
		"/registration/complete?token="+registrationToken.String())
	templateValues := map[string]any{"URL": registrationURL}

	salt, err := uuid.NewRandom()
	if err != nil {
		return sc.InternalError("Failed to generate signup salt")
	}

	htmlEmailContent, err := sc.HtmlTemplates().Execute("registration-email.tmpl", templateValues)
	if err != nil {
		return sc.InternalError("Failed to execute email template")
	}

	hash := sha256.New()
	hash.Write(salt[:])
	hash.Write([]byte(user.Password))
	passwordHash := hex.EncodeToString(hash.Sum(nil))

	signupObject := data.Object{
		ID:      registrationToken.String(),
		OwnerID: user.Email,
		Version: 1,
		Attributes: map[string]any{
			"email":    user.Email,
			"password": passwordHash,
			"salt":     salt.String(),
		},
	}

	err = sc.Tables().Insert("signup", signupObject)
	if err != nil {
		return sc.InternalError("Failed to store user data")
	}

	to := mail.NewEmail(user.Email, user.Email)
	from := mail.NewEmail("LinuxFleet Support", "support@linuxfleet.com")
	message := mail.NewSingleEmail(from, "LinuxFleet Registration", to, "", htmlEmailContent)

	err = sc.SendEmail(message)
	if err != nil {
		return sc.InternalError("Failed to send registration email")
	}

	return sc.OK("User signup initiated successfully")
}
