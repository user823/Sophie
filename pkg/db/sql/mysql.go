package sql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MysqlConfig struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}

func NewMysqlDB(config any) (*gorm.DB, error) {
	cfg, ok := config.(*MysqlConfig)
	if !ok {
		return nil, ErrConfigTypeInvalid
	}

	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s`,
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Database,
		true,
		"Local",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 cfg.Logger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(cfg.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	return db, nil
}

//mysqlLoggerConfig := logger.Config{
//SlowThreshold:             1 * time.Second,
//Colorful:                  false,
//IgnoreRecordNotFoundError: true,
//ParameterizedQueries:      true,
//LogLevel:                  (logger.LogLevel)(opts.MySQLOptions.LogLevel),
//}
//mysql := &sql.MysqlConfig{
//Host:                  opts.MySQLOptions.Host,
//Username:              opts.MySQLOptions.Username,
//Password:              opts.MySQLOptions.Password,
//Database:              opts.MySQLOptions.Database,
//MaxConnectionLifeTime: opts.MySQLOptions.MaxConnectionLifeTime,
//MaxOpenConnections:    opts.MySQLOptions.MaxOpenConnections,
//MaxIdleConnections:    opts.MySQLOptions.MaxIdleConnections,
//LogLevel:              opts.MySQLOptions.LogLevel,
//Logger:                logger.New(log.Default(), mysqlLoggerConfig),
//}
