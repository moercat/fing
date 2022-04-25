package db

import (
	"fing/pkg/config"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"xorm.io/xorm"
)

var (
	Main        *xorm.Engine
	Gain        *gorm.DB
	RedisClient *redis.Client
	EsClient    *elastic.Client
)

func init() {
	var err error

	Main, err = xorm.NewEngine("mysql", config.Config.DataSource.Main)
	if err != nil {
		panic(err)
	}

	// 初始化GORM日志配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	Gain, err = gorm.Open(mysql.Open(config.Config.DataSource.Main), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         config.Config.Redis.Addr,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxConnAge:   10 * time.Minute,
	})
	_, err = RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Config.Es.EsUrl),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	}
	if config.Config.Es.EsUsername != "" {
		options = append(options, elastic.SetBasicAuth(config.Config.Es.EsUsername, config.Config.Es.EsPassword))
	}

	EsClient, err = elastic.NewClient(options...)
	if err != nil {
		panic(err)
	}

}
