package api_v1

import (
	"time"

	"github.com/gin-gonic/gin"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
	"github.com/mrtdeh/centor/pkg/kive"
)

type kvRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func PutKV(c *gin.Context) {
	var kvr kvRequest
	if err := c.ShouldBind(&kvr); err != nil {
		c.JSON(400, gin.H{"error": "invalid input data"})
		return
	}
	ts := time.Now().Unix()
	kv, err := kive.Put(kvr.Key, kvr.Value, ts)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	grpc_server.Sync(*kv, "")

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

	ts := time.Now().Unix()
	kv, err := kive.Del(key, ts)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	grpc_server.Sync(*kv, "")

	c.JSON(200, gin.H{"status": "ok"})
}
