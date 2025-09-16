/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	//"k8s.io/client-go/informers/networking"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "y2505.com/bookapp/api/v1"
)

// BookappReconciler reconciles a Bookapp object
type BookappReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *BookappReconciler) updateBookAppStatus(ctx context.Context, bookapp *webappv1.Bookapp) error {
	dep := appsv1.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      bookapp.Name,
		Namespace: bookapp.Namespace,
	}, &dep)
	if err != nil {
		return err
	}
	running := dep.Status.AvailableReplicas
	total := dep.Status.Replicas
	notRunning := total - running

	bookapp.Status.Runing = &running
	bookapp.Status.NotRuning = &notRunning
	if err := r.Status().Update(ctx, bookapp); err != nil {
		return err
	}
	return nil

}

// +kubebuilder:rbac:groups=webapp.y2505.com,resources=bookapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.y2505.com,resources=bookapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.y2505.com,resources=bookapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bookapp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
// 调和函数，死循环
func (r *BookappReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)
	klog.Infof("Start to reconcile Bookapp")
	//定义一个全局变量，存储返回的数据
	bookapp := webappv1.Bookapp{}
	// TODO(user): your logic here
	// 自己的逻辑代码
	// 拿到key去本地缓存查看是否有指定数据，没有说明已经被删除，如果有数据，就需要执行添加或者修改的逻辑
	//根据 req.NamespacedName（资源的名字和命名空间），获取对应的 Bookapp 对象，并把结果填充到传入的结构体里

	err := r.Client.Get(ctx, req.NamespacedName, &bookapp)
	if err != nil {
		//如果没有找到对应的数据
		if errors.IsNotFound(err) {
			klog.Infof("bookapp:%v-%v is delete\n", req.NamespacedName, req.Name)
			return ctrl.Result{}, nil
		}
		//如果有错就重新入队列，超过5次取消

		return ctrl.Result{}, err
	}
	//有对应的数据，创建和更新的逻辑
	//如果找到了对应的数据，需要判断是否在删除
	//如果这个时间戳不为空，表示该资源正在被删除
	if bookapp.DeletionTimestamp != nil {
		klog.Infof("bookapp:%v-%v is delete\n", req.NamespacedName, req.Name)
		//需要重新入队列
		return ctrl.Result{
			Requeue: true,
		}, nil
	}
	//add bookapp

	//复制一个新的deployment
	//newdep := dep.DeepCopyObject()

	//根据 req.NamespacedName（也就是 namespace + name）来获取对应的 Deployment 对象，并把结果填充到 dep 里
	//err = r.Client.Get(ctx, req.NamespacedName, dep)
	//
	//if err != nil {
	//	if errors.IsNotFound(err) {
	//		//如果没有这个deployment那就创建这个deployment
	//		//定义deployment的相关信息 dep的名称和名称空间与bookapp一致
	//		klog.Infof("创建deployment")
	//
	//		//定义完成后，创建deployment
	//		err = r.Client.Create(ctx, dep)
	//		//如果创建失败，打印日志并重新入队列
	//		if err != nil {
	//			klog.Errorf("Create deployment fail:%v-%v-\n", err, req.NamespacedName)
	//			return ctrl.Result{}, err
	//		}
	//		return ctrl.Result{}, nil
	//		//容器
	//		//klog.Infof("bookapp:%v-%v is delete\n", req.NamespacedName, req.Name)
	//	}
	//	//如果是其他错误，你就重新入队列
	//	return ctrl.Result{}, err
	//}
	////如果找到了相应的deployment，那就更新deployment
	//if !reflect.DeepEqual(newdep.(*appsv1.Deployment), dep.Spec) {
	//	klog.Infof("add or update bookapp:%v-%v\n", req.NamespacedName, req.Name)
	//	dep.Spec = newdep.(*appsv1.Deployment).Spec
	//	err = r.Client.Update(ctx, dep)
	//	if err != nil {
	//		klog.Errorf("更新失败")
	//	}
	//	return ctrl.Result{}, err
	//
	//}
	//构建deployment
	//定义一个空的deployment
	dep := &appsv1.Deployment{}
	dep.Name = bookapp.Name
	dep.Namespace = bookapp.Namespace
	//内置的创建更新方法
	_, err = ctrl.CreateOrUpdate(ctx, r.Client, dep, func() error {
		klog.Infof("add or update bookapp-deployment:%v-%v\n", req.NamespacedName, req.Name)
		dep.Labels = bookapp.Labels
		//副本数量
		dep.Spec.Replicas = bookapp.Spec.Size
		//selector
		selector := &metav1.LabelSelector{
			MatchLabels: bookapp.Labels,
		}
		dep.Spec.Selector = selector
		dep.Spec.Template.Labels = bookapp.Labels
		//创建的pod的信息
		contaners := []corev1.Container{
			{
				Name:            bookapp.Name,
				Image:           bookapp.Spec.Image,
				ImagePullPolicy: bookapp.Spec.ImagePullPolicy,
			},
		}
		dep.Spec.Template.Spec.Containers = contaners
		//配置deployment的拥有者，实现，删除bookapp时，可以同时删除相应的资源,例如deployment
		err = ctrl.SetControllerReference(&bookapp, dep, r.Scheme)
		return err
	})
	//return ctrl.Result{}, err
	if err != nil {
		return ctrl.Result{}, err
	}

	//构建service
	svc := &corev1.Service{}
	svc.Name = bookapp.Name
	svc.Namespace = bookapp.Namespace
	_, err = ctrl.CreateOrUpdate(ctx, r.Client, svc, func() error {
		klog.Infof("add or update bookapp-service:%v-%v\n", req.NamespacedName, req.Name)
		svc.Labels = bookapp.Labels
		svc.Spec.Selector = bookapp.Labels
		ports := []corev1.ServicePort{
			{
				Name:       bookapp.Spec.Port.Name,
				Port:       bookapp.Spec.Port.Port,
				TargetPort: bookapp.Spec.Port.TargetPort,
			},
		}
		svc.Spec.Ports = ports
		err = ctrl.SetControllerReference(&bookapp, svc, r.Scheme)
		return err
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	//构建ingress
	ing := &networkingv1.Ingress{}
	ing.Name = bookapp.Name
	ing.Namespace = bookapp.Namespace
	ing.Labels = bookapp.Labels
	_, err = ctrl.CreateOrUpdate(ctx, r.Client, ing, func() error {
		klog.Infof("add or update bookapp-ingress:%v-%v\n", req.NamespacedName, req.Name)
		pathType := networkingv1.PathTypePrefix
		ing.Spec.IngressClassName = pointer.String("nginx")
		ing.Spec.Rules = []networkingv1.IngressRule{
			{
				Host: bookapp.Spec.ServerName, // CRD 中的 serverName
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: bookapp.Name, // 对应 Service 名称
										Port: networkingv1.ServiceBackendPort{
											Number: bookapp.Spec.Port.Port,
										},
									},
								},
							},
						},
					},
				},
			},
		}
		err = ctrl.SetControllerReference(&bookapp, ing, r.Scheme)
		return err
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	//状态更新
	if err := r.updateBookAppStatus(ctx, &bookapp); err != nil {
		klog.Error(err, "更新bookapp状态失败：%v")
		return ctrl.Result{}, err
	}

	klog.Infof("Reconcile success")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookappReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Bookapp{}).
		//Owns(&appsv1.Deployment{}) 作用删除deployment后自动重建
		Named("bookapp").Owns(&appsv1.Deployment{}).Owns(&corev1.Service{}).Owns(&networkingv1.Ingress{}).
		Complete(r)
}
