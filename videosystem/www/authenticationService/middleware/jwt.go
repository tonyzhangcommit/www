package middleware

import (
	"auth/global"
	"auth/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

/*
jwt 鉴权中间件
*/

func JWTAUTH(GuardName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取jwt token
		tokenStr := ctx.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.JwtTokenErrorFail(ctx, "登录已过期，请重新登录")
			ctx.Abort()
			return
		}
		// 这里获取token 的字符串为bearer **************
		tokenStr = tokenStr[len(response.TokenType)+1:]
		token, err := jwt.ParseWithClaims(tokenStr, &response.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(global.App.Config.Jwt.Secretkey), nil
		})
		if err != nil {
			response.JwtTokenErrorFail(ctx, "登录已过期，请重新登录")
			ctx.Abort()
			return
		}
		claims := token.Claims.(*response.CustomClaims)
		if claims.Issuer != GuardName {
			response.JwtTokenErrorFail(ctx, "登录已过期，请重新登录")
			ctx.Abort()
			return
		}
		ctx.Set("token", token)
		ctx.Set("userid", claims.Id)
	}
}
