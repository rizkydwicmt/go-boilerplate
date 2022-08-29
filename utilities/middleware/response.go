package middleware

import (
	"github.com/gin-gonic/gin"
)

func ResponseHandler(c *gin.Context) {
	data, _ := c.Get("response.data")
	c.JSON(c.GetInt("response.status"), gin.H{
		"rc":   c.GetInt("response.rc"),
		"rd":   c.GetString("response.rd"),
		"data": data,
	})
}
