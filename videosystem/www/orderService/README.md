订单服务：
1. 日常订单服务
2. 秒杀系统设计
    1.1  如何处理商品库存的实时更新和校验，以避免超卖
    1.2  消息的持久化、消息确认机制以及消费者的并发消费配置
    1.3  关于用户认证和限流，是否已经有一套策略，尤其是在高并发的秒杀场景下

秒杀服务
    1. 针对目前业务需求中商品相关储存仅限于单库单表单字段，设计考虑到并发效果，这里利用redis原子性操作保证并发安全的情况下，设计缓存机制，并且需要考虑缓存击穿的情况
    2. 业务使用rabbitMQ进行消息削峰操作，并结合websocket进行异步通讯
    3. 认证服务中需要做到请求限流，请求过滤，订单服务只接收符合条件的请求并处理
    4. 系统完成部署后，需要进行压力测试































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