package utils

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"quiz-api-service/internal/config"
)

// TODO: add a string generator and split the connection string into host, db, user and pass - figure out how api container should connect to db container

// GetDBConnection uses string passed to connect to mysql database
func GetDBConnection(conn string, logger *zap.Logger) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		logger.Error("❌ something went wrong getting the db connection", zap.String("method", "GetDBConnection"), zap.Error(err))
		return nil, err
	}

	return db, nil
}

// SetUpDBConnection gets the connection and applies all the configs to it
func SetUpDBConnection(conn string, logger *zap.Logger) (*gorm.DB, error) {

	db, err := GetDBConnection(conn, logger)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("❌ something went wrong extracting the db from the gorm conn", zap.String("method", "SetUpDBConnection"), zap.Error(err))
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(config.CurrentConfigs.MaxIdleConnections)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(config.CurrentConfigs.MaxConnections)

	lifeTime, err := time.ParseDuration(fmt.Sprintf("%vh", config.CurrentConfigs.MaxLifetime))
	if err != nil {
		logger.Error("❌ something went wrong formatting the lifetime duration from the configurations", zap.String("method", "SetUpDBConnection"), zap.Error(err))
		return nil, err
	}

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(lifeTime)

	return db, nil
}
