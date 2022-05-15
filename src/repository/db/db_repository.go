package db

import (
	"github.com/gerdagi/bookstore_oauth-api/src/clients/cassandra"
	accesstoken "github.com/gerdagi/bookstore_oauth-api/src/domain/access_token"
	resterrors "github.com/gerdagi/bookstore_utils-go/rest_errors"

	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id , expires FROM access_tokens WHERE access_token = ?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id , expires) VALUES(?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*accesstoken.AccessToken, *resterrors.RestError)
	Create(accesstoken.AccessToken) *resterrors.RestError
	UpdateExpirationTime(accesstoken.AccessToken) *resterrors.RestError
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) UpdateExpirationTime(at accesstoken.AccessToken) *resterrors.RestError {
	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return resterrors.NewInternalServerError(err.Error(), nil)
	}

	return nil
}

func (r *dbRepository) Create(at accesstoken.AccessToken) *resterrors.RestError {
	if err := cassandra.GetSession().Query(queryCreateAccessToken, at.AccessToken, at.UserId, at.ClientId, at.Expires).Exec(); err != nil {
		return resterrors.NewInternalServerError(err.Error(), nil)
	}

	return nil
}

func (r *dbRepository) GetById(id string) (*accesstoken.AccessToken, *resterrors.RestError) {
	var result accesstoken.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, resterrors.NewNotFoundError("no access token found with given id")
		}
		return nil, resterrors.NewInternalServerError(err.Error(), nil)
	}

	return &result, nil
}
