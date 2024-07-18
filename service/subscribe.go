package service

import (
	"SSE/db"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Subscribe(sse *SSEvent) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Query("user")
		if _, ok := db.Users[user]; !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("user not found"))
		}

		client := Client{
			name: user,
			ch:   make(ClientChan),
		}

		sse.NewClient <- client

		defer func() {
			sse.CloseClient <- client
		}()

		c.Stream(func(w io.Writer) bool {
			if message, ok := <-client.ch; ok {
				c.SSEvent("message", message)
				return true
			}
			return false
		})

	}
}

func Unsubscribe(sse *SSEvent) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Query("user")
		if _, ok := db.Users[user]; !ok {
			c.AbortWithError(http.StatusBadRequest, errors.New("user not found"))
		}

		client := Client{
			name: user,
			ch:   make(ClientChan),
		}

		sse.CloseClient <- client

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
	}
}
