package access_token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {
	at := NewAccessToken()
	assert.NotNil(t, at, "access token should not be nil")

	assert.EqualValues(t, at.AccessToken, "", "access token has to be an empty string")
	assert.EqualValues(t, at.UserID, 0, "UserID has to be equual to 0")
	assert.EqualValues(t, at.CliendID, 0, "CliendID has to be equual to 0")
	assert.False(t, at.IsExpired(), "brand new access token should not be expired")
}

func TestIsExpired(t *testing.T) {
	at := AccessToken{}
	assert.True(t, at.IsExpired(), "empty access token should be expired by default")

	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "access token created three hours from now should NOT be expired")
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, expirationTime, 24*time.Hour, "expirationTime has to be equal to 24 hours")
}
