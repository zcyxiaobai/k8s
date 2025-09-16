package views

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"y2505.com/bookapp/k8s"
)

// 定义一个pod的结构体，需要一个pod的客户端
type PodView struct {
	PodClient *k8s.PodClient
}

//type MetricsView struct {
//	MetricsClient *k8s.MetricsClient
//}

// 定义pod的结构体，返回的pod的数据模板
type PodData struct {
	Namespace    string `json:"namespace"`
	Name         string `json:"name"`
	ContainerNum int    `json:"containernum"`
}

// 定义一个结构体，声明pod镜像的相关属性
type Container struct {
	Name            string `json:"name"`
	Image           string `json:"image"`
	ImagePullPolicy string `json:"imagepullpolicy"`
	Port            int32  `json:"port"`
}

// 定义一个结构体，用于创建pod
type ReqPodData struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Labels     map[string]string `json:"labels"`
	Containers []Container       `json:"containers"`
}

// 定义一个结构体，用于返回查找到的单个pod的关键信息
//type ReqGetPodDate struct {
//}

// 定义一个结构体，用于删除和查询单个pod
type ReqPodDelete struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// 定义节点监控信息响应结构体
type NodeMetrics struct {
	Name       string `json:"name"`
	CPUUsage   string `json:"cpuUsage"`
	CPUPercent string `json:"cpuPercent"`
	MemUsage   string `json:"memUsage"`
	MemPercent string `json:"memPercent"`
}

// Pod 状态常量
const (
	PodStatusRunning   = "Running"
	PodStatusFailed    = "Failed"
	PodStatusPending   = "Pending"
	PodStatusSucceeded = "Succeeded"
)

// NodePodStats 响应结构
type NodePodStats struct {
	NodeName       string `json:"nodeName"`
	RunningPods    int    `json:"runningPods"`
	NonRunningPods int    `json:"nonRunningPods"`
}

// 定义修改pod的结构体
// UpdatePodRequest 定义前端传入的数据结构
type UpdatePodRequest struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

// 这个方法是用于前端请求返回数据的
func (p *PodView) List(c *gin.Context) {
	//获取前端传过来的名称空间
	namespace := c.Param("namespace")
	//得到所有满足要求的pod,在这里我就可以修改获取pod的各种不同属性
	pods, err := p.PodClient.List(namespace)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found pods",
		})
		return
	}
	//[]PodData这是一个PodData结构体类型的切片
	podsdata := make([]PodData, len(pods))
	//遍历返回pod
	for i, pod := range pods {
		podsdata[i] = PodData{
			Namespace:    pod.Namespace,
			Name:         pod.Name,
			ContainerNum: len(pod.Spec.Containers),
		}
	}
	//将
	c.JSON(http.StatusOK, gin.H{
		"data": podsdata,
	})

}

// 查找单个的pod
func (p *PodView) GetPod(c *gin.Context) {
	//用于接收用户输入的pod的信息（namespace 和name）
	resPod := ReqPodDelete{}
	//用于输出查询到的pod的关键信息
	newPod := ReqPodData{}
	//接收的数据做反序列化
	err := c.ShouldBind(&resPod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "序列化pod失败",
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	pod := &corev1.Pod{}
	pod.Name = resPod.Name
	pod.Namespace = resPod.Namespace
	podyx, err := p.PodClient.Get(pod.Namespace, resPod.Name)
	newPod.Name = podyx.Name
	newPod.Namespace = podyx.Namespace
	newPod.Labels = podyx.Labels
	containers := make([]Container, len(podyx.Spec.Containers))
	for i, container := range podyx.Spec.Containers {
		containers[i].Name = container.Name
		containers[i].Image = container.Image
		containers[i].ImagePullPolicy = string(container.ImagePullPolicy)
	}
	newPod.Containers = containers
	//只有返回需要的数据
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "查询pod失败",
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询pod成功",
		"podyx":   newPod,
	})
}

// 获取容器日志
//
//	func (p *PodView) GetLog(c *gin.Context) {
//		//将http服务升级为WebSocket
//		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
//		if err != nil {
//			klog.Error(err)
//			c.JSON(http.StatusInternalServerError, gin.H{
//				"message": "websocket upgrade failed",
//			})
//			return
//		}
//		//获取用户传入的参数
//		namespace := c.Param("namespace")
//		podname := c.Param("podname")
//		cname := c.Param("cname")
//
//		err = p.PodClient.GetPodLog(namespace, podname, cname, conn)
//		if err != nil {
//			conn.WriteMessage(websocket.TextMessage, []byte("获取日志失败: "+err.Error()))
//			return
//		}
//
// }

