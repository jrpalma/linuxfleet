package secret

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pquerna/otp/totp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTOTP(t *testing.T) {
	Convey("Scenario: The user attempts to validate TOTP code", t, func() {
		now := time.Now()
		userUUID := uuid.NewString()
		saltUUID := uuid.NewString()
		secret := generateSecretForTOTP(userUUID, saltUUID)

		code, err := totp.GenerateCode(secret, now)
		So(err, ShouldBeNil)

		valid := totp.Validate(code, secret)
		So(valid, ShouldBeTrue)
	})
}
