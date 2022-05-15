package accesstoken

import (
	"fmt"
	"strings"
	"time"

	cryptoutils "github.com/gerdagi/bookstore_oauth-api/src/utils/crypto_utils"
	resterrors "github.com/gerdagi/bookstore_utils-go/rest_errors"
)

const (
	expirationTime             = 24
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
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

// Web Frontend Client Id = 123
// Android APP Client Id = 678
// IOS APP Client Id = 443

func (at *AccessTokenRequest) Validate() *resterrors.RestError {
	if at.GrantType != grantTypePassword ||
		at.GrantType != grantTypeClientCredentials {
		return resterrors.NewBadRequestError("invalid grant type")
	}

	//TODO: Validate parameters for each grant_type
	switch at.GrantType {
	case grantTypePassword:
		if strings.TrimSpace(at.Username) == "" ||
			strings.TrimSpace(at.Password) == "" {
			return resterrors.NewBadRequestError("invalid username or password")
		}
		break
	case grantTypeClientCredentials:
		if strings.TrimSpace(at.ClientId) == "" ||
			strings.TrimSpace(at.ClientSecret) == "" {
			return resterrors.NewBadRequestError("invalid client id  or client secret")
		}
		break
	default:
		return resterrors.NewBadRequestError("invalid grant type parameter")
	}

	return nil
}

func (at *AccessToken) Validate() *resterrors.RestError {
	if len(strings.TrimSpace(at.AccessToken)) == 0 {
		return resterrors.NewBadRequestError("invalid access token id")
	}

	if at.UserId <= 0 {
		return resterrors.NewBadRequestError("invalid  User Id")
	}

	if at.ClientId <= 0 {
		return resterrors.NewBadRequestError("invalid  Client Id")
	}

	if at.Expires <= 0 {
		return resterrors.NewBadRequestError("invalid  Expiration Time")
	}

	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = cryptoutils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
