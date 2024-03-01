package timingtask

import (
	"auth/global"
	"auth/middleware"
	"fmt"
	"time"
)

// cleanupLimiters 定期清理超过 6 小时未使用的 limiter
func CleanupLimiters() {
	for {
		time.Sleep(30 * time.Minute) // 每 30 分钟执行一次清理
		now := time.Now()
		middleware.IpLimiter.Range(func(key, value interface{}) bool {
			global.SendLogs("info", "定时任务执行检测iplimiter")
			limiterInfo := value.(*middleware.Limiterinfo)
			fmt.Println(limiterInfo)
			if now.Sub(limiterInfo.LastUsed) > 2*time.Hour {
				middleware.IpLimiter.Delete(key)
				global.SendLogs("info", "定时任务执行删除iplimiter")
			}
			return true
		})
	}
}
