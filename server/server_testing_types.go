package server

import (
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type TestContext struct {
	EchoContext  echo.Context
	HttpResponse *httptest.ResponseRecorder
}

type EmailSenderMock struct {
	res *rest.Response
	err error
}

func (esm *EmailSenderMock) Send(email *mail.SGMailV3) (*rest.Response, error) {
	return esm.res, esm.err
}
