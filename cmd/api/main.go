package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(testMiddelware)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:8081") // listen and serve on 0.0.0.0:8080
}

func testMiddelware(c *gin.Context) {
	fmt.Println("testing!")
}
