package rest

import (
	"net/http"
	"os"
	"testing"

	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	httpmock.ActivateNonDefault(userRestClient.GetClient())
	defer httpmock.DeactivateAndReset()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Reset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(60 * time.Millisecond)
			return httpmock.NewJsonResponse(200, map[string]interface{}{})
		},
	)

	repo := NewRepository()

	result, err := repo.LoginUser("Bill", "Kill")

	assert.Nil(t, result)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Reset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, map[string]interface{}{
				"message": "invalid login credentials",
				"status":  "400",
				"error":   "not_found",
			})
		},
	)

	repo := NewRepository()

	result, err := repo.LoginUser("Bill", "Kill")

	assert.Nil(t, result)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid response when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Reset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, map[string]interface{}{
				"message": "invalid login credentials",
				"status":  400,
				"error":   "not_found",
			})
		},
	)

	repo := NewRepository()

	result, err := repo.LoginUser("Bill", "Kill")

	assert.Nil(t, result)
	assert.EqualValues(t, http.StatusBadRequest, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidLoginJsonResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Reset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, `{name}`)
		},
	)

	repo := NewRepository()

	result, err := repo.LoginUser("Bill", "Kill")

	assert.Nil(t, result)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid response when trying to login user", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Reset()

	httpmock.RegisterResponder("POST", "http://localhost:8080/users/login",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"id":         1,
				"first_name": "Bill",
				"last_name":  "Kill",
				"email":      "kill.bill@gmail.com",
			})
		},
	)

	repo := NewRepository()

	result, err := repo.LoginUser("Bill", "Kill")
	assert.Nil(t, err)
	assert.EqualValues(t, 1, result.ID)
	assert.EqualValues(t, "Bill", result.FirstName)
	assert.EqualValues(t, "Kill", result.LastName)
	assert.EqualValues(t, "kill.bill@gmail.com", result.Email)
}
