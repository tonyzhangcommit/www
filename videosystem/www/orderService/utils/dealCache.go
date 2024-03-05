package utils

import (
	"context"
	"errors"
	"time"
	"userservice/global"

	"github.com/go-redis/redis/v8"
)

/*
	缓存操作封装
*/

// 验证码存入缓存
func SetVirifCode(key, value string, dur time.Duration) error {
	if global.App.Redis == nil {
		global.SendLogs("error", "初始化redis服务失败") // 这里应该加上预警提示
		// 尝试使用本地缓存
		global.Store.Set(key, value, dur)
	} else {
		err := global.App.Redis.Set(context.Background(), key, value, dur).Err()
		if err != nil {
			global.SendLogs("error", "保存验证码到redis报错", err)
			// 尝试使用本地缓存
			global.Store.Set(key, value, dur)
		}
	}

	return nil
}

// 验证码从缓存读取
func GetVirifCode(key string) string {
	if global.App.Redis == nil {
		global.SendLogs("error", "初始化redis服务失败") // 这里应该加上预警提示
		// 尝试使用本地缓存
		if valInterface, okv := global.Store.Get(key); okv {
			valstr, ok := valInterface.(string)
			if ok {
				return valstr
			} else {
				err := errors.New("验证码类型不是string")
				global.SendLogs("error", "从本地缓存获取验证码报错", err)
				return ""
			}
		}
	} else {
		val, err := global.App.Redis.Get(context.Background(), key).Result()
		if err == redis.Nil {
			return ""
		} else if err != nil {
			global.SendLogs("error", "从redis获取验证码报错", err)
			if valInterface, okv := global.Store.Get(key); okv {
				valstr, ok := valInterface.(string)
				if ok {
					val = valstr
				} else {
					err = errors.New("验证码类型不是string")
					global.SendLogs("error", "从本地缓存获取验证码报错", err)
				}
			}
		}
		return val
	}
	return ""
}
