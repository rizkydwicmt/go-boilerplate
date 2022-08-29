package handler

import (
	"go-boilerplate/handler/auth/request"
	authService "go-boilerplate/service/auth"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService    authService.Service
	authMiddleware *jwt.GinJWTMiddleware
}

func NewHandler(authService authService.Service, authMiddleware *jwt.GinJWTMiddleware) *authHandler {
	return &authHandler{authService, authMiddleware}
}

func (h *authHandler) Router(router *gin.RouterGroup) {

	router.POST("/login", h.authMiddleware.LoginHandler)

	router.POST("/register", func(c *gin.Context) {
		var registerRequest request.Register
		err := c.ShouldBindJSON(&registerRequest)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"rc":   500,
				"rd":   "failed bind json request",
				"data": err.Error(),
			})
			return
		}
		user, _ := h.authService.Register(registerRequest.ToModel())
		c.JSON(http.StatusOK, gin.H{
			"rc":   200,
			"rd":   "sukses root",
			"data": user,
		})
	})
}
