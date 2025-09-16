package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	ctrl "sigs.k8s.io/controller-runtime"
	webappv1 "y2505.com/bookapp/api/v1"
	"y2505.com/bookapp/bookhttp"
	y2505jwt "y2505.com/bookapp/jwt"
	"y2505.com/bookapp/k8s"
	"y2505.com/bookapp/views"
)

// 注册
var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(webappv1.AddToScheme(scheme))
}

// -----------------以上为另一种接口管理方式----------------------
// 使用token的中间件  AuthMiddleware中间件
func AuthMiddleware(jwtToken *y2505jwt.JwtToken) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		//从请求头获取 Authorization 意思是前台在请求过程中必须有key为Authorization
		//可以是：Authorization   Bearer +  合法的token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid Authorization header format",
				})
				c.Abort()
				return

			}
		} else {
			// 如果请求头没有，从 query 参数获取 token (用于 WebSocket)
			tokenStr = c.Query("token")
			if tokenStr == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "missing token",
				})
				c.Abort()
				return
			}
		}
		//去掉Bearer前缀

		//解析token
		claims, err := jwtToken.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}
		//把username 保存到gin.Context,后续handler可以使用
		if username, ok := (*claims)["username"].(string); ok {
			if username != "admin" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "user is not admin",
				})
			}
			c.Set("username", username)
		}
		//放行
		c.Next()
	}
}

//func AuthMiddleware(jwtToken *y2505jwt.JwtToken) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//从请求头获取 Authorization 意思是前台在请求过程中必须有key为Authorization
//		//可以是：Authorization   Bearer +  合法的token
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"error": "missing authorization header",
//			})
//			c.Abort()
//			return
//
//		}
//		//去掉Bearer前缀
//		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
//		if tokenStr == authHeader {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"error": "invalid Authorization header format",
//			})
//			c.Abort()
//			return
//
//		}
//		//解析token
//		claims, err := jwtToken.ParseToken(tokenStr)
//		if err != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"error": "invalid or expired token",
//			})
//			c.Abort()
//			return
//		}
//		//把username 保存到gin.Context,后续handler可以使用
//		if username, ok := (*claims)["username"].(string); ok {
//			if username != "admin" {
//				c.JSON(http.StatusUnauthorized, gin.H{
//					"error": "user is not admin",
//				})
//			}
//			c.Set("username", username)
//		}
//		//放行
//		c.Next()
//	}
//}

//// WebSocket的中间件
//func WsAuthMiddleware(jwtToken *y2505jwt.JwtToken) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		tokenStr := c.Query("token") // 从 query 获取 token
//		if tokenStr == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
//			c.Abort()
//			return
//		}
//
//		// 去掉可能的 Bearer 前缀
//		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
//
//		claims, err := jwtToken.ParseToken(tokenStr)
//		if err != nil {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
//			c.Abort()
//			return
//		}
//
//		if username, ok := (*claims)["username"].(string); ok {
//			c.Set("username", username)
//		}
//
//		c.Next()
//	}
//}

