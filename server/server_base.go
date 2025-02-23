package server

import (
	"io"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/jrpalma/linuxfleet/data"
	"github.com/jrpalma/linuxfleet/html"
)

type Server struct {
	tables    *data.Tables
	email     EmailSender
	echo      *echo.Echo
	templates *html.Templates
}

func NewServer(tables *data.Tables, templates *html.Templates, email EmailSender) *Server {
	server := &Server{
		templates: templates,
		echo:      echo.New(),
		tables:    tables,
		email:     email,
	}
	server.echo.HideBanner = true
	server.echo.POST("/api/signup", server.signup)
	return server
}

func (s *Server) ServerContext(c echo.Context) *ServerContext {
	return &ServerContext{
		templates: s.templates,
		tables:    s.tables,
		email:     s.email,
		ec:        c,
	}
}

func (s *Server) EchoTestContext(method string, target string, reader io.Reader) *TestContext {
	request := httptest.NewRequest(method, target, reader)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	return &TestContext{
		EchoContext:  s.echo.NewContext(request, recorder),
		HttpResponse: recorder,
	}
}
