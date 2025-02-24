package server

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegistration(t *testing.T) {
	Convey("Scenario: The admin initiates the registration for the first time", t, func() {
		server := testServer()
		Convey("When POST /api/registration/initiate without a body", func() {
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/initiate", nil)
			server.initiateRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/initiate with invalid email", func() {
			initiateRequest := &initiateRegistrationRequest{Email: "invalid", Password: "abc"}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/initiate", initiateRequest)
			server.initiateRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/initiate with invalid password", func() {
			initiateRequest := &initiateRegistrationRequest{Email: "user@example.com", Password: "abc"}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/initiate", initiateRequest)
			server.initiateRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/initiate with valid request", func() {
			initiateRequest := &initiateRegistrationRequest{Email: "user@example.com", Password: "abc123#8"}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/initiate", initiateRequest)
			server.initiateRegistrationHandler(tc.EchoContext)

			So(tc.HttpResponse.Code, ShouldEqual, http.StatusOK)
			response := &initiateRegistrationResponse{}
			err := tc.UnmarshalResponse(response)
			So(err, ShouldBeNil)
			So(response.Token, ShouldNotEqual, "")
		})
	})

	Convey("Scenario: The admin tries to complete the registration", t, func() {
		server := testServer()
		Convey("When POST /api/registration/complete with invalid request", func() {
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/complete", nil)
			server.initiateRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/complete with a bad seed", func() {
			completeRequest := &completeRegistrationRequest{Token: uuid.NewString(), Seed: "1234"}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/complete", completeRequest)
			server.completeRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When POST /api/registration/complete without existing registration", func() {
			completeRequest := &completeRegistrationRequest{Token: uuid.NewString(), Seed: "12345678"}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/complete", completeRequest)
			server.completeRegistrationHandler(tc.EchoContext)
			So(tc.HttpResponse.Code, ShouldEqual, http.StatusNotFound)
		})
		Convey("Given POST /api/registration/initiate with valid request", func() {
			initiateRequest := &initiateRegistrationRequest{Email: "user@example.com", Password: "abc123#8"}
			tc := server.EchoTestContext(http.MethodPost, "/api/registration/initiate", initiateRequest)
			server.initiateRegistrationHandler(tc.EchoContext)

			So(tc.HttpResponse.Code, ShouldEqual, http.StatusOK)
			response := &initiateRegistrationResponse{}
			err := tc.UnmarshalResponse(response)
			So(err, ShouldBeNil)
			So(response.Token, ShouldNotEqual, "")

			Convey("Given POST /api/registration/complete with valid request", func() {
				completeRequest := &completeRegistrationRequest{Token: response.Token, Seed: "12345678"}
				tc := server.EchoTestContext(http.MethodPost, "/api/registration/complete?token="+response.Token, completeRequest)
				server.completeRegistrationHandler(tc.EchoContext)
				So(tc.HttpResponse.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}
