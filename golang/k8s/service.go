package k8s

import (
	"context"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 查询所有的service
func (p *PodClient) ServiceAllList(c *gin.Context) (*corev1.ServiceList, error) {
	svcList, err := p.client.CoreV1().Services("").List(c.Request.Context(), metav1.ListOptions{})
	return svcList, err
}

// 删除service
func (p *PodClient) ServiceDelete(namespace, svcname string) error {
	err := p.client.CoreV1().Services(namespace).Delete(context.TODO(), svcname, metav1.DeleteOptions{})
	return err
}

// 创建service
func (p *PodClient) ServiceCreate(service *corev1.Service) error {
	_, err := p.client.CoreV1().Services(service.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	return err
}

// 查询指定的service
func (p *PodClient) ServiceGet(namespace, svcname string) (*corev1.Service, error) {
	svc, err := p.client.CoreV1().Services(namespace).Get(context.TODO(), svcname, metav1.GetOptions{})
	return svc, err
}

// 修改service
func (p *PodClient) ServiceUpdate(namespace, name string, service *corev1.Service) error {
	service.Name = name // 确保一致
	_, err := p.client.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	return err
}
