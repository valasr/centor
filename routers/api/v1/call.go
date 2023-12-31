package api_v1

import (
	"context"

	"github.com/gin-gonic/gin"
)

func Call(c *gin.Context) {
	// h := getServerAPI()
	tags, err := h.Call(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"result": tags,
	})
}
