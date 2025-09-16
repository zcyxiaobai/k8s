package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// IngressInfo 用于序列化返回前端的数据
type IngressInfo struct {
	APIVersion       string             `json:"apiVersion"`
	Kind             string             `json:"kind"`
	Name             string             `json:"name"`
	Namespace        string             `json:"namespace"`
	Labels           map[string]string  `json:"labels,omitempty"`
	Annotations      map[string]string  `json:"annotations,omitempty"`
	IngressClassName string             `json:"ingressClassName,omitempty"`
	Rules            []IngressRuleInfo  `json:"rules,omitempty"`
	TLS              []IngressTLSInfo   `json:"tls,omitempty"`
	Status           []LoadBalancerInfo `json:"status,omitempty"`
}

// IngressRuleInfo 路由规则
type IngressRuleInfo struct {
	Host  string            `json:"host,omitempty"`
	Paths []IngressPathInfo `json:"paths,omitempty"`
}

// IngressPathInfo 路径配置
type IngressPathInfo struct {
	Path     string `json:"path"`
	PathType string `json:"pathType"`
	Service  string `json:"service"`
	Port     int32  `json:"port"`
}

// IngressTLSInfo TLS 配置
type IngressTLSInfo struct {
	Hosts      []string `json:"hosts,omitempty"`
	SecretName string   `json:"secretName,omitempty"`
}

