package api_v1

import (
	"sort"

	"github.com/gin-gonic/gin"
	grpc_server "github.com/mrtdeh/centor/pkg/grpc/server"
)

func GetNodes(c *gin.Context) {
	// h := getServerAPI()

	res := h.GetClusterNodes()
	var r []any
	for _, v := range res {
		r = append(r, v)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i].(grpc_server.NodeInfo).Id < r[j].(grpc_server.NodeInfo).Id
	})
	c.JSON(200, gin.H{
		"result": r,
	})
}
