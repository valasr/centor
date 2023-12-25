package api_v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mrtdeh/centor/pkg/kive"
)

type kvRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func PutKV(c *gin.Context) {
	var kv kvRequest
	if err := c.ShouldBind(&kv); err != nil {
		c.JSON(400, gin.H{"error": "invalid input data"})
		return
	}

	err := kive.Put(kv.Key, kv.Value)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func GetKV(c *gin.Context) {
	key := c.Param("key")

	res, err := kive.Get(key)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "result": res})
}

func DeleteKV(c *gin.Context) {
	key := c.Param("key")

	err := kive.Del(key)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}
