package k8s

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// pod的增删改查
// 定义pod的客户端
type PodClient struct {
	client        *kubernetes.Clientset
	MetricsClient *metricsv.Clientset
}

// 获取k8s中的数据需要一个客户端  这个方法就是返回一个客户端 这是一个 构造函数（工厂方法），用来创建 PodClient 的实例
func NewPodClient(client *kubernetes.Clientset, metricsClient *metricsv.Clientset) *PodClient {
	return &PodClient{
		client:        client,
		MetricsClient: metricsClient,
	}
}

//func NewMetricsClient(client *metricsv.Clientset) *MetricsClient {
//	return &MetricsClient{client}
//}

// 作用：查询某个 namespace 下的所有 Pod，并返回 Pod 列表p.client.CoreV1().Pods(namespace)：
// 使用 Kubernetes API 的 CoreV1 组（核心 API，包括 Pod、Node、Service 等），获取某个命名空间下的 Pod 资源
// metav1.ListOptions{} 是 Kubernetes List/Get 请求的参数，里面的条件可以用来做过滤、分页、字段选择等，相当于 kubectl get 的一些选项
// 查询pod
func (p *PodClient) List(namespace string) ([]corev1.Pod, error) {
	podlist, err := p.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	//直接返回，调用值再进行判断
	return podlist.Items, err
}

// 查询指定的单个pod
func (p *PodClient) Get(namespace, name string) (*corev1.Pod, error) {
	podyx, err := p.client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	//返回查找到的pod和err
	return podyx, err
}

// 创建pod  指向的类型是 corev1.Pod（Pod 资源对象的结构体）
func (p *PodClient) Create(pod *corev1.Pod) error {
	//Kubernetes 官方 Go Client（client-go）库里内置提供的,用于创建指定的pod
	_, err := p.client.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	return err
}

// 删除pod
func (p *PodClient) Delete(namespace, podname string) error {
	err := p.client.CoreV1().Pods(namespace).Delete(context.TODO(), podname, metav1.DeleteOptions{})
	return err
}

// 更新修改pod
func (p *PodClient) Update(namespace, name string, pod *corev1.Pod) error {
	pod.Name = name // 确保一致
	_, err := p.client.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	return err
}

// 获取node的数量
func (p *PodClient) GetNode(c *gin.Context) (*corev1.NodeList, error) {
	nodes, err := p.client.CoreV1().Nodes().List(c.Request.Context(), metav1.ListOptions{})
	return nodes, err
}

// 获取node Metrics
func (p *PodClient) GetNodeMetrics(c *gin.Context) (*v1beta1.NodeMetricsList, error) {
	nodeMetrics, err := p.MetricsClient.MetricsV1beta1().NodeMetricses().List(c.Request.Context(), metav1.ListOptions{})
	return nodeMetrics, err
}

// 获取所有的名称空间名称
func (p *PodClient) GetNameSpace(c *gin.Context) (*corev1.NamespaceList, error) {
	namespacelist, err := p.client.CoreV1().Namespaces().List(c.Request.Context(), metav1.ListOptions{})
	return namespacelist, err
}

// 按照节点查询pod
func (p *PodClient) GetNodePod(c *gin.Context, name string) (*corev1.PodList, error) {
	podlist, err := p.client.CoreV1().Pods("").List(c.Request.Context(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", name),
	})
	return podlist, err
}

// 获取pod的总数量
func (p *PodClient) GetPodAll(c *gin.Context) (*corev1.PodList, error) {
	number, err := p.client.CoreV1().Pods("").List(c.Request.Context(), metav1.ListOptions{})
	return number, err
}

// 获取deployment的总数量
func (p *PodClient) GetDepAll(c *gin.Context) (int32, error) {
	number, err := p.client.AppsV1().Deployments("").List(c.Request.Context(), metav1.ListOptions{})
	return int32(len(number.Items)), err
}

// 获取service的总数量
func (p *PodClient) GetSvcAll(c *gin.Context) (int32, error) {
	number, err := p.client.CoreV1().Services("").List(c.Request.Context(), metav1.ListOptions{})
	return int32(len(number.Items)), err
}

// 获取ingress的总数量
func (p *PodClient) GetIngAll(c *gin.Context) (int32, error) {
	number, err := p.client.NetworkingV1().Ingresses("").List(c.Request.Context(), metav1.ListOptions{})
	return int32(len(number.Items)), err
}
