package database

import (
	"errors"
	"fmt"
	"guide-u/cafe-crawl/config"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var connector *gorm.DB

func newLogger() logger.Interface {
	return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		LogLevel:                  logger.Error,
		IgnoreRecordNotFoundError: true, // 레코드를 찾을 수 없을 경우 에러를 출력하는 지
		Colorful:                  true,
		ParameterizedQueries:      true, // 쿼리에 매개변수를 포함하지 않는지
	})
}

func getConnectionString(config databaseConfig) string {
	var dsn string

	switch config.dbtype {
	case strings.ToLower("mysql"):
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
			config.user, config.password, config.host, config.port, config.dbname,
		)
	case strings.ToLower("postgresql"):
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
			config.host, config.user, config.password, config.dbname, config.port,
		)
	}

	return dsn
}

func DatabaseConnection() error {
	if !strings.Contains(config.Environment("DATABASE_CREDENTIALS"), ":") {
		return errors.New("cannot split database credentials. need id:password")
	}

	credentials := strings.Split(config.Environment("DATABASE_CREDENTIALS"), ":")

	config := databaseConfig{
		host:     config.Environment("DATABASE_HOST"),
		port:     config.Environment("DATABASE_PORT"),
		user:     credentials[0],
		password: credentials[1],
		dbname:   config.Environment("DATABASE_NAME"),
		dbtype:   config.Environment("DATABASE_TYPE"),
	}

	dsn := getConnectionString(config)

	var err error
	switch config.dbtype {
	case strings.ToLower("mysql"):
		connector, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger(),
		})
		connector.AutoMigrate(Post{})
	case strings.ToLower("postgresql"):
		connector, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger(),
		})

		connector.AutoMigrate(&Post{})
	default:
		return errors.New("cannot read database type")
	}
	if err != nil {
		return err
	}

	log.Printf("database ok, type : %s\n", config.dbtype)
	return nil
}
