package global

import (
	"auth/utils"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

/*
	定义并发锁
*/

// 定义接口
type Interface interface {
	Get() bool
	Block(second int64) bool
	Release() bool
	ForseRelease()
}

// 定义锁结构
type lock struct {
	context context.Context
	name    string // 锁名称
	owner   string // 锁标识
	seconds int64  // 有效期
}

// lock 实现 Interface,获取锁
func (l *lock) Get() bool {
	fmt.Println(l.name, l.owner, time.Duration(l.seconds)*time.Second)
	return App.Redis.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

// 阻塞一段时间，尝试获取锁
func (l *lock) Block(second int64) bool {
	starting := time.Now().Unix()
	for {
		if !l.Get() {
			time.Sleep(time.Duration(1))
			if time.Now().Unix()-second >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

// 释放锁
func (l *lock) Release() bool {
	luaScript := redis.NewScript(releaseLockLuaStript)
	result := luaScript.Run(l.context, App.Redis, []string{l.name}, l.owner).Val().(int64)
	return result != 0
}

// 强制释放锁
func (l *lock) ForseRelease() {
	App.Redis.Del(l.context, l.name)
}

// 释放锁Lua脚本，防止客户端都能解锁
const releaseLockLuaStript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
	return redis.call("del,KEYS[1]")
else
	return 0
end
`

// 生成锁
func Lock(name string, seconds int64) Interface {
	return &lock{
		context.Background(),
		name,
		utils.RandString(16),
		seconds,
	}
}
