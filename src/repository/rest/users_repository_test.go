package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start test cases..")
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore/com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"gerdagi@gmail.com", "password":"the-password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("gerdagi@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore/com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"gerdagi@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials", "status":"404", "error":"not_found"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("gerdagi@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore/com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"gerdagi@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials", "status":404, "error":"not_found"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("gerdagi@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "https://api.bookstore/com/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"gerdagi@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{
			"id": "1",
			"first_name": "Gültekin",
			"last_name": "Erdağı",
			"email": "gerdagi@gmail.com"
		}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("gerdagi@gmail.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)

}
func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore/com/users/login",
		ReqBody:      `{"email":"gerdagi@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{
			"id": 1,
			"first_name": "Gültekin",
			"email": "gerdagi@gmail.com"
		}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("gerdagi@gmail.com", "the-password")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, http.StatusOK, err.Status)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "gerdagi@gmail.com", user.Email)
}
