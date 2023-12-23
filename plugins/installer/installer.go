package installer_plugin

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	PluginKits "github.com/mrtdeh/centor/plugins/assets"
)

type PluginProvider struct {
	PluginKits.PluginProps
}

func (p *PluginProvider) SetHandler(h PluginKits.CoreHandlerInterface) {
	p.Handler = h
}

func (p *PluginProvider) SetRouter(r http.Handler) {
	p.Router = r
}

func (p *PluginProvider) Init() error {
	r, ok := p.Router.(*gin.Engine)
	if !ok {
		return fmt.Errorf("router is not a gin router")
	}
	r.POST("/send-file", sendFile)

	p.Router = r
	return nil

}

var h PluginKits.CoreHandlerInterface

// Run method for ExamplePlugin1
func (p *PluginProvider) Run() {
	fmt.Printf("Plugin %s is running...\n", p.Name)
	h = p.Handler

	err := h.WaitForReady(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

}

type SendFileRequest struct {
	Filename    string                `json:"filename"`
	Data        string                `json:"data"`
	NodeId      string                `json:"node_id"`
	PackageFile *multipart.FileHeader `form:"deb" json:"-"`
}

func sendFile(c *gin.Context) {
	var req SendFileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var filename string
	contentType := c.GetHeader("Content-Type")
	if contentType != "application/json" { // if data is json

		// save uploaded file
		filename += req.PackageFile.Filename
		err := c.SaveUploadedFile(req.PackageFile, filename)
		if err != nil {
			log.Println("error in save file : ", err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		data, _ := os.ReadFile(filename)
		req.Data = string(data)
		req.Filename = path.Base(filename)
		req.NodeId = c.PostForm("node_id")
	}

	err := h.SendFile(context.Background(), req.NodeId, req.Filename, []byte(req.Data))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	cmd := "dpkg -i /tmp/centor-recieved/" + req.Filename
	installRes, err := h.Exec(context.Background(), req.NodeId, cmd)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	cmd = "rm /tmp/centor-recieved/" + req.Filename
	_, err = h.Exec(context.Background(), req.NodeId, cmd)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Exec request on %s : \n%s\n", req.NodeId, installRes)
	c.JSON(200, "ok")
}
