package bootstrap

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"userservice/global"
	"userservice/models"
	"userservice/utils"

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
		// 初始化 role permission
		initializeRPSA(db)
		creatTestAdmin(db)
		return db
	}
}

func initializeTables(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{}, &models.Profile{}, &models.AgentManagement{}, &models.Role{}, &models.Permission{}, &models.UserActivity{}, &models.VerificationCodeRecord{},
	)
	if err != nil {
		panic(fmt.Sprintf("初始化数据库表失败，err:%v", err))
	}
}

/*
项目初始化角色，权限，新建超级管理员
通过判断roles表，如果为空，则执行
*/

func initializeRPSA(db *gorm.DB) {
	var count int64
	if err := db.Find(&[]models.Role{}).Count(&count).Error; err != nil {
		panic(fmt.Sprintf("查询Role表失败，err:%v", err))
	}
	if count == 0 {
		// 初始化role 和 premission table
		// 获取所有权限名称
		permissmap := make(map[string]string)
		initPermission := global.App.Config.Permission
		for _, rolePermission := range [][]string{initPermission.SuperAdmin, initPermission.Admin, initPermission.RegularUser, initPermission.MonthlyVip, initPermission.QuarterlyVip, initPermission.AnnualVip} {
			for _, perm := range rolePermission {
				// 分割权限名和权限解释
				parts := strings.SplitN(perm, "-", 2)
				if len(parts) == 2 {
					// 权限存在重复的，这里加一步过滤
					if _, ok := permissmap[parts[0]]; !ok {
						permissmap[parts[0]] = parts[1]
					}
				}
			}
		}
		for name, desc := range permissmap {
			tempPermiss := models.Permission{
				PermissionName: name,
				Description:    desc,
			}
			// 初始化权限表信息
			if err := db.Create(&tempPermiss).Error; err != nil {
				fmt.Println("创建失败", tempPermiss)
			}
		}
		// 初始化角色表
		// 超级管理员
		roles := []models.Role{}
		rolePermissionMap := make(map[string][]models.Permission)
		// 初始化 superAdmin
		var superadminpernamelist []string
		var superadminperstructlist []models.Permission

		for _, rolePermission := range [][]string{initPermission.SuperAdmin} {
			for _, perm := range rolePermission {
				// 分割权限名和权限解释
				parts := strings.SplitN(perm, "-", 2)
				if len(parts) == 2 {
					superadminpernamelist = append(superadminpernamelist, parts[0])
				}
			}
		}
		// 查询
		db.Where("permissionname IN ?", superadminpernamelist).Find(&superadminperstructlist)
		rolePermissionMap["superAdmin"] = superadminperstructlist
		// 初始化 admin
		var adminpernamelist []string
		var adminperstructlist []models.Permission

		for _, rolePermission := range [][]string{initPermission.Admin} {
			for _, perm := range rolePermission {
				// 分割权限名和权限解释
				parts := strings.SplitN(perm, "-", 2)
				if len(parts) == 2 {
					adminpernamelist = append(adminpernamelist, parts[0])
				}
			}
		}
		// 查询
		db.Where("permissionname IN ?", adminpernamelist).Find(&adminperstructlist)
		rolePermissionMap["admin"] = adminperstructlist
		// 初始化 regularUser
		var regularUserpernamelist []string
		var regularUserperstructlist []models.Permission

		for _, rolePermission := range [][]string{initPermission.RegularUser} {
			for _, perm := range rolePermission {
				// 分割权限名和权限解释
				parts := strings.SplitN(perm, "-", 2)
				if len(parts) == 2 {
					regularUserpernamelist = append(regularUserpernamelist, parts[0])
				}
			}
		}
		// 查询
		db.Where("permissionname IN ?", regularUserpernamelist).Find(&regularUserperstructlist)
		rolePermissionMap["regularUser"] = regularUserperstructlist
		// 初始化 monthlyVip
		var monthlypernamelist []string
		var monthlyperstructlist []models.Permission

		for _, rolePermission := range [][]string{initPermission.MonthlyVip} {
			for _, perm := range rolePermission {
				// 分割权限名和权限解释
				parts := strings.SplitN(perm, "-", 2)
				if len(parts) == 2 {
					monthlypernamelist = append(monthlypernamelist, parts[0])
				}
			}
		}
		// 查询
		db.Where("permissionname IN ?", monthlypernamelist).Find(&monthlyperstructlist)
		rolePermissionMap["monthlyVip"] = monthlyperstructlist
		// 初始化 quarterlyVip
		var quarterlypernamelist []string
		var quarterlyperstructlist []models.Permission

		for _, rolePermission := range [][]string{initPermission.QuarterlyVip} {
			for _, perm := range rolePermission {
				// 分割权限名和权限解释
				parts := strings.SplitN(perm, "-", 2)
				if len(parts) == 2 {
					quarterlypernamelist = append(quarterlypernamelist, parts[0])
				}
			}
		}
		// 查询
		db.Where("permissionname IN ?", quarterlypernamelist).Find(&quarterlyperstructlist)
		rolePermissionMap["quarterlyVip"] = quarterlyperstructlist
		// 初始化 quarterlyVip
		var annualpernamelist []string
		var annualperstructlist []models.Permission

		for _, rolePermission := range [][]string{initPermission.AnnualVip} {
			for _, perm := range rolePermission {
				// 分割权限名和权限解释
				parts := strings.SplitN(perm, "-", 2)
				if len(parts) == 2 {
					annualpernamelist = append(annualpernamelist, parts[0])
				}
			}
		}
		// 查询
		db.Where("permissionname IN ?", annualpernamelist).Find(&annualperstructlist)
		rolePermissionMap["annualVip"] = annualperstructlist

		for _, name := range global.App.Config.Roles.NameList {
			roles = append(roles, models.Role{RoleName: name, Permissions: rolePermissionMap[name]})
		}
		for _, role := range roles {
			if err := db.Create(&role).Error; err != nil {
				fmt.Println("创建失败", err)
			}
		}

	} else {
		fmt.Println("table Role had init")
	}
}

// 初始化数据库，创建默认管理员账户
func creatTestAdmin(db *gorm.DB) {
	var admin models.User
	if err := db.Where("username = ?", "desupadmin").First(&admin).Error; err != nil {
		// 没有测试用户，开始新建
		pwd := "123456!@#"
		if hashpwd := utils.BcryptMake([]byte(pwd)); hashpwd != "" {
			var superAdminRole models.Role
			if err := db.Where("rolename = ?", "superAdmin").First(&superAdminRole).Error; err != nil {
				global.SendLogs("error", "查询角色信息错误，初始管理员失败", err)
				return
			} else {
				admin.Username = "desupadmin"
				admin.Password = hashpwd
				admin.Roles = []models.Role{superAdminRole}
				admin.AgentCode = utils.GenerateRCode(6)
				if err := db.Create(&admin).Error; err != nil {
					global.SendLogs("error", "创建测试管理员失败", err)
					return
				}
			}
		} else {
			global.SendLogs("error", "加密失败导致创建初始化管理员失败")
		}
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
