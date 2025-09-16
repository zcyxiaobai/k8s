// 用户登陆模块
package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"y2505.com/bookapp/dao"
	"y2505.com/bookapp/jwt"
)

//定义用户信息的结构体

type UserDate struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserView struct {
	*jwt.JwtToken
}

func (u *UserView) Login(c *gin.Context) {
	userdata := UserDate{}
	err := c.ShouldBind(&userdata)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "登陆失败",
		})
		klog.Error(err.Error())
		return
	}
	user, err := dao.QueryUser(userdata.Username, userdata.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "登陆失败",
		})
		return
	}
	//生成token
	token, err := u.JwtToken.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "登陆失败",
		})
		klog.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "登陆成功",
		"username": user.Username,
		//登陆成功后将token返回给前端
		"token": token,
	})

}
