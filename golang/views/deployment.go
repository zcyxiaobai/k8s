package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type ContainerInfo struct {
	Name            string  `json:"name"`
	Image           string  `json:"image"`
	Ports           []int32 `json:"ports,omitempty"`
	ImagePullPolicy string  `json:"imagePullPolicy,omitempty"` // 新增字段
}

// 自定义返回结构
type DeploymentInfo struct {
	APIVersion       string            `json:"apiVersion"`
	Kind             string            `json:"kind"`
	Name             string            `json:"name"`
	Namespace        string            `json:"namespace"`
	Replicas         int32             `json:"replicas"`
	ImagePullSecrets []string          `json:"imagePullSecrets"`
	Containers       []ContainerInfo   `json:"containers"`
	Labels           map[string]string `json:"labels"`   // Pod 模板 labels
	Selector         map[string]string `json:"selector"` // selector.matchLabels
}

// ReqDeploymentData 定义前端传入的 Deployment 数据
type ReqDeploymentData struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Replicas   int32             `json:"replicas"`
	MetaLabels map[string]string `json:"metaLabels"`
	Labels     map[string]string `json:"labels"`
	Containers []ContainerReq    `json:"containers"`
}

// ContainerReq 定义容器请求结构体
type ContainerReq struct {
	Name            string `json:"name"`
	Image           string `json:"image"`
	ImagePullPolicy string `json:"imagePullPolicy"`
	Port            int32  `json:"port"`
}

// 定义修改deployment的结构体
// UpdateDeploymentRequest 定义请求体
type UpdateDeploymentRequest struct {
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Replicas  int32             `json:"replicas"`
}

// 查询所有的deployment
func (p *PodView) GetDeploymentList(c *gin.Context) {
	deplist, err := p.PodClient.DeploymentAllList(c)
	if err != nil {
		klog.Error("查询deployment失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询deployment失败",
		})
		return
	}
	var result []DeploymentInfo
	for _, dep := range deplist.Items {
		info := DeploymentInfo{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       dep.Name,
			Namespace:  dep.Namespace,
			Replicas:   1,
			Labels:     dep.Spec.Template.Labels,
			Selector:   dep.Spec.Selector.MatchLabels,
		}

		// 副本数
		if dep.Spec.Replicas != nil {
			info.Replicas = *dep.Spec.Replicas
		}

		// 容器信息
		for _, ctn := range dep.Spec.Template.Spec.Containers {
			ci := ContainerInfo{
				Name:            ctn.Name,
				Image:           ctn.Image,
				ImagePullPolicy: string(ctn.ImagePullPolicy), // 这里获取镜像拉取策略
			}
			for _, port := range ctn.Ports {
				ci.Ports = append(ci.Ports, port.ContainerPort)
			}
			info.Containers = append(info.Containers, ci)
		}

		// 镜像拉取 Secret
		for _, secret := range dep.Spec.Template.Spec.ImagePullSecrets {
			info.ImagePullSecrets = append(info.ImagePullSecrets, secret.Name)
		}

		result = append(result, info)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询deployment成功",
		"date":    result,
	})
}

// 删除指定的deployment
func (p *PodView) DeploymentDelete(c *gin.Context) {
	resdev := ReqPodDelete{}
	//接收的数据做反序列化
	err := c.ShouldBind(&resdev)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		klog.Error(err.Error())
		return
	}
	dev := &appsv1.Deployment{}
	dev.Name = resdev.Name
	dev.Namespace = resdev.Namespace
	//删除dev
	err = p.PodClient.DeploymentDelete(dev.Namespace, resdev.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除dev失败",
			"poderr":  err.Error(),
		})
		//打印错误日志
		klog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除dev成功",
	})
}

// 创建deployment
func (p *PodView) CreateDeployment(c *gin.Context) {
	req := ReqDeploymentData{}
	// 反序列化
	if err := c.ShouldBind(&req); err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "序列化 deployment 失败",
		})
		return
	}
	// 构造 Deployment 对象
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    req.MetaLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &req.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: req.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: req.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: buildContainers(req.Containers),
				},
			},
		},
	}
	err := p.PodClient.DeploymentCraete(deployment)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建deployment失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"meaage": "创建deployment成功",
	})

}

// 辅助函数：构建容器列表
func buildContainers(cs []ContainerReq) []corev1.Container {
	containers := make([]corev1.Container, len(cs))
	for i, c := range cs {
		containers[i] = corev1.Container{
			Name:            c.Name,
			Image:           c.Image,
			ImagePullPolicy: corev1.PullPolicy(c.ImagePullPolicy), // 注意这里强转
		}
		if c.Port > 0 {
			containers[i].Ports = []corev1.ContainerPort{
				{ContainerPort: c.Port},
			}
		}
	}
	return containers
}

// 修改deployment
func (p *PodView) UpdateDeployment(c *gin.Context) {
	req := UpdateDeploymentRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		return
	}
	dep := &appsv1.Deployment{}
	dep.Name = req.Name
	dep.Namespace = req.Namespace
	depnew, err := p.PodClient.DeploymentGet(dep.Namespace, dep.Name)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取deployment失败",
		})
		return
	}

	//更新副本数
	// 3. 更新副本数
	if req.Replicas > 0 {
		replicas := req.Replicas
		depnew.Spec.Replicas = &replicas
	}
	// 更新 Labels（metadata.labels）
	if depnew.Labels == nil {
		depnew.Labels = make(map[string]string)
	} else {
		for k := range depnew.Labels {
			delete(depnew.Labels, k) // 清空旧的
		}
	}
	//添加新的标签
	for k, v := range req.Labels {
		depnew.Labels[k] = v
	}
	// 更新 Pod 模板 labels
	depnew.Spec.Template.Labels = make(map[string]string)
	for k, v := range req.Labels {
		depnew.Spec.Template.Labels[k] = v
	}
	//提交更新
	err = p.PodClient.DeploymentUpdate(dep.Namespace, dep.Name, depnew)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改deployment失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改deployment成功",
	})
}
