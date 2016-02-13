package main

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

type Login struct {
    User     string `form:"user" json:"user" binding:"required"`
    Password string `form:"password" json:"password" binding:"required"`
}

func process() {
    for {
        fmt.Println("processing")
        time.Sleep(time.Millisecond * 1000)
    }
}

func ping(c *gin.Context) {
    go process()

    c.JSON(200, gin.H{
        "message": "pong",
    })

}

func loginJSON(c *gin.Context) {
    var json Login
    if c.BindJSON(&json) == nil {
        fmt.Println(json.User)
        fmt.Println(json.Password)
    }
    fmt.Println(c)
}

func main() {

    // Create a default gin stack
    r := gin.Default()

    // Routes
    r.GET("/ping", ping)
    r.POST("/loginJSON", loginJSON)

    // Unleash!
    r.Run() // listen and serve on port defined by environment variable PORT
}
