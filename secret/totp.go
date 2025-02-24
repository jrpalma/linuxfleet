package secret

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"

	"github.com/pquerna/otp/totp"
)

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

func generateSecretForTOTP(userUUID string, saltUUID string) string {
	authCode := hmac.New(sha1.New, []byte(saltUUID))
	seed := authCode.Sum([]byte(userUUID))
	return b32NoPadding.EncodeToString(seed)
}

func ValidateCodeForTOTP(userUUID string, saltUUID string, code string) bool {
	seed := generateSecretForTOTP(userUUID, saltUUID)
	result := totp.Validate(code, seed)
	return result
}
