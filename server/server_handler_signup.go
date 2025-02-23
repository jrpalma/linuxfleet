package server

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/jrpalma/linuxfleet/data"
)

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Server) signup(c echo.Context) error {
	sc := h.ServerContext(c)

	var user signupRequest
	if err := c.Bind(&user); err != nil {
		return sc.BadRequest(c, "Invalid request payload")
	}

	signupLink, err := uuid.NewRandom()
	if err != nil {
		return sc.InternalError(c, "Failed to generate signup link")
	}

	salt, err := uuid.NewRandom()
	if err != nil {
		return sc.InternalError(c, "Failed to generate signup salt")
	}

	htmlContent, err := sc.HtmlTemplates().Execute("signup-email", signupLink.String())
	if err != nil {
		return sc.InternalError(c, "Failed to execute email template")
	}

	hash := sha256.New()
	hash.Write(salt[:])
	hash.Write([]byte(user.Password))
	passwordHash := hex.EncodeToString(hash.Sum(nil))

	signupObject := data.Object{
		ID:      signupLink.String(),
		OwnerID: user.Email,
		Version: 1,
		Attributes: map[string]any{
			"email":     user.Email,
			"password":  passwordHash,
			"salt":      salt.String(),
			"createdAt": time.Now(),
		},
	}

	err = sc.Tables().Insert("signup", signupObject)
	if err != nil {
		return sc.InternalError(c, "Failed to store user data")
	}

	from := mail.NewEmail("LinuxFleet Support", "support@linuxfleet.com")
	subject := "LinuxFleet Registration"
	to := mail.NewEmail(user.Email, user.Email)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	err = sc.SendEmail(message)
	if err != nil {
		return sc.InternalError(c, "Failed to send registration email")
	}

	return sc.OK(c, "User signup initiated successfully")
}
