package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/jrpalma/linuxfleet/data"
	"github.com/jrpalma/linuxfleet/html"
)

type TestContext struct {
	EchoContext  echo.Context
	HttpResponse *httptest.ResponseRecorder
}

func (tc *TestContext) UnmarshalResponse(model any) error {
	err := json.Unmarshal(tc.HttpResponse.Body.Bytes(), model)
	return err
}

type EmailSenderMock struct {
	res *rest.Response
	err error
}

func (esm *EmailSenderMock) Send(email *mail.SGMailV3) (*rest.Response, error) {
	return esm.res, esm.err
}

func testServer() *Server {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err.Error())
	}

	tables, err := data.NewTables(db)
	if err != nil {
		log.Fatal(err.Error())
	}

	templates := html.NewTemplates()
	if err != nil {
		log.Fatal(err.Error())
	}

	emailSernder := &EmailSenderMock{}
	server := NewServer(tables, templates, emailSernder)
	return server
}
