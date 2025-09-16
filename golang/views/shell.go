package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

func (p *PodView) GetShell(conf *rest.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			klog.Error(err.Error())
			c.JSON(http.StatusNotAcceptable, gin.H{
				"message": "获取shell失败",
			})
			return
		}
		namespace := c.Param("namespace")
		podname := c.Param("podname")
		cname := c.Param("cname")
		err = p.PodClient.GetShell(namespace, podname, cname, conn, conf)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
	}
}
