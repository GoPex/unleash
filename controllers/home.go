package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHome is the handler for the GET / route.
// This will respond by rendering the home html page.
func GetHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{})
}