// LoadBalancerInfo 负载均衡信息
type LoadBalancerInfo struct {
	IP       string `json:"ip,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

// 修改ingress的结构体
type UpdateIngressRequest struct {
	Name        string              `json:"name" binding:"required"`
	Namespace   string              `json:"namespace" binding:"required"`
	Labels      map[string]string   `json:"labels"`
	Annotations map[string]string   `json:"annotations"`
	Rules       []IngressRuleUpdate `json:"rules"`
}

type IngressRuleUpdate struct {
	Host  string              `json:"host"`
	Paths []IngressPathUpdate `json:"paths"`
}

type IngressPathUpdate struct {
	Path        string `json:"path"`
	PathType    string `json:"pathType"`
	ServiceName string `json:"serviceName"`
	ServicePort int32  `json:"servicePort"`
}

// 创建ingress的结构体
// 请求结构体
type ReqIngressData struct {
	Name             string            `json:"name" binding:"required"`
	Namespace        string            `json:"namespace" binding:"required"`
	IngressClassName string            `json:"ingressClassName"`
	Annotations      map[string]string `json:"annotations"`
	Rules            []struct {
		Host  string `json:"host"`
		Paths []struct {
			Path        string `json:"path"`
			PathType    string `json:"pathType"`
			ServiceName string `json:"serviceName"`
			ServicePort int    `json:"servicePort"`
		} `json:"paths"`
	} `json:"rules"`
}

func (p *PodView) GetIngressAll(c *gin.Context) {
	ingList, err := p.PodClient.IngressAllLst(c)
	if err != nil {
		klog.Error("查询ingress失败")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询 ingress 失败",
		})
		return
	}
	var result []IngressInfo
	for _, ing := range ingList.Items {
		info := IngressInfo{
			APIVersion:       "networking.k8s.io/v1",
			Kind:             "Ingress",
			Name:             ing.Name,
			Namespace:        ing.Namespace,
			Labels:           ing.Labels,
			Annotations:      ing.Annotations,
			IngressClassName: "",
		}

		if ing.Spec.IngressClassName != nil {
			info.IngressClassName = *ing.Spec.IngressClassName
		}

		// 规则
		for _, rule := range ing.Spec.Rules {
			ruleInfo := IngressRuleInfo{
				Host: rule.Host,
			}
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					pathInfo := IngressPathInfo{
						Path:     path.Path,
						PathType: string(*path.PathType),
					}
					if path.Backend.Service != nil {
						pathInfo.Service = path.Backend.Service.Name
						if path.Backend.Service.Port.Number != 0 {
							pathInfo.Port = path.Backend.Service.Port.Number
						}
					}
					ruleInfo.Paths = append(ruleInfo.Paths, pathInfo)
				}
			}
			info.Rules = append(info.Rules, ruleInfo)
		}

		// TLS
		for _, tls := range ing.Spec.TLS {
			tlsInfo := IngressTLSInfo{
				Hosts:      tls.Hosts,
				SecretName: tls.SecretName,
			}
			info.TLS = append(info.TLS, tlsInfo)
		}

		// 状态
		for _, lb := range ing.Status.LoadBalancer.Ingress {
			info.Status = append(info.Status, LoadBalancerInfo{
				IP:       lb.IP,
				Hostname: lb.Hostname,
			})
		}

		result = append(result, info)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取ingress成功",
		"date":    result,
	})

}

//删除指定ingress

func (p *PodView) IngressDelete(c *gin.Context) {
	resing := ReqPodDelete{}
	//接收的数据做反序列化
	err := c.ShouldBind(&resing)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		klog.Error(err.Error())
		return
	}
	ing := &networkingv1.Ingress{}
	ing.Name = resing.Name
	ing.Namespace = resing.Namespace
	//删除ing
	err = p.PodClient.DeleteIngress(ing.Namespace, ing.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除ingress失败",
			"poderr":  err.Error(),
		})
		klog.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除ingress成功",
	})
}

// 创建ingrss
func (p *PodView) CreateIngress(c *gin.Context) {
	req := ReqIngressData{}
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "发序列化失败",
		})
		return
	}
	ingress := networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Annotations: req.Annotations,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &req.IngressClassName,
		},
	}
	var rules []networkingv1.IngressRule
	for _, r := range req.Rules {
		var paths []networkingv1.HTTPIngressPath
		for _, p := range r.Paths {
			pathType := networkingv1.PathType(p.PathType)
			paths = append(paths, networkingv1.HTTPIngressPath{
				Path:     p.Path,
				PathType: &pathType,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: p.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Number: int32(p.ServicePort),
						},
					},
				},
			})
		}
		rules = append(rules, networkingv1.IngressRule{
			Host: r.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		})
	}
	ingress.Spec.Rules = rules
	err = p.PodClient.IngressCreate(ingress)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建ingress失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "创建ingress成功",
	})

}

// 修改ingress
func (p *PodView) UpdateIngress(c *gin.Context) {
	req := UpdateIngressRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "反序列化失败",
		})
		return
	}
	ingress := networkingv1.Ingress{}
	ingress.Name = req.Name
	ingress.Namespace = req.Namespace
	//获取相应的ingress
	ingnew, err := p.PodClient.IngressGet(ingress.Namespace, ingress.Name)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取ingress失败",
		})
		return
	}
	// 更新 Labels（如果请求不为空，先清空再添加）
	if ingnew.Labels == nil {
		ingnew.Labels = make(map[string]string)
	} else {
		for k := range ingnew.Labels {
			delete(ingnew.Labels, k) // 清空旧的
		}
	}
	for k, v := range req.Labels {
		ingnew.Labels[k] = v
	}
	// 更新 Annotations
	if ingnew.Annotations == nil {
		ingnew.Annotations = make(map[string]string)
	} else {
		for k := range ingnew.Annotations {
			delete(ingnew.Annotations, k)
		}
	}
	for k, v := range req.Annotations {
		ingnew.Annotations[k] = v
	}
	//更新rules
	ingnew.Spec.Rules = []networkingv1.IngressRule{} // 直接清空 slice
	for _, r := range req.Rules {
		paths := []networkingv1.HTTPIngressPath{}
		for _, p := range r.Paths {
			pt := networkingv1.PathType(p.PathType)
			paths = append(paths, networkingv1.HTTPIngressPath{
				Path:     p.Path,
				PathType: &pt,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: p.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Number: p.ServicePort,
						},
					},
				},
			})
		}
		ingnew.Spec.Rules = append(ingnew.Spec.Rules, networkingv1.IngressRule{
			Host: r.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: paths,
				},
			},
		})
	}
	//更新ingress
	err = p.PodClient.IngressUpdate(ingress.Namespace, ingress.Name, ingnew)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改ingress失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改ingress成功",
	})
}
