package routers

import (
	"io"

	"github.com/gin-gonic/gin"
	api_v1 "github.com/mrtdeh/centor/routers/api/v1"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func pingHandler(c *gin.Context) {
	c.JSON(200, `pong`)
}

func InitRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/call", api_v1.Call)
	r.GET("/nodes", api_v1.GetNodes)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
