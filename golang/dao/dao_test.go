// 单元测试
package dao

import (
	"testing"
)

// 测试创建表
func TestCreateDB(t *testing.T) {
	CreateDB()
}

// 测试插入数据
func TestInsertUserData(t *testing.T) {
	//cases := []UserModdel{
	//	UserModdel{
	//		Username: "user1",
	//		Password: "123456",
	//	},
	//	UserModdel{
	//		Username: "user2",
	//		Password: "123456",
	//	},
	//}
	//for _, c := range cases {
	//	err := InsertUserData(c)
	//	if err != nil {
	//		fmt.Errorf("插入数据失败")
	//		return
	//	}
	//	fmt.Printf("插入数据成功")
	//}

	cases := []struct {
		Name    string
		User    UserModdel
		WantErr bool
	}{
		{
			Name: "Normal Create User",
			User: UserModdel{
				Username: "admin",
				Password: "123456",
			},
			WantErr: false,
		},
		//唯一键冲突
		{
			Name: "Duplicate Create User",
			User: UserModdel{
				Username: "admin",
				Password: "123456",
			},
			WantErr: true,
		},
	}
	for _, c := range cases {
		err := InsertUserData(c.User)
		//想出错，没有错
		if c.WantErr && err == nil {
			t.Errorf("TestCreateDB(%s): expected error but got nil", c.Name)
			continue
		}
		//不想出错，出了错
		if err != nil && !c.WantErr {
			t.Errorf("TestCreateDB(%s): expected no error but got %v", c.Name, err)
		}
	}
}

// 测试用户登陆
func TestQueryUser(t *testing.T) {
	cases := []struct {
		Name         string
		Username     string
		Password     string
		WantErr      bool
		wantUserName string
	}{
		{
			Name:         "Normal Query User",
			Username:     "user1",
			Password:     "123456",
			WantErr:      false,
			wantUserName: "user1",
		},
		{
			Name:     "Not User Exist",
			Username: "admin1",
			Password: "123456",
			WantErr:  true,
		}, {
			Name:     "Wrong Password",
			Username: "user1",
			Password: "1234567869",
			WantErr:  true,
		},
	}
	for _, c := range cases {
		newuser, err := QueryUser(c.Username, c.Password)
		if c.WantErr && err == nil {
			t.Errorf("TestQueryUser(%s): expected error but got nil", c.Name)
			continue
		}
		if err != nil {
			if !c.WantErr {
				t.Errorf("TestQueryUser(%s): expected no error but got %v", c.Name, err)
			}
			continue
		}
		//if newuser == nil {
		//	t.Errorf("TestQueryUser(%s): expected not nil", c.Name)
		//	continue
		//}
		if newuser.Username != c.wantUserName {
			t.Errorf("%s expected: %s, got: %s", c.Name, c.wantUserName, newuser.Username)
		}
	}
}
