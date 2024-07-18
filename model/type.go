package model

import (
	"SSE/db"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Kind string `json:"kind" binding:"required"`
	From string `json:"from" binding:"required"`
	To   string `json:"to,omitempty" binding:"required"`
	Data string `json:"data,omitempty" binding:"required"`
}

func filterUsers(username string) (filterUsers []gin.H) {
	for k := range db.Users {
		if k != username {
			filterUsers = append(filterUsers, gin.H{
				"Username": k,
				"Online":   false,
			})
		}
	}
	return
}
