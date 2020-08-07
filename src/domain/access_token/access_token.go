package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/olmuz/bookstore_oauth-api/src/utils/crypto_utils"
	"github.com/olmuz/bookstore_oauth-api/src/utils/errors"
)

const (
	expirationTime             = 24 * time.Hour
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client credentials grant type
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return errors.NewBadRequestError("invalid grant_type parameter")
	}

	return nil
}

// AccessToken represents access token
type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	CliendID    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

// NewAccessToken creates token and shifts expiration time to 24 hours
func NewAccessToken(userID int64) *AccessToken {
	return &AccessToken{
		UserID:  userID,
		Expires: time.Now().UTC().Add(expirationTime).Unix(),
	}
}

// Validate validates AccessToken struct
func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token id")
	}
	if at.UserID <= 0 {
		return errors.NewBadRequestError("invalid user_id")
	}
	if at.CliendID <= 0 {
		return errors.NewBadRequestError("invalid client id")
	}
	if at.Expires <= 0 {
		return errors.NewBadRequestError("invalid expiration time")
	}
	return nil
}

// IsExpired checks if token expired
func (at *AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
