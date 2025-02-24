package server

import (
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegistration(t *testing.T) {
	Convey("Given the server", t, func() {
		server := testServer()
		Convey("When POST /api/registration/start without a body", func() {
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/start", nil)
			server.startRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/start with invalid email", func() {
			body := &startRegistrationRequest{
				Email:    "invalid",
				Password: "abc",
			}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/start", body)
			server.startRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/start with invalid password", func() {
			body := &startRegistrationRequest{
				Email:    "user@example.com",
				Password: "abc",
			}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/start", body)
			server.startRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/start with valid request", func() {
			body := &startRegistrationRequest{
				Email:    "user@example.com",
				Password: "abc123#8",
			}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/start", body)
			server.startRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusOK)
		})
	})
}