// 创建新增pod
func (p *PodView) CreatePod(c *gin.Context) {
	resPod := ReqPodData{}
	//接收的数据做反序列化
	err := c.ShouldBind(&resPod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "序列化pod失败",
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	pod := &corev1.Pod{}
	pod.Name = resPod.Name
	pod.Namespace = resPod.Namespace
	pod.Labels = resPod.Labels
	containers := make([]corev1.Container, len(resPod.Containers))
	for i, container := range resPod.Containers {
		containers[i].Name = container.Name
		containers[i].Image = container.Image
		containers[i].ImagePullPolicy = corev1.PullPolicy(container.ImagePullPolicy)
		containers[i].Ports = []corev1.ContainerPort{{ContainerPort: container.Port}}
	}
	pod.Spec.Containers = containers
	//创建pod
	err = p.PodClient.Create(pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "创建pod失败",
		})
		klog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "创建pod成功",
	})
}

// 删除pod
func (p *PodView) Delete(c *gin.Context) {
	resPod := ReqPodDelete{}
	//接收的数据做反序列化
	err := c.ShouldBind(&resPod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	pod := &corev1.Pod{}
	pod.Name = resPod.Name
	pod.Namespace = resPod.Namespace
	//删除pod
	err = p.PodClient.Delete(pod.Namespace, pod.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "删除pod失败",
			"poderr":  err.Error(),
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除pod成功",
	})
}

// 更新pod,功能未完成
func (p *PodView) Update(c *gin.Context) {
	resPod := ReqPodData{}
	err := c.ShouldBind(&resPod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	pod := &corev1.Pod{}
	pod.Name = resPod.Name
	pod.Namespace = resPod.Namespace
	pod.Labels = resPod.Labels
	containers := make([]corev1.Container, len(resPod.Containers))
	for i, container := range resPod.Containers {
		containers[i].Name = container.Name
		containers[i].Image = container.Image
		containers[i].ImagePullPolicy = corev1.PullPolicy(container.ImagePullPolicy)

	}
	pod.Spec.Containers = containers
	err = p.PodClient.Update(pod.Namespace, pod.Name, pod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改pod失败,有些字段无法修改",
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改pod成功",
	})

}

// 获取k8s集群中各个资源的总数
func (p *PodView) GetResourceCount(c *gin.Context) {

	//获取Node的数量
	nodes, err := p.PodClient.GetNode(c)
	if err != nil {
		klog.Error("获取 Node 失败: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取node失败",
		})
		return
	}
	//获取pod的数量
	pods, err := p.PodClient.GetPodAll(c)
	if err != nil {
		klog.Error("获取pod 失败: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取pod失败",
		})
		return
	}
	//获取deployment的数量
	deps, err := p.PodClient.GetDepAll(c)
	if err != nil {
		klog.Error("获取deployment失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取deployment失败",
		})
		return
	}
	//获取service的数量
	svcs, err := p.PodClient.GetSvcAll(c)
	if err != nil {
		klog.Error("获取service失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取service失败",
		})
		return
	}
	//获取ingress的数量
	ings, err := p.PodClient.GetIngAll(c)
	if err != nil {
		klog.Error("获取ingress失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取ingress失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据查询成功",
		"node":    int32(len(nodes.Items)),
		"pods":    int32(len(pods.Items)),
		"deps":    deps,
		"svcs":    svcs,
		"ings":    ings,
	})
}

// 获取所有节点的监控信息
func (p *PodView) GetNodeMetrics(c *gin.Context) {
	//获取Node列表
	nodes, err := p.PodClient.GetNode(c)
	if err != nil {
		klog.Error("获取Node失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取Node失败",
		})
		return
	}
	//获取Node Metrics
	nodemetrics, err := p.PodClient.GetNodeMetrics(c)
	if err != nil {
		klog.Error("获取NodeMetrics失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取NodeMetrics失败",
		})
		return
	}
	//转换结果
	var results []NodeMetrics
	for _, nm := range nodemetrics.Items {
		// 找到对应 Node 信息（为了获取 capacity）
		var node corev1.Node
		for _, n := range nodes.Items {
			if n.Name == nm.Name {
				node = n
				break
			}
		}
		cpuUsage := nm.Usage["cpu"]
		memUsage := nm.Usage["memory"]

		cpuCapacity := node.Status.Capacity["cpu"]
		memCapacity := node.Status.Capacity["memory"]
		// 转换为数值（m 核 & MiB）
		cpuUsageMillicores := cpuUsage.MilliValue()
		cpuCapacityMillicores := cpuCapacity.MilliValue()
		cpuPercent := float64(cpuUsageMillicores) / float64(cpuCapacityMillicores) * 100

		memUsageBytes := memUsage.Value()
		memCapacityBytes := memCapacity.Value()
		memPercent := float64(memUsageBytes) / float64(memCapacityBytes) * 100
		results = append(results, NodeMetrics{
			Name:       nm.Name,
			CPUUsage:   cpuUsage.String(),
			CPUPercent: formatPercent(cpuPercent),
			MemUsage:   memUsage.String(),
			MemPercent: formatPercent(memPercent),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取 Node Metrics 成功",
		"data":    results,
	})
}
func formatPercent(p float64) string {
	return fmt.Sprintf("%.2f%%", p)
}

