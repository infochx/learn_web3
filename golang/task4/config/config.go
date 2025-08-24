package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Mysql MysqlConfig `yaml:"mysql"`
	Jwt   JwtConfig   `yaml:"jwt"`
}
type MysqlConfig struct {
	DSN          string `yaml:"dsn"`
	MaxIdleConns int64  `yaml:"maxIdleConns"`
	MaxOpenConns int64  `yaml:"maxOpenConns"`
}

type JwtConfig struct {
	JwtSecret string `yaml:"jwtSecret"`
}

var DB *gorm.DB

func InitDB() error {
	//加载数据库配置文件
	//读取配置文件内容
	path := "config.yaml"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read file %s error: %v", path, err))
	}
	//创建一个空的config，yaml.Unmarshal需要
	dbConfig := &Config{}
	err = yaml.Unmarshal(data, dbConfig)
	if err != nil {
		panic(fmt.Errorf("unmarshal file data %s error: %v", path, err))
	}
	if dbConfig.Mysql.DSN == "" {
		panic(fmt.Errorf("dsn is empty"))
	}

	if err != nil {
		panic("配置文件加载错误：" + err.Error())
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}

	//连接数据库
	DB, err = gorm.Open(mysql.Open(dbConfig.Mysql.DSN), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().In(loc)
		},
	})
	if err != nil {
		return err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(int(dbConfig.Mysql.MaxOpenConns))
	sqlDB.SetMaxIdleConns(int(dbConfig.Mysql.MaxIdleConns))
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	return err
}
