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

/*
	设置缓冲池大小

	package main

import (
    "github.com/gorilla/websocket"
    "sync"
)

// 定义一个实现了 websocket.BufferPool 接口的结构体
type BufferPool struct {
    pool sync.Pool
}

// NewBuffer 方法用于获取新的缓冲区
func (p *BufferPool) Get() []byte {
    return p.pool.Get().([]byte)
}

// Put 方法用于回收使用完的缓冲区
func (p *BufferPool) Put(b []byte) {
    p.pool.Put(b)
}

func main() {
    // 初始化缓冲池
    bufferPool := &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                // 指定缓冲区的大小，这里假设为1024字节
                return make([]byte, 1024)
            },
        },
    }

    // 在 Upgrader 中使用自定义的缓冲池
    var upgrader = websocket.Upgrader{
        WriteBufferPool: bufferPool,
        // 其他字段配置...
    }

    // 使用 upgrader 进行 WebSocket 握手和后续处理...
}

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
			delete(services.WebSocketclients, userID)
			services.ClientsLock.Unlock()
			return
		default:
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				// 处理读取错误
				global.SendLogs("error", fmt.Sprintf("读取消息报错：用户ID:%s", userID), err)
				return
			}
			_ = messageType
			_ = p
		}
	}
}