// 获取每个node上不同运行状态的pod
func (p *PodView) GetNodePodStatus(c *gin.Context) {
	//获取node列表
	nodes, err := p.PodClient.GetNode(c)
	if err != nil {
		klog.Error("获取Node失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取Node失败",
		})
		return
	}
	// 存储每个节点的 Pod 状态统计
	var nodePodStats []NodePodStats
	// 遍历每个 Node，获取该节点上的 Pods 状态
	for _, n := range nodes.Items {
		//获取当前节点上的Pod列表
		pods, err := p.PodClient.GetNodePod(c, n.Name)
		if err != nil {
			klog.Error("获取Pod失败", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("获取节点 %s 的 Pod 失败", n.Name),
			})
			return
		}
		// 统计当前节点上正常和非正常的 Pod 数量
		var runningCount, failedCount int
		for _, pod := range pods.Items {
			// 判断 Pod 是否处于 Running 状态
			if pod.Status.Phase == PodStatusRunning {
				runningCount++
			} else {
				failedCount++
			}
		}
		// 将结果添加到节点 Pod 状态统计中
		nodePodStats = append(nodePodStats, NodePodStats{
			NodeName:       n.Name,
			RunningPods:    runningCount,
			NonRunningPods: failedCount,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取节点的pod成功",
		"data":    nodePodStats,
	})
}

// 获取所有的名称空间
func (p *PodView) GetNamespaceList(c *gin.Context) {
	namespacelist, err := p.PodClient.GetNameSpace(c)
	if err != nil {
		klog.Error("获取Namespace失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取namespace失败",
		})
		return
	}
	var namespaces []string
	for _, n := range namespacelist.Items {
		namespaces = append(namespaces, n.Name)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取namespace成功",
		"data":    namespaces,
	})
}

// 获取集群中所有pod的详细信息
func (p *PodView) GetPodAllList(c *gin.Context) {
	podlist, err := p.PodClient.GetPodAll(c)
	if err != nil {
		klog.Error("获取pod失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取pod失败",
		})
		return
	}

	result := []gin.H{}
	for _, pod := range podlist.Items {
		// 解析容器信息
		containers := []gin.H{}
		for _, ctn := range pod.Spec.Containers {
			// 获取端口信息
			ports := []int32{}
			for _, port := range ctn.Ports {
				ports = append(ports, port.ContainerPort)
			}

			containers = append(containers, gin.H{
				"name":            ctn.Name,
				"image":           ctn.Image,
				"imagePullPolicy": string(ctn.ImagePullPolicy),
				"ports":           ports,
			})
		}

		// 拼接 Pod 信息
		result = append(result, gin.H{
			"name":       pod.Name,
			"status":     string(pod.Status.Phase),
			"namespace":  pod.Namespace,
			"node":       pod.Spec.NodeName,
			"labels":     pod.Labels,
			"containers": containers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取pod成功",
		"data":    result,
	})
}

// 修改pod
func (p *PodView) UpdatePod(c *gin.Context) {
	resPod := UpdatePodRequest{}
	err := c.ShouldBind(&resPod)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		return
	}

	pod := &corev1.Pod{}
	pod, err = p.PodClient.Get(resPod.Namespace, resPod.Name)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取pod失败",
		})
		return
	}
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	} else {
		for k := range pod.Labels {
			//清空旧的pod
			delete(pod.Labels, k)
		}
	}
	// 添加新 labels
	for k, v := range resPod.Labels {
		pod.Labels[k] = v
	}
	//更新pod
	err = p.PodClient.Update(pod.Namespace, pod.Name, pod)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "更新pod失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更新pod成功",
	})
}
