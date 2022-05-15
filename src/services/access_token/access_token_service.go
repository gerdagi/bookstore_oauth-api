package access_token

import (
	"strings"

	accesstoken "github.com/gerdagi/bookstore_oauth-api/src/domain/access_token"
	"github.com/gerdagi/bookstore_oauth-api/src/repository/db"
	"github.com/gerdagi/bookstore_oauth-api/src/repository/rest"
	resterrors "github.com/gerdagi/bookstore_utils-go/rest_errors"
)

type Service interface {
	GetById(string) (*accesstoken.AccessToken, *resterrors.RestError)
	Create(accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, *resterrors.RestError)
	UpdateExpirationTime(accesstoken.AccessToken) *resterrors.RestError
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*accesstoken.AccessToken, *resterrors.RestError) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, resterrors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, resterrors.NewBadRequestError("something wrong trying to get by id")
	}
	return accessToken, nil
}

func (s *service) Create(request accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, *resterrors.RestError) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	//TODO: Support both grant types: client_credentials and password

	// Authenticate the user against the Users API:
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token:
	at := accesstoken.GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at accesstoken.AccessToken) *resterrors.RestError {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}
