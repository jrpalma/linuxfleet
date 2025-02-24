package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/jrpalma/linuxfleet/data"
	"github.com/jrpalma/linuxfleet/html"
)

type Server struct {
	tables    *data.Tables
	email     EmailSender
	echo      *echo.Echo
	templates *html.Templates
	validator *validator.Validate
}

func NewServer(tables *data.Tables, templates *html.Templates, email EmailSender) *Server {
	server := &Server{
		validator: validator.New(validator.WithRequiredStructEnabled()),
		templates: templates,
		echo:      echo.New(),
		tables:    tables,
		email:     email,
	}
	server.echo.HideBanner = true
	server.echo.POST("/api/registration/initiate", server.initiateRegistrationHandler)
	server.echo.POST("/api/registration/complete", server.completeRegistrationHandler)
	return server
}

func (s *Server) ServerContext(c echo.Context) *ServerContext {
	return &ServerContext{
		validator: s.validator,
		templates: s.templates,
		tables:    s.tables,
		email:     s.email,
		ec:        c,
	}
}

func (s *Server) EchoTestContext(method string, target string, body any) *TestContext {
	data, err := json.Marshal(body)
	if err != nil {
		log.Fatal("could not marshal echo test context body")
	}
	reader := bytes.NewReader(data)
	request := httptest.NewRequest(method, target, reader)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	return &TestContext{
		EchoContext:  s.echo.NewContext(request, recorder),
		HttpResponse: recorder,
	}
}
