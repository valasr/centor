package api_v1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrtdeh/centor/pkg/kive"
)

type kvRequest struct {
	Namespace string `json:"ns"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func PutKV(c *gin.Context) {
	var kvr kvRequest
	if err := c.ShouldBind(&kvr); err != nil {
		c.JSON(400, gin.H{"error": "invalid input data"})
		return
	}
	ts := time.Now().Unix()
	dc := h.GetMyDC()
	kv, err := kive.Put(dc, kvr.Namespace, kvr.Key, kvr.Value, ts)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := kive.Sync(kv); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func GetKV(c *gin.Context) {
	ns := c.Param("ns")
	key := c.Param("key")

	res, err := kive.Get(ns, key)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": res})
}

func DeleteKV(c *gin.Context) {
	ns := c.Param("ns")
	key := c.Param("key")

	ts := time.Now().Unix()
	kv, err := kive.Del(ns, key, ts)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := kive.Sync(kv); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}
