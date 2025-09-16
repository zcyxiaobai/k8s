package bookhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	webappv1 "y2505.com/bookapp/api/v1"
	"y2505.com/bookapp/jwt"
)

type BookHttp struct {
	client.Client
	//用来保存 Gin 的路由引擎实例
	Router *gin.Engine
	//Token
	JwtToken *jwt.JwtToken
}

// 定义一个结构体，用于接收前端传过来的创建bookapp的数据
type BookAppRequest struct {
	Name      string               `json:"name"`
	Namespace string               `json:"namespace"`
	Labels    map[string]string    `json:"labels"`
	Spec      webappv1.BookappSpec `json:"spec"`
}

// 定义删除bookapp时使用的结构体
// 定义请求结构
type BookAppDelete struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// 定义一个修改的结构体
type BookAppUpdate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Labels    map[string]string      `json:"labels"`
	Spec      map[string]interface{} `json:"spec"`
}

func NewBookHttp(c client.Client, r *gin.Engine, jwtToken *jwt.JwtToken) *BookHttp {
	bookhttp := &BookHttp{}
	bookhttp.Client = c
	bookhttp.Router = r
	bookhttp.JwtToken = jwtToken
	return bookhttp
}

// 注册路由
// 查找整个集群所有的bookapp

//	func (b *BookHttp) ListBookApps(c *gin.Context) {
//		ctx := context.Background()
//
//		bookappList := &webappv1.BookappList{}
//		// 不传 namespace 参数，直接查询全局
//		err := b.List(ctx, bookappList)
//		if err != nil {
//			klog.Errorf("list bookapps err:%v", err)
//			c.JSON(http.StatusInternalServerError, gin.H{
//				"message": "查询失败",
//			})
//			return
//		}
//
//		// 组装返回数据
//		resp := make([]map[string]interface{}, 0, len(bookappList.Items))
//		for _, ba := range bookappList.Items {
//			age := time.Since(ba.CreationTimestamp.Time).Round(time.Second)
//			resp = append(resp, map[string]interface{}{
//				"name":      ba.Name,
//				"namespace": ba.Namespace,
//				"status":    ba.Status,
//				"image":     ba.Spec.Image,
//				"size":      ba.Spec.Size,
//				"age":       age.String(),
//			})
//		}
//
//		c.JSON(http.StatusOK, gin.H{
//			"message": "success",
//			"data":    resp,
//		})
//	}
//
// 查找整个集群所有的bookapp
func (b *BookHttp) ListBookApps(c *gin.Context) {
	ctx := context.Background()

	bookappList := &webappv1.BookappList{}
	// 不传 namespace 参数，直接查询全局
	err := b.List(ctx, bookappList)
	if err != nil {
		klog.Errorf("list bookapps err:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询失败",
		})
		return
	}

	// 组装返回数据
	resp := make([]map[string]interface{}, 0, len(bookappList.Items))
	for _, ba := range bookappList.Items {
		age := time.Since(ba.CreationTimestamp.Time).Round(time.Second)

		// 这里直接把整个 Spec 和 Status 放进去
		resp = append(resp, map[string]interface{}{
			"name":        ba.Name,
			"namespace":   ba.Namespace,
			"labels":      ba.Labels,
			"annotations": ba.Annotations,
			"spec":        ba.Spec,   // 直接返回完整的 spec
			"status":      ba.Status, // 直接返回完整的 status
			"age":         age.String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resp,
	})
}

// 查找指定名称空间下的bookapp
// 查询指定 namespace 下所有的 bookapp
func (b *BookHttp) ListBookAppsByNamespace(c *gin.Context) {
	ctx := context.Background()
	namespace := c.Param("namespace") // 从路由参数获取命名空间，比如 /bookapps/:namespace

	bookappList := &webappv1.BookappList{}
	// 只列出指定 namespace 下的资源
	err := b.List(ctx, bookappList, &client.ListOptions{Namespace: namespace})
	if err != nil {
		klog.Errorf("list bookapps in namespace %s err: %v", namespace, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询失败",
		})
		return
	}

	// 格式化返回数据
	var result []gin.H
	for _, item := range bookappList.Items {
		age := time.Since(item.CreationTimestamp.Time).Round(time.Second)
		result = append(result, gin.H{
			"name":      item.Name,
			"namespace": item.Namespace,
			"status":    item.Status,
			"image":     item.Spec.Image,
			"size":      item.Spec.Size,
			"age":       age.String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result,
	})
}

