package bootstrap

import (
	"auth/global"
	"context"

	"github.com/go-redis/redis/v8"
)

/*
初始化 缓存
*/
func InitializeRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.App.Config.Redis.Host + ":" + global.App.Config.Redis.Port,
		Password: "",                         // 密码，没有则留空
		DB:       global.App.Config.Redis.DB, // 使用默认DB
	})
	// 测试redis
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		global.SendLogs("error", "redis 初始化错误", err)
	} else {
		global.SendLogs("info", "redis 初始化成功")
		global.App.Redis = rdb
	}
}
