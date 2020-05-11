package main

import (
	"github.com/blake86/gin-poppy"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
)

func test2(ctx *gin.Context) {
	log.Println("in TEST2:")
}

func test3(ctx *gin.Context) {
	log.Println("in TEST3:")
}

func main() {
	r := gin.New()

	gin_poppy.Register(r)

	r.GET("/test", func(c *gin.Context) {
		r := rand.Int() % 10
		switch r {
		case 1:
			c.JSON(http.StatusOK, nil)
		case 2:
			c.JSON(http.StatusNoContent, nil)
		case 3:
			c.JSON(http.StatusBadRequest, nil)
		case 4:
			c.JSON(http.StatusInternalServerError, nil)
		default:
			c.JSON(http.StatusOK, nil)
		}
		log.Println("in /test1:")
	})

	r.GET("/test2", test2)

	r.GET("/test3/:id", test3)

	// Listen and serve on 0.0.0.0:8080
	r.Run("localhost:8999")
}
