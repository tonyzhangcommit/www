package bootstrap

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"userservice/global"
	"userservice/models"

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
		// 初始化视频类型
		initVideoType(db)
		// 初始化会员表
		initMembership(db)
		// 初始化秒杀活动
		createFlashSaleEvent(db)
		// 关联秒杀商品

		return db
	}
}

func initializeTables(db *gorm.DB) {
	err := db.AutoMigrate(&models.VideoType{}, &models.Video{}, &models.Membership{}, &models.FlashSaleEvent{}, &models.FlashSaleEvent{}, &models.FlashSaleEventProduct{})
	if err != nil {
		panic(fmt.Sprintf("初始化数据库表失败，err:%v", err))
	}
}

func InitializeRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.App.Config.Redis.Host + ":" + global.App.Config.Redis.Port,
		Password: "", // 密码，没有则留空
		DB:       0,  // 使用默认DB
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

func initVideoType(db *gorm.DB) {
	vT := make(map[string]string)
	vT["电影"] = "电影"
	vT["电视剧"] = "电视剧"
	vT["动漫"] = "动漫"
	if err := db.Where("name = ?", "电影").First(&models.VideoType{}).Error; err != nil {
		for k, v := range vT {
			if err := db.Create(&models.VideoType{Name: k, Description: v}).Error; err != nil {
				global.SendLogs("error", "初始化视频类型失败", err)
				return
			}
		}
	}
}

func initMembership(db *gorm.DB) {
	type menber struct {
		Name     string
		Price    float64
		Duration int
	}
	type memberlist []menber
	m1 := menber{
		"月会员",
		30.0,
		31,
	}
	m2 := menber{
		"季会员",
		100.0,
		93,
	}
	m3 := menber{
		"年会员",
		300.0,
		365,
	}
	var ml memberlist
	ml = append(ml, m1, m2, m3)
	var count int64
	db.Table("membership").Count(&count)
	if count == 0 {
		for _, v := range ml {
			if err := db.Create(&models.Membership{Name: v.Name, Price: v.Price, Duration: v.Duration}).Error; err != nil {
				global.SendLogs("error", "初始化会员失败", err)
				return
			}
		}
	}

}

// 定义一个秒杀活动
func createFlashSaleEvent(db *gorm.DB) {
	var mber models.Membership
	if err := db.Where("name = ?", "月会员").First(&mber).Error; err != nil {
		global.SendLogs("error", "没有产品信息", err)
		return
	}
	totalcount := 10000
	userlimit := 1
	starttime := time.Now()
	endtime := starttime.Add(time.Hour * 24)
	flashevent := models.FlashSaleEvent{
		Name:      "测试秒杀活动",
		Condition: "月会员",
		StartTime: starttime,
		EndTime:   endtime,
	}
	if err := db.Create(&flashevent).Error; err != nil {
		global.SendLogs("error", "创建秒杀活动失败", err)
		return
	}
	global.SendLogs("info", "创建秒杀活动成功")
	// 添加商品和活动关联信息
	if err := db.Create(&models.FlashSaleEventProduct{
		EventID:           flashevent.ID,
		ProductID:         mber.ID,
		OriginalPrice:     mber.Price,
		FlashSalePrice:    mber.Price * 0.8,
		Quantity:          totalcount,
		RemainingQuantity: totalcount,
		LimitPerUser:      userlimit,
	}).Error; err != nil {
		global.SendLogs("error", "创建秒杀活动失败", err)
		return
	}
	global.SendLogs("info", "绑定商品成功")
}