// 查找指定bookapp的接口
func (b *BookHttp) GetBookApp(c *gin.Context) {
	ctx := context.Background()
	appname := c.Param("appname")
	namespace := c.Param("namespace")
	bookapp := webappv1.Bookapp{}
	err := b.Get(ctx, types.NamespacedName{Namespace: namespace, Name: appname}, &bookapp)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Bookapp Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "内部错误",
		})
		klog.Errorf("get bookapp err:%v", err)
		return
	}
	age := time.Since(bookapp.CreationTimestamp.Time).Round(time.Second)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		//"data":    bookapp,
		"bookappname": appname,
		"namespace":   namespace,
		"status":      bookapp.Status,
		"image":       bookapp.Spec.Image,
		"size":        bookapp.Spec.Size,
		"age":         age.String(),
	})
}

// 创建bookapp的接口
func (b *BookHttp) CreateBookApp(c *gin.Context) {
	ctx := context.Background()
	baq := &BookAppRequest{}
	err := c.ShouldBind(baq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "序列化失败",
		})
		return
	}
	bookapp := webappv1.Bookapp{}
	bookapp.Name = baq.Name
	bookapp.Namespace = baq.Namespace
	bookapp.Labels = baq.Labels
	bookapp.Spec = baq.Spec
	err = b.Create(ctx, &bookapp)
	if err != nil {
		klog.Errorf("创建失败:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
	})
}

// 删除指定BookApp的方法
func (b *BookHttp) DeletebookApp(c *gin.Context) {
	ctx := context.Background()
	bookapp := BookAppDelete{}
	err := c.ShouldBind(&bookapp)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "反序列化失败",
		})
		return
	}
	bookappnew := &webappv1.Bookapp{}
	bookappnew.Name = bookapp.Name
	bookappnew.Namespace = bookapp.Namespace
	if bookapp.Namespace == "" || bookapp.Name == "" {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "namespacce 或 name不能为空",
		})
		return
	}
	//删除方法
	err = b.Delete(ctx, bookappnew)
	if err != nil {
		klog.Errorf("删除 bookapp %s/%s 失败: %v", bookapp.Namespace, bookapp.Name, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})

}

// 修改指定的bookapp
func (b *BookHttp) UpdateBookApp(c *gin.Context) {
	ctx := context.Background()
	bookapp := BookAppUpdate{}
	err := c.ShouldBind(&bookapp)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "反序列化失败",
		})
		return
	}
	if bookapp.Namespace == "" || bookapp.Name == "" {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Name 和 Namespace 不能为空",
		})
		return
	}
	bookappnew := &webappv1.Bookapp{}

	// 1️⃣ 先获取原对象，拿到 resourceVersion
	err = b.Get(ctx, types.NamespacedName{Namespace: bookapp.Namespace, Name: bookapp.Name}, bookappnew)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "内部错误",
		})
		return
	}
	// 2️⃣ 更新 Labels
	if bookapp.Labels != nil {
		bookappnew.Labels = bookapp.Labels
	}
	// 3️⃣ 更新 Spec（这里根据你 Bookapp 的结构，可自行映射）
	if specSize, ok := bookapp.Spec["size"].(float64); ok {
		sizeVal := int32(specSize)
		bookappnew.Spec.Size = &sizeVal
		//bookappnew.Spec.Size = int32(specSize)
	}
	if image, ok := bookapp.Spec["image"].(string); ok {
		bookappnew.Spec.Image = image
	}
	if imagePullPolicy, ok := bookapp.Spec["imagePullPolicy"].(string); ok {
		bookappnew.Spec.ImagePullPolicy = corev1.PullPolicy(imagePullPolicy)
	}
	if portMap, ok := bookapp.Spec["port"].(map[string]interface{}); ok {
		if name, ok := portMap["name"].(string); ok {
			bookappnew.Spec.Port.Name = name
		}
		if protocol, ok := portMap["protocol"].(string); ok {
			bookappnew.Spec.Port.Protocol = corev1.Protocol(protocol)
		}
		if portValue, ok := portMap["port"].(float64); ok {
			bookappnew.Spec.Port.Port = int32(portValue)
		}
		if targetPort, ok := portMap["targetPort"].(float64); ok {
			bookappnew.Spec.Port.TargetPort = intstr.FromInt(int(targetPort))
		}
	}
	if serverName, ok := bookapp.Spec["serverName"].(string); ok {
		bookappnew.Spec.ServerName = serverName
	}
	// 4️⃣ 调用 Update，resourceVersion 已在 existing 中
	err = b.Update(ctx, bookappnew)
	if err != nil {
		klog.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "修改bookapp失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "修改bookapp成功",
	})
}

// manager托管api,start方法
func (b *BookHttp) Start(ctx context.Context) error {
	//routergroup := b.Router.Group("/bookapp")
	//
	//{
	//	routergroup.GET("/:namespace/:appname", b.GetBookApp)
	//	routergroup.POST("/", b.CreateBookApp)
	//	routergroup.GET("/bookapps", b.ListBookApps)
	//	routergroup.GET("/bookapps/:namespace", b.ListBookAppsByNamespace)
	//}
	//return b.Router.Run(":8090")
	return nil
}
