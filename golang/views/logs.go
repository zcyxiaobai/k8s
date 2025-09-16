package views

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"
)

// 初始化一个 WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //允许跨域
	},
}

// 获取容器日志
func (p *PodView) GetLog(c *gin.Context) {
	// 将 HTTP 升级为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		klog.Error(err)
		return
	}

	namespace := c.Param("namespace")
	podname := c.Param("podname")
	cname := c.Param("cname")

	// 创建可取消的 context
	ctx, cancel := context.WithCancel(context.Background())

	// 当连接关闭时自动取消 context
	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				// 客户端关闭 WebSocket
				cancel()
				return
			}
		}
	}()

	// 后台推送日志
	err = p.PodClient.GetPodLog(ctx, namespace, podname, cname, conn)
	if err != nil && !errors.Is(err, context.Canceled) {
		conn.WriteMessage(websocket.TextMessage, []byte("日志推送出错: "+err.Error()))
	}
}
