package main

import (
	"SSE/model"
	"SSE/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	sse := service.NewSSEvent()
	r.GET("/subscribe", service.Subscribe(sse))
	r.GET("/unsubscribe", service.Unsubscribe(sse))
	r.POST("/send", func(c *gin.Context) {
		var msg model.Message
		if err := c.BindJSON(&msg); err != nil {
			return
		}
		sse.Message <- msg
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})
	r.Run(":8000")

}
