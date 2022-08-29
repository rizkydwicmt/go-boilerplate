package middleware

import (
	"time"

	helper "go-boilerplate/utilities"
	"go-boilerplate/utilities/redis"

	"go-boilerplate/handler/auth/request"
	authService "go-boilerplate/service/auth"
	userModel "go-boilerplate/service/auth/model"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	dayToSecond = int((time.Hour * 24).Seconds())
)

type middleware struct {
	authService authService.Service
}

func NewAuth(authService authService.Service) *middleware {
	return &middleware{authService}
}

func (m *middleware) Auth() (*jwt.GinJWTMiddleware, error) {
	var identityKey = "name"
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*userModel.User); ok {
				return jwt.MapClaims{
					identityKey: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &userModel.User{
				Name: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginRequest request.Login
			if err := c.ShouldBind(&loginRequest); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			user, err := m.authService.Login(loginRequest.Username, loginRequest.Password)
			if err != nil {
				return "", jwt.ErrFailedAuthentication
			} else {
				c.Set("redis.user_JWT", user)
				return &userModel.User{
					Name: user.Name,
				}, nil
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*userModel.User); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"rc":   code,
				"rd":   message,
				"data": nil,
			})
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			data, _ := c.Get("redis.user_JWT")
			redis.Set("aci:jwt:"+message, data, dayToSecond)
			c.JSON(code, gin.H{
				"rc":   code,
				"rd":   "sukses ambil jwt token",
				"data": gin.H{"token": message},
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		helper.Log("error")("JWT Error", err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		helper.Log("error")("authMiddleware.MiddlewareInit() Error", errInit.Error())
	}

	return authMiddleware, err
}
