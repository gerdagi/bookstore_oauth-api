package app

import (
	http "github.com/gerdagi/bookstore_oauth-api/src/http"
	"github.com/gerdagi/bookstore_oauth-api/src/repository/db"
	"github.com/gerdagi/bookstore_oauth-api/src/repository/rest"
	accesstoken "github.com/gerdagi/bookstore_oauth-api/src/services/access_token"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewHandler(
		accesstoken.NewService(rest.NewRepository(), db.NewRepository()))
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	//router.PUT("/oauth/access_token", atHandler.Update)
	router.Run(":8080")
}
