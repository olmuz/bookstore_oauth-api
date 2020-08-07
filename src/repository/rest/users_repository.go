package rest

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/olmuz/bookstore_oauth-api/src/domain/users"
	"github.com/olmuz/bookstore_oauth-api/src/utils/errors"
)

const (
	baseURL = "http://localhost:8080"
)

var (
	userRestClient *resty.Client
)

func init() {
	userRestClient = resty.New()
	userRestClient.SetTimeout(100 * time.Millisecond)
}

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct{}

var endpoints = map[string]string{
	"login": "/users/login",
}

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(username string, password string) (*users.User, *errors.RestErr) {
	url := baseURL + endpoints["login"]
	request := users.UserLoginRequest{
		Email:    username,
		Password: password,
	}
	var result users.User
	var restErr errors.RestErr
	resp, err := userRestClient.R().
		SetBody(&request).
		SetResult(&result).
		SetError(&restErr).
		Post(url)

	if err != nil {
		return nil, errors.NewInternalServerError("invalid response when trying to login user")
	}
	if resp.StatusCode() != 200 {
		return nil, &restErr
	}
	return &result, nil
}
