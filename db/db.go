package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alpody/fiber-realworld/model"
	iris "github.com/caretdev/gorm-iris"
	iriscontainer "github.com/caretdev/testcontainers-iris-go"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() *gorm.DB {
	// dsn := "iris://_SYSTEM:SYS@iris:1972/USER"
	connectionString := os.Getenv("GORM_DSN")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * 10, // Slow SQL threshold
			LogLevel:                  logger.Info,           // Log level
			IgnoreRecordNotFoundError: false,                 // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                  // Disable color
		},
	)

	db, err := gorm.Open(iris.New(iris.Config{
		DSN: connectionString,
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Println("storage err: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("storage err: ", err)
	}

	sqlDB.SetMaxIdleConns(3)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

var container *iriscontainer.IRISContainer

func TestDB(useContainer bool) *gorm.DB {
	var err error
	var connectionString = "iris://_SYSTEM:SYS@localhost:1972/USER"

	if useContainer {
		options := []testcontainers.ContainerCustomizer{
			iriscontainer.WithNamespace("TEST"),
			iriscontainer.WithUsername("testuser"),
			iriscontainer.WithPassword("testpassword"),
		}
		ctx := context.Background()
		container, err = iriscontainer.RunContainer(ctx, options...)
		if err != nil {
			log.Println("Failed to start container:", err)
			os.Exit(1)
		}
		connectionString = container.MustConnectionString(ctx)
		fmt.Println("Container started: ", connectionString)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // Disable color
		},
	)

	db, err := gorm.Open(iris.New(iris.Config{
		DSN: connectionString,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if !useContainer {
		_ = db.Exec("DROP DATABASE TEST").Error
		_ = db.Exec("CREATE DATABASE TEST").Error
		_ = db.Exec("USE DATABASE TEST").Error
	}

	if err != nil {
		fmt.Println("storage err: ", err)
	}
	return db
}

func DropTestDB() error {
	if container != nil {
		container.Terminate(context.Background())
	}
	container = nil
	return nil
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Article{},
		&model.Comment{},
		&model.Tag{},
	)
}
