package controllers

import (
	"github.com/gin-gonic/gin"
)

func Users(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "controllers OK",
	})
}
