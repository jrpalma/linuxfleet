package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
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
	validator *validator.Validate
}

func (sc *ServerContext) BadRequest(message string) error {
	return sc.errorResponse(http.StatusBadRequest, message)
}

func (sc *ServerContext) NotFound(message string) error {
	return sc.errorResponse(http.StatusNotFound, message)
}

func (sc *ServerContext) InternalError(message string) error {
	return sc.errorResponse(http.StatusInternalServerError, message)
}

func (sc *ServerContext) OK(message string) error {
	return sc.ec.JSON(http.StatusOK, map[string]string{"message": message})
}

func (sc *ServerContext) OKJSON(body any) error {
	return sc.ec.JSON(http.StatusOK, body)
}

func (sc *ServerContext) SendEmail(email *mail.SGMailV3) error {
	_, err := sc.email.Send(email)
	return err
}

func (sc *ServerContext) ExecuteTemplate(name string, data any) (string, error) {
	return sc.templates.Execute(name, data)
}

func (sc *ServerContext) GetEnv(name string) string {
	return os.Getenv(name)
}

func (sc *ServerContext) errorResponse(status int, errMessage string) error {
	return sc.ec.JSON(status, map[string]string{"error": errMessage})
}

func (sc *ServerContext) DataListByOwner(tableName string, ownerID string) ([]data.Object, error) {
	return sc.tables.ListByOwner(tableName, ownerID)
}

func (sc *ServerContext) DataInsert(tableName string, obj data.Object) error {
	return sc.tables.Insert(tableName, obj)
}

func (sc *ServerContext) DataDeleteByID(tableName string, id string) error {
	return sc.tables.DeleteByID(tableName, id)
}

func (sc *ServerContext) DataGetByID(tableName string, id string) (data.Object, error) {
	return sc.tables.GetByID(tableName, id)
}

func (sc *ServerContext) DataUpdateByID(tableName string, id string, obj data.Object) error {
	return sc.tables.UpdateByID(tableName, id, obj)
}

func (sc *ServerContext) BindModel(model any) error {
	if err := sc.ec.Bind(model); err != nil {
		return sc.BadRequest("Invalid request payload")
	}

	if err := sc.validator.Struct(model); err != nil {
		return sc.BadRequest(err.Error())
	}
	return nil
}

func (sc *ServerContext) FormatURL(pathFormat string, args ...any) string {
	fullURL := filepath.Join(sc.GetEnv("BASE_URL"), fmt.Sprintf(pathFormat, args...))
	return fullURL
}
