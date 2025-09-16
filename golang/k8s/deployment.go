package k8s

import (
	"context"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 查询所有的deployment
func (p *PodClient) DeploymentAllList(c *gin.Context) (*appsv1.DeploymentList, error) {
	deplist, err := p.client.AppsV1().Deployments("").List(c.Request.Context(), metav1.ListOptions{})
	return deplist, err
}

// 删除指定的deployment
func (p *PodClient) DeploymentDelete(namespace, depname string) error {
	err := p.client.AppsV1().Deployments(namespace).Delete(context.TODO(), depname, metav1.DeleteOptions{})
	return err
}

// 创建新增deployment
func (p *PodClient) DeploymentCraete(deployment *appsv1.Deployment) error {
	_, err := p.client.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

// 查询指定的deployment
func (p *PodClient) DeploymentGet(namespace, depname string) (*appsv1.Deployment, error) {
	dep, err := p.client.AppsV1().Deployments(namespace).Get(context.TODO(), depname, metav1.GetOptions{})
	return dep, err
}

// 修改deployment
func (p *PodClient) DeploymentUpdate(namespace, name string, deployment *appsv1.Deployment) error {
	deployment.Name = name // 确保一致
	_, err := p.client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}