func main() {
	//注册
	//scheme := runtime.NewScheme()
	//utilruntime.Must(webappv1.AddToScheme(scheme))
	//cli, err := client.New(ctrl.GetConfigOrDie(), client.Options{
	//	Scheme: scheme,
	//})
	//if err != nil {
	//	klog.Fatal(err)
	//	return
	//}
	//
	//bh := bookhttp.NewBookHttp(cli, gin.Default())
	//rg := bh.Router.Group("/bookapp")
	//{
	//	rg.GET("/:namespace/:appname", bh.GetBookApp)
	//	rg.POST("/", bh.CreateBookApp)
	//	rg.GET("/bookapps", bh.ListBookApps)
	//	rg.GET("/bookapps/:namespace", bh.ListBookAppsByNamespace)
	//}
	//bh.Router.Run(":8080")
	//---------------------以上代码可以正常使用---------------------------

	//manager托管API
	//创建并启动一个 控制器运行环境，统一管理 client、缓存、控制器和 webhook 等组件
	router := gin.Default()
	// 允许所有来源请求,防止前端访问出现跨域问题
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 获取 Kubernetes 配置
	cfg := ctrl.GetConfigOrDie() // *rest.Config

	// 直接用 Clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	metricsClient, err := metricsv.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	// 原有 PodClient 不用改
	pv := views.PodView{
		PodClient: k8s.NewPodClient(clientset, metricsClient),
	}

	//初始化token
	klog.Info("初始化token")
	jt, err := y2505jwt.NewJwtToken(
		"/root/projects/bookapp/keys/privatekey.pem",
		"/root/projects/bookapp/keys/publickey.pem")
	if err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
	//登陆的路由
	uv := views.UserView{}
	uv.JwtToken = jt
	//登陆接口
	router.POST("/login", uv.Login)
	//获取所有总数
	router.GET("/listall", pv.GetResourceCount)
	//获取node的信息
	router.GET("/nodemetrics", pv.GetNodeMetrics)
	//获取每个节点上正常和非正常的pod
	router.GET("/getnodepodstatus", pv.GetNodePodStatus)
	//获取所有的namespace
	router.GET("/getnamespacelist", pv.GetNamespaceList)
	//获取元素总数接口

	// 在这里使用 `AuthMiddleware` 进行路由保护
	//router.Use(AuthMiddleware(jt)) // 保护整个 api

	//pod的相关接口
	{
		pod := router.Group("/pod")
		//查询所有的pod
		//查询pod
		pod.GET("/podalllist", pv.GetPodAllList)
		//pod日志
		pod.GET("/log/:namespace/:podname/:cname", AuthMiddleware(jt), pv.GetLog)
		//pod  webshell
		pod.GET("/shell/:namespace/:podname/:cname", AuthMiddleware(jt), pv.GetShell(cfg))
		//删除pod
		pod.POST("/Delete", AuthMiddleware(jt), pv.Delete)
		//创建新增pod
		pod.POST("/Create", AuthMiddleware(jt), pv.CreatePod)
		//修改pod
		pod.POST("/Update", AuthMiddleware(jt), pv.UpdatePod)

	}
	//deployment的相关接口
	{
		dep := router.Group("/dep")
		//查看deployment
		dep.GET("/GetDeploymentList", pv.GetDeploymentList)
		//删除deployment
		dep.POST("/DeleteDev", AuthMiddleware(jt), pv.DeploymentDelete)
		//创建deployment
		dep.POST("/CreateDev", AuthMiddleware(jt), pv.CreateDeployment)
		//修改deployment
		dep.POST("/UpdateDev", AuthMiddleware(jt), pv.UpdateDeployment)
	}
	//service的相关接口
	{
		svc := router.Group("/service")
		//查看service
		svc.GET("/GetServiceInfo", pv.GetServiceInfo)
		//删除service
		svc.POST("/DeleteService", AuthMiddleware(jt), pv.DeleteService)
		//创建service
		svc.POST("/CreateService", AuthMiddleware(jt), pv.CreateService)
		//修改service
		svc.POST("/UpdateService", AuthMiddleware(jt), pv.UpdateService)
	}
	//ingress的相关接口
	{
		ing := router.Group("/ingest")
		//获取ingress
		ing.GET("/GetIngressAll", pv.GetIngressAll)
		//删除ingress
		ing.POST("/IngressDelete", AuthMiddleware(jt), pv.IngressDelete)
		//创建ingress
		ing.POST("/CreateIngress", AuthMiddleware(jt), pv.CreateIngress)
		//修改ingress
		ing.POST("/UpdateIngress", AuthMiddleware(jt), pv.UpdateIngress)
	}

	//创建一个 controller-manager 的实例，ctrl.NewManager 是 Kubernetes 控制器（controller）框架的一部分，它会创建一个用于管理控制器、API 客户端、缓存等组件的管理器
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		klog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	//创建BookHttp示例
	klog.Info("创建BookHttp示例")
	bh := bookhttp.NewBookHttp(mgr.GetClient(), router, jt)
	//使用中间件保护bookhttp路由
	bookappRouter := bh.Router.Group("/bookapp")
	//使用token验证
	//bookappRouter.Use(AuthMiddleware(jt))
	// 将路由的具体处理交给 BookHttp 类
	bookappRouter.GET("/:namespace/:appname", bh.GetBookApp)
	//创建bookapp
	bookappRouter.POST("/", AuthMiddleware(jt), bh.CreateBookApp)
	//查看所有的bookapp
	bookappRouter.GET("/bookapps", bh.ListBookApps)
	bookappRouter.GET("/bookapps/:namespace", bh.ListBookAppsByNamespace)
	//删除指定的bookapp
	bookappRouter.POST("/deletebookapp", AuthMiddleware(jt), bh.DeletebookApp)
	//修改指定的bookapp
	bookappRouter.POST("/UpdateBookApp", AuthMiddleware(jt), bh.UpdateBookApp)
	//将 bh（即 BookHttp 类型的实例）添加到 Kubernetes 控制器管理器（manager）
	err = mgr.Add(bh)
	if err != nil {
		klog.Error(err, "unable to register controller")
		os.Exit(1)
	}

	// 监听并在 0.0.0.0:8080 上启动服务
	// 启动 manager 在 Goroutine 中，以便不阻塞主程序

	go func() {
		klog.Info("starting manager")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			klog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}()

	//启动Gin
	//if err := router.Run(":8090"); err != nil {
	//	klog.Error(err, "problem running manager")
	//	os.Exit(1)
	//}
	router.Run(":8090")
}
