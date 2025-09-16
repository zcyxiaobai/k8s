package k8s

import (
	"context"

	"github.com/gin-gonic/gin"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 获取所有ingress
func (p *PodClient) IngressAllLst(c *gin.Context) (*networkingv1.IngressList, error) {
	ingList, err := p.client.NetworkingV1().Ingresses("").List(c.Request.Context(), metav1.ListOptions{})
	return ingList, err
}

// 删除指定的ingress
func (p *PodClient) DeleteIngress(namespace, ingname string) error {
	err := p.client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingname, metav1.DeleteOptions{})
	return err
}

// 创建新增ingress
func (p *PodClient) IngressCreate(ingress networkingv1.Ingress) error {
	_, err := p.client.NetworkingV1().Ingresses(ingress.Namespace).Create(context.TODO(), &ingress, metav1.CreateOptions{})
	return err
}

// 查询指定的ingress
func (p *PodClient) IngressGet(namespace, name string) (*networkingv1.Ingress, error) {
	ing, err := p.client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return ing, err
}

// 修改ingress
func (p *PodClient) IngressUpdate(namespace, name string, ingress *networkingv1.Ingress) error {
	ingress.Name = name
	_, err := p.client.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	return err
}
