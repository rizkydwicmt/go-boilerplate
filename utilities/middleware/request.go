package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestHandler(c *gin.Context) {
	c.Set("response.status", http.StatusInternalServerError)
	c.Set("response.rc", http.StatusInternalServerError)
	c.Set("response.rd", "Internal Server Error")
	c.Set("response.data", gin.H{})
}
