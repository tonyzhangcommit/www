package management

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"userservice/app/response"
	"userservice/app/services"
	"userservice/global"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/*
	websocker 接口视图函数
*/

// 自定义jwt的claims，考虑到扩展性，这里只多存放用户角色信息
type CustomClaims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

// 初始化 Upgrader 结构体， 用于升级http协议至websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		query := r.URL.Query()
		tokenStr := query.Get("token")
		if tokenStr == "" {
			return false
		} else {
			// 这里验证jwt token 合法性代码
			token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(global.App.Config.WebSocket.JwtSecretkey), nil
			})
			if err != nil {
				return false
			}
			claims := token.Claims.(*CustomClaims)
			return claims.Issuer == "app"
		}
	},
	HandshakeTimeout:  time.Duration(global.App.Config.WebSocket.HandShakeTimeout) * time.Second,
	EnableCompression: global.App.Config.WebSocket.EnableCompression,
}

// 视图函数
func WebsocketHandler(c *gin.Context) {
	userID := c.Query("userid")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.BusinessFail(c, "订单查询失败，稍后请在”我的-订单“查询订单")
		return
	}
	defer ws.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	services.WebSocketclients[userID] = ws
	// 保持连接
	for {
		select {
		case <-ctx.Done():
			// 上下文超时或被取消，关闭WebSocket连接
			services.ClientsLock.Lock()
			services.WebSocketclients[userID].Close()
			delete(services.WebSocketclients, userID)
			services.ClientsLock.Unlock()
			return
		default:
			_, p, err := ws.ReadMessage()
			if err != nil {
				// 处理读取错误
				global.SendLogs("error", fmt.Sprintf("读取消息报错：用户ID:%s", userID), err)
				return
			}
			if string(p) == "close" {
				fmt.Println("用户主动关闭连接")
				services.ClientsLock.Lock()
				services.WebSocketclients[userID].Close()
				delete(services.WebSocketclients, userID)
				services.ClientsLock.Unlock()
				fmt.Println("关闭连接成功")
				return
			}
		}
	}
}
