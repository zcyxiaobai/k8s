//单元测试

package jwt

import (
	"testing"
	"time"
)

// 测试生成token
func TestJwtToken_GenerateToken(t *testing.T) {
	jwttoken, err := NewJwtToken("F:\\北大青鸟培训\\golang\\go_pod\\keys\\privatekey.pem", "")
	if err != nil {
		t.Fatal(err)
	}

	nowFunc = func() time.Time {
		return time.Unix(1756829884, 0)
	}
	//这里的一个测试目的是，测试代码生成的token是否与官网中生成的一致
	cases := struct {
		username string
		token    string
	}{
		username: "admin",
		//官网中生成的token
		token: ``,
	}
	token, err := jwttoken.GenerateToken(cases.username)
	if err != nil {
		t.Fatal(err)
	}
	//比较两个token是否一致
	if token != cases.token {
		t.Errorf("gost:%s\nwant:%s", token, cases.token)
	}
}

// 测试校验token
func TestJwtToken_ParseToken(t *testing.T) {
	jwttoken, err := NewJwtToken("F:\\北大青鸟培训\\golang\\go_pod\\keys\\privatekey.pem", "F:\\北大青鸟培训\\golang\\go_pod\\keys\\publickey.pem")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name     string
		token    string
		wantErr  bool
		wantname string
	}{
		{
			name: "valid token",
			//合法的token
			token:    ``,
			wantErr:  false,
			wantname: "admin",
		},
		{
			name: "expire token",
			//过期的token
			token:   ``,
			wantErr: true,
		},
		{
			name: "fake token",
			//伪造的token
			token:   ``,
			wantErr: true,
		},
		{
			//无效的token
			name:    "invalid token",
			token:   ``,
			wantErr: true,
		},
	}
	for _, tt := range cases {
		claim, err := jwttoken.ParseToken(tt.token)
		if err == nil && tt.wantErr {
			t.Errorf("Test %s want err but get no", tt.name)
			continue
		}
		if err != nil {
			if !tt.wantErr {
				t.Errorf("Test %s want no error but got err %v", tt.name, err)
			}
			continue
		}
		if username, ok := (*claim)["username"]; ok {
			if username, ok = username.(string); ok {
				if username != tt.wantname {
					t.Errorf("Test %s want %s but got %v", tt.name, tt.wantname, username)
					continue
				}
			} else {
				t.Fatalf("%s:username error", tt.name)
			}

		} else {
			t.Fatalf("%s:no username", tt.name)
		}

	}

}
