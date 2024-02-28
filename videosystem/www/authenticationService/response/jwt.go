package response

import (
	"auth/global"
	"time"

	"github.com/dgrijalva/jwt-go"
)

/*
创建jwt
*/

type jwtService struct {
}

var JwtService = new(jwtService)

// 自定义jwt的claims，考虑到扩展性，这里只多存放用户角色信息
type CustomClaims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

// 生成Jwt Token 函数  GuardName 区分客户端
func (JwtService *jwtService) CreateJwtToken(GuardName string, userId string, roles []string) (tokenOut TokenOutPut, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + global.App.Config.Jwt.Jwttil,
				Id:        userId,
				Issuer:    GuardName,
			},
			Roles: roles,
		},
	)
	tokenstr, err := token.SignedString([]byte(global.App.Config.Jwt.Secretkey))
	tokenOut = TokenOutPut{
		tokenstr,
		int(global.App.Config.Jwt.Jwttil),
		TokenType,
	}
	return
}
