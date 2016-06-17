package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Root is the handler for the GET / route.
// This will respond by rendering the home html page.
func Root(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{})
}
