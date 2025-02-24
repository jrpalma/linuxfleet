package server

import (
	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	TOTP     string
}

func (h *Server) loginHandler(c echo.Context) error {
	sc := h.ServerContext(c)

	var request loginRequest
	if err := sc.BindModel(&request); err != nil {
		return sc.BadRequest(err.Error())
	}

	return nil
}
