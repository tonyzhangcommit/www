package bootstrap

import (
	"context"
	"fmt"
	"strconv"
	"userservice/global"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
	初始化 mysql 数据库
*/

func InitializeDatabase() {
	if global.App.Config.Database.Driver == "mysql" {
		global.App.DB = InitMysql()
		return
	} else {
		// 这里默认为mysql， 如果存在其他数据库，这里编写对应的逻辑，gorm 支持 MySQL, PostgreSQL, SQLite, SQL Server 和 TiDB
		global.App.DB = InitMysql()
		return
	}
}

func InitMysql() *gorm.DB {
	dbConfig := global.App.Config.Database
	if dbConfig.DBName == "" {
		panic("数据库配置文件错误")
	}
	dsn := dbConfig.Username + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" +
		dbConfig.DBName + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: nil,
	}); err != nil {
		errlog := global.LogMessage{
			Level:   "error",
			Message: err,
		}
		errlog.SendInfoToRabbitMQ()
		panic(fmt.Sprintf("mysql connect failed! err : %v\n", err))
	} else {
		// 设置连接池连接数和最大连接数
		DB, _ := db.DB()
		DB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		DB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		initializeTables(db)
		return db
	}
}

func initializeTables(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		panic(fmt.Sprintf("初始化数据库表失败，err:%v", err))
	}
}

func InitializeRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.App.Config.Redis.Host + ":" + global.App.Config.Redis.Port,
		Password: "", // 密码，没有则留空
		DB:       2,  // 使用默认DB
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
