package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库的操作使用的是gorm
type UserModdel struct {
	gorm.Model
	//gorm:"unique"表示唯一键
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password"`
}

//var cases = []UserModdel{
//	UserModdel{
//		Username: "user1",
//		Password: "123456",
//	},
//	UserModdel{
//		Username: "user2",
//		Password: "123456",
//	},
//}

// 定义数据库连接
var dsn = "root:Aa@123456@tcp(192.168.216.95:3306)/y2505?charset=utf8mb4&parseTime=True&loc=Local"

//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

func CreateDB() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//操作是数据库，根据结构体创建表
	if err != nil {
		panic(err)
	}
	//判断是否存在这张表，没有才创建
	if !db.Migrator().HasTable(&UserModdel{}) {
		err = db.Migrator().CreateTable(&UserModdel{})
		if err != nil {
			panic(err)
		}
	}
}

//往表里面插入数据

func InsertUserData(user UserModdel) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	res := db.Create(&user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("创建用户失败")
	}
	return nil
}

// 实现根据用户名和密码登陆的简单功能
func QueryUser(username, password string) (*UserModdel, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	user := UserModdel{}
	if err != nil {
		return nil, err
	}
	tx := db.Where("username = ? AND password = ?", username, password).Find(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}
