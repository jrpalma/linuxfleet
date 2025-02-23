package server

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/jrpalma/linuxfleet/data"
	"github.com/jrpalma/linuxfleet/html"
)

type EmailSender interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

type ServerContext struct {
	ec        echo.Context
	tables    *data.Tables
	email     EmailSender
	templates *html.Templates
}

func (sc *ServerContext) BadRequest(message string) error {
	return sc.errorResponse(http.StatusBadRequest, message)
}

func (h *ServerContext) InternalError(message string) error {
	return h.errorResponse(http.StatusInternalServerError, message)
}

func (h *ServerContext) OK(message string) error {
	return h.ec.JSON(http.StatusOK, map[string]string{"message": message})
}

func (h *ServerContext) SendEmail(email *mail.SGMailV3) error {
	_, err := h.email.Send(email)
	return err
}

func (c *ServerContext) Tables() *data.Tables {
	return c.tables
}

func (c *ServerContext) HtmlTemplates() *html.Templates {
	return c.templates
}

func (h *ServerContext) GetEnv(name string) string {
	return os.Getenv(name)
}

func (h *ServerContext) errorResponse(status int, errMessage string) error {
	return h.ec.JSON(status, map[string]string{"error": errMessage})
}
