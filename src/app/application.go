package app

import (
	"github.com/gin-gonic/gin"
	"github.com/olmuz/bookstore_oauth-api/src/http"
	"github.com/olmuz/bookstore_oauth-api/src/repository/db"
	"github.com/olmuz/bookstore_oauth-api/src/repository/rest"
	"github.com/olmuz/bookstore_oauth-api/src/services/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {
	dbRepository := db.NewRepository()
	httpRepository := rest.NewRepository()
	atService := access_token.NewService(httpRepository, dbRepository)
	atHandler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetByID)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8081")
}
