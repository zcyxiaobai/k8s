package views

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog/v2"
)

// ServiceInfo 用于返回给前端的结构体
type ServiceInfo struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Selector   map[string]string `json:"selector"`
	Ports      []ServicePortInfo `json:"ports"`
	Type       string            `json:"type"`
	Labels     map[string]string `json:"labels"`
}

// ServicePortInfo 端口信息
type ServicePortInfo struct {
	Port       int32  `json:"port"`
	TargetPort string `json:"targetPort"`
	NodePort   int32  `json:"nodePort,omitempty"`
	Protocol   string `json:"protocol"`
}

//创建service的接口

type ReqServiceData struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Selectors map[string]string `json:"selectors"`
	Type      string            `json:"type"` // ClusterIP, NodePort, LoadBalancer
	Ports     []PortReq         `json:"ports"`
}

type PortReq struct {
	Protocol   string `json:"protocol"`   // TCP/UDP
	Port       int32  `json:"port"`       // Service 对外端口
	TargetPort int32  `json:"targetPort"` // Pod 容器端口
}

// 修改service的结构体
type UpdateServiceRequest struct {
	Name      string                     `json:"name" binding:"required"`
	Namespace string                     `json:"namespace" binding:"required"`
	Labels    map[string]string          `json:"labels"`
	Type      string                     `json:"type"`
	Ports     []ServicePortUpdateRequest `json:"ports"`
}

type ServicePortUpdateRequest struct {
	Protocol   string `json:"protocol"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
}

// 查看所有的service
func (p *PodView) GetServiceInfo(c *gin.Context) {
	svcList, err := p.PodClient.ServiceAllList(c)
	if err != nil {
		klog.Error("获取service失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "获取service失败",
		})
		return
	}

	var result []ServiceInfo
	for _, svc := range svcList.Items {
		info := ServiceInfo{
			APIVersion: "v1",
			Kind:       "Service",
			Name:       svc.Name,
			Namespace:  svc.Namespace,
			Selector:   svc.Spec.Selector,
			Type:       string(svc.Spec.Type),
			Labels:     svc.Labels, // 新增：返回 Service 的 labels
		}

		// 端口信息
		for _, port := range svc.Spec.Ports {
			pi := ServicePortInfo{
				Port:     port.Port,
				Protocol: string(port.Protocol),
			}

			// targetPort 可能是 IntOrString 类型，这里转成 string 返回
			if port.TargetPort.Type == intstr.Int {
				pi.TargetPort = fmt.Sprintf("%d", port.TargetPort.IntVal)
			} else {
				pi.TargetPort = port.TargetPort.StrVal
			}

			// NodePort 只有 NodePort / LoadBalancer 类型才有
			if port.NodePort != 0 {
				pi.NodePort = port.NodePort
			}

			info.Ports = append(info.Ports, pi)
		}

		result = append(result, info)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取service成功",
		"date":    result,
	})
}

// 删除指定的service
func (p *PodView) DeleteService(c *gin.Context) {
	ressvc := ReqPodDelete{}
	//接收的数据做反序列化
	err := c.ShouldBind(&ressvc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		klog.Error(err.Error())
		return
	}
	svc := &corev1.Service{}
	svc.Name = ressvc.Name
	svc.Namespace = ressvc.Namespace
	//删除svc
	err = p.PodClient.ServiceDelete(svc.Namespace, svc.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除service失败",
			"poderr":  err.Error(),
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除service成功",
	})
}

// 创建service
func (p *PodView) CreateService(c *gin.Context) {
	req := ReqServiceData{}
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "序列化失败",
		})
		return
	}
	// 构造 Service 对象
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    req.Labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: req.Selectors,
			Type:     corev1.ServiceType(req.Type),
		},
	}
	for _, p := range req.Ports {
		svcPort := corev1.ServicePort{
			Protocol:   corev1.Protocol(p.Protocol),
			Port:       p.Port,
			TargetPort: intstr.FromInt(int(p.TargetPort)),
		}
		svc.Spec.Ports = append(svc.Spec.Ports, svcPort)
	}
	//创建service
	err = p.PodClient.ServiceCreate(svc)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建service失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "创建service成功",
	})

}

// 修改service
func (p *PodView) UpdateService(c *gin.Context) {
	req := UpdateServiceRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		return

	}
	svc := &corev1.Service{}
	svc.Name = req.Name
	svc.Namespace = req.Namespace
	//获取指定的service
	svcnew, err := p.PodClient.ServiceGet(svc.Namespace, svc.Name)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取service失败",
		})
		return
	}
	// 更新 Labels
	if svcnew.Labels != nil {
		svcnew.Labels = make(map[string]string)
	} else {
		for k := range svcnew.Labels {
			//清空旧的
			delete(svcnew.Labels, k)
		}
	}
	//添加labels
	for k, v := range req.Labels {
		svcnew.Labels[k] = v
	}
	//更新service类型
	if req.Type != "" {
		svcnew.Spec.Type = corev1.ServiceType(req.Type)
	}
	// 更新 Ports
	if len(req.Ports) > 0 {
		newPorts := []corev1.ServicePort{}
		for _, p := range req.Ports {
			newPorts = append(newPorts, corev1.ServicePort{
				Name:       "",
				Protocol:   corev1.Protocol(p.Protocol),
				Port:       p.Port,
				TargetPort: intstr.FromInt(int(p.TargetPort)),
			})
		}
		svcnew.Spec.Ports = newPorts
	}
	//提交更新
	err = p.PodClient.ServiceUpdate(svc.Namespace, svc.Name, svcnew)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改service失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改service成功",
	})
}
