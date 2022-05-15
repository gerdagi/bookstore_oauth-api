package rest

import (
	"encoding/json"
	"time"

	resterrors "github.com/gerdagi/bookstore_utils-go/rest_errors"

	"github.com/gerdagi/bookstore_oauth-api/src/domain/users"
	"github.com/mercadolibre/golang-restclient/rest"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "https://api.bookstore.com",
		Timeout: 100 * time.Microsecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *resterrors.RestError)
}

type usersRepository struct{}

func NewRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email, password string) (*users.User, *resterrors.RestError) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := usersRestClient.Post("/users/login", request)

	if response == nil || response.Response == nil {
		return nil, resterrors.NewInternalServerError("invalid restclient response when trying to login user", nil)
	}

	if response.StatusCode > 299 {
		var restErr resterrors.RestError
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, resterrors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, resterrors.NewInternalServerError("error when trying to unmarshal user response", err)
	}

	return &user, nil
}
