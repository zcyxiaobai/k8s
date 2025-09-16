package k8s

import (
	"bufio"
	"context"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
)

// 获取容器日志
// 获取容器日志并通过 WebSocket 推送
func (p *PodClient) GetPodLog(ctx context.Context, namespace, podname, cname string, conn *websocket.Conn) error {
	podLogOpts := corev1.PodLogOptions{
		Container: cname,
		Follow:    true,
		TailLines: int64Ptr(100),
	}

	req := p.client.CoreV1().Pods(namespace).GetLogs(podname, &podLogOpts)

	stream, err := req.Stream(ctx)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("获取日志失败: "+err.Error()))
		return err
	}
	defer stream.Close()
	defer conn.Close()

	scanner := bufio.NewScanner(stream)
	buf := make([]byte, 0, 1024*64)
	scanner.Buffer(buf, 1024*1024)

	pingTicker := time.NewTicker(20 * time.Second)
	defer pingTicker.Stop()
	writeDeadline := 10 * time.Second

	for {
		select {
		case <-ctx.Done():
			conn.WriteMessage(websocket.TextMessage, []byte("日志推送已断开"))
			return nil // 不返回错误
		case <-pingTicker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				return err
			}
		default:
			if scanner.Scan() {
				line := scanner.Text()
				_ = conn.SetWriteDeadline(time.Now().Add(writeDeadline))
				if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
					return err
				}
			} else {
				if err := scanner.Err(); err != nil && !errors.Is(err, context.Canceled) {
					conn.WriteMessage(websocket.TextMessage, []byte("日志读取错误: "+err.Error()))
				}
				return nil // 流结束或者 context canceled
			}
		}
	}
}
func int64Ptr(i int64) *int64 { return &i }
