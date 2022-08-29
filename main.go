package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	handlerAuth "go-boilerplate/handler/auth"
	repoAuth "go-boilerplate/repository/auth"
	serviceAuth "go-boilerplate/service/auth"
	helper "go-boilerplate/utilities"

	"go-boilerplate/utilities/database"
	"go-boilerplate/utilities/middleware"
	"go-boilerplate/utilities/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("configs/config_local.json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// only debug devel
	if strings.ToLower(viper.GetString("ENVIRONMENT")) == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	/* use utilities */
	db := database.ConnDB()
	err = redis.Setup()
	router := gin.Default()
	server := &http.Server{
		Addr:    viper.GetString("APPPORT"),
		Handler: router,
	}

	if err != nil {
		panic(err)
	}

	/* use repo */
	authRepository := repoAuth.NewRepository(db)

	/* use service */
	authService := serviceAuth.NewService(authRepository)

	/* use global middleware */
	router.Use(helper.LoggerToFile())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.RequestHandler)
	authMiddleware := middleware.NewAuth(authService)
	auth, _ := authMiddleware.Auth()

	/* use handler */
	authHandler := handlerAuth.NewHandler(authService, auth)

	/* Group Handler */
	{
		authHandler.Router(router.Group("/auth"))
	}

	serverErr := server.ListenAndServe()
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
