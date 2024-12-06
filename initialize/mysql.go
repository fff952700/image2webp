package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"image2webp/global"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMysql() {
	// 构建 DSN (Data Source Name) 字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.Conf.MysqlInfo.User,
		global.Conf.MysqlInfo.Password,
		global.Conf.MysqlInfo.Host,
		global.Conf.MysqlInfo.Port,
		global.Conf.MysqlInfo.DB)

	// 设置 logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阀值
			LogLevel:      logger.Info,
			Colorful:      true, // 禁用色彩打印
		},
	)

	// 打开 MySQL 数据库连接
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 保持原有表名，不使用复数形式
			NameReplacer:  nil,  // 名称替换器（此处未使用）
		},
	})
	if err != nil {
		zap.S().Panicw("mysql init failed", "err", err)
	}

	// 获取底层的 *sql.DB 对象
	sqlDB, err := DB.DB()
	if err != nil {
		zap.S().Panicw("failed to get sql.DB", "err", err)
	}

	// 配置连接池参数
	sqlDB.SetMaxOpenConns(global.Conf.MysqlInfo.MaxOpenConn) // 设置最大连接数
	sqlDB.SetMaxIdleConns(global.Conf.MysqlInfo.MaxIdleConn) // 设置最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)                      // 设置连接最大生命周期

	// 将 DB 实例保存到 global.DB
	global.DB = DB

	// 输出日志
	zap.S().Info("MySQL initialized successfully")
}
