package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
)

// 交互的结构体，接管输入和输出
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// 消息内容
type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// 读数据的方法
// 返回值int是读成功了多少数据
func (t *TerminalSession) Read(p []byte) (int, error) {
	//从ws中读取消息
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		klog.Errorf("读取消息错误: %v", err)
		return 0, err
	}
	//反序列化
	var msg terminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		klog.Errorf("读取消息语法错误: %v", err)
		return 0, err
	}
	//逻辑判断
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		klog.Errorf("消息类型错误'%s'", msg.Operation)
		return 0, fmt.Errorf("消息类型错误'%s'", msg.Operation)
	}
}

// 写数据的方法,拿到apiserver的返回内容，向web端输出
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		klog.Errorf("写消息语法错误: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		klog.Errorf("写消息错误: %v", err)
		return 0, err
	}

	return len(p), nil
}

// 标记关闭的方法
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// 关闭的方法
func (t *TerminalSession) Close() {
	t.wsConn.Close()
}

// resize方法，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

//	func (p *PodClient) GetShell(namespace, podname, cname string, conn *websocket.Conn, conf *rest.Config) error {
//		pty := &TerminalSession{
//			wsConn:   conn,
//			sizeChan: make(chan remotecommand.TerminalSize),
//			doneChan: make(chan struct{}),
//		}
//		defer pty.Close()
//
//		// 1. 检查 /bin/bash 是否存在
//		checkCmd := []string{"/bin/sh", "-c", "if [ -e /bin/bash ]; then echo ok; else echo no; fi"}
//
//		req := p.client.CoreV1().RESTClient().Post().
//			Resource("pods").
//			Name(podname).
//			Namespace(namespace).
//			SubResource("exec").
//			VersionedParams(&corev1.PodExecOptions{
//				Stdin:     false,
//				Stdout:    true,
//				Stderr:    true,
//				TTY:       false,
//				Container: cname,
//				Command:   checkCmd,
//			}, scheme.ParameterCodec)
//
//		executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
//		if err != nil {
//			return err
//		}
//
//		var outputBuf bytes.Buffer
//		err = executor.Stream(remotecommand.StreamOptions{
//			Stdout: &outputBuf,
//			Stderr: &outputBuf,
//		})
//		if err != nil {
//			return err
//		}
//
//		shell := "/bin/sh"
//		if strings.Contains(outputBuf.String(), "ok") {
//			shell = "/bin/bash"
//		}
//
//		// 2. 构造统一提示符，显示真实路径
//		// Bash 使用 $PWD 展示真实路径，sh/ash 使用 $(pwd)
//		var ps1Cmd string
//		if shell == "/bin/bash" {
//			ps1Cmd = fmt.Sprintf("export PS1='root@%s:$PWD# '; exec bash --noprofile --norc", podname)
//		} else {
//			ps1Cmd = fmt.Sprintf("export PS1='root@%s:$(pwd)# '; exec sh", podname)
//		}
//
//		// 3. 执行 shell
//		req = p.client.CoreV1().RESTClient().Post().
//			Resource("pods").
//			Name(podname).
//			Namespace(namespace).
//			SubResource("exec").
//			VersionedParams(&corev1.PodExecOptions{
//				Stdin:     true,
//				Stdout:    true,
//				Stderr:    true,
//				TTY:       true,
//				Container: cname,
//				Command:   []string{shell, "-c", ps1Cmd},
//			}, scheme.ParameterCodec)
//
//		executor, err = remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
//		if err != nil {
//			return err
//		}
//
//		// 4. Stream 交互式终端
//		err = executor.Stream(remotecommand.StreamOptions{
//			Stdin:             pty,
//			Stdout:            pty,
//			Stderr:            pty,
//			Tty:               true,
//			TerminalSizeQueue: pty,
//		})
//		if err != nil {
//			klog.Error("执行 pod 命令失败: " + err.Error())
//			pty.Write([]byte("执行 pod 命令失败: " + err.Error()))
//			pty.Done()
//		}
//
//		return nil
//	}
func (p *PodClient) GetShell(namespace, podname, cname string, conn *websocket.Conn, conf *rest.Config) error {
	pty := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	defer pty.Close()

	// 1. 检查 /bin/bash 是否存在
	checkCmd := []string{"/bin/sh", "-c", "if [ -e /bin/bash ]; then echo ok; else echo no; fi"}

	req := p.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Container: cname,
			Command:   checkCmd,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return err
	}

	var outputBuf bytes.Buffer
	err = executor.Stream(remotecommand.StreamOptions{
		Stdout: &outputBuf,
		Stderr: &outputBuf,
	})
	if err != nil {
		return err
	}

	shell := "/bin/sh"
	if strings.Contains(outputBuf.String(), "ok") {
		shell = "/bin/bash"
	}

	// 2. 构造水印信息
	watermark := fmt.Sprintf(`
🌐  已连接到 Pod       : %s
📦  Namespace         : %s
🛠  Container         : %s
`, podname, namespace, cname)

	// 3. 构造 shell 命令，先打印水印再启动交互式 shell
	var cmd string
	if shell == "/bin/bash" {
		cmd = fmt.Sprintf("echo \"%s\"; export PS1='root@%s:$PWD# '; exec bash --noprofile --norc", watermark, podname)
	} else {
		cmd = fmt.Sprintf("echo \"%s\"; export PS1='root@%s:$(pwd)# '; exec sh", watermark, podname)
	}

	// 4. 执行 shell
	req = p.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podname).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: cname,
			Command:   []string{shell, "-c", cmd},
		}, scheme.ParameterCodec)

	executor, err = remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return err
	}

	// 5. Stream 交互式终端
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})
	if err != nil {
		klog.Error("执行 pod 命令失败: " + err.Error())
		pty.Write([]byte("执行 pod 命令失败: " + err.Error()))
		pty.Done()
	}

	return nil
}
