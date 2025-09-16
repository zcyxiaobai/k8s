package jwt

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
	//"io/ioutil"
	"os"
	"time"
)

// JwtToken 是一个结构体，里面保存了一个 RSA 私钥对象
type JwtToken struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// 返回当前时间
var nowFunc = func() time.Time {
	return time.Now()
}

// 这是一个工厂函数，用来 创建 JwtToken 对象
func NewJwtToken(key1, key2 string) (*JwtToken, error) {
	privatekey, err := loadPrivateKey(key1)
	if err != nil {
		return nil, err
	}
	publickey, err := loadPublicKey(key2)
	if err != nil {
		return nil, err
	}
	//返回公钥和私钥
	return &JwtToken{privateKey: privatekey, publicKey: publickey}, nil
}

// 生成token
func (j *JwtToken) GenerateToken(username string) (string, error) {
	now := nowFunc()
	//now := time.Now()
	//过期时间两个小时
	expirationTime := now.Add(2 * time.Hour)
	claims := jwt.MapClaims{
		"sub":      "user_token",
		"username": username,
		"iat":      now.Unix(),
		"exp":      expirationTime.Unix(),
	}

	//使用RS512
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	//使用私钥签名
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// 从文件里读取私钥 PEM 内容
// 调用 jwt.ParseRSAPrivateKeyFromPEM 解析成 Go 的 *rsa.PrivateKey 对象
// 返回这个私钥
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyDate, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyDate)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
