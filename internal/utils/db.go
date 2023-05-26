package utils

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrateMysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/knave-de-coeur/user-api-service/internal/config"
	"github.com/knave-de-coeur/user-api-service/internal/pkg"
)

// GetDBConnectionString uses configs to generate a connection string to the db
func GetDBConnectionString(user, pass, host, dbName string) string {
	return fmt.Sprintf(config.CurrentConfigs.DBConnectionFormat, user, pass, host, dbName)
}

// GetDBConnection uses string passed to connect to mysql database
func GetDBConnection(user, pass, host, dbName string, logger *zap.Logger) (*gorm.DB, error) {

	dsn := GetDBConnectionString(user, pass, host, dbName)

	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Unix(time.Now().Unix(), 0).UTC()
		},
	})
	if err != nil {
		logger.Error("❌ something went wrong getting the db connection", zap.String("method", "GetDBConnection"), zap.Error(err))
		return nil, err
	}

	return db, nil
}

// SetUpDBConnection gets the connection and applies all the configs to it
func SetUpDBConnection(user, pass, host, dbName string, logger *zap.Logger) (*gorm.DB, error) {

	db, err := GetDBConnection(user, pass, host, dbName, logger)
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

// SetUpSchema uses the current structs to generate tables
func SetUpSchema(db *gorm.DB, logger *zap.Logger) (err error) {
	// set up schema
	err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&pkg.User{},
	)
	if err != nil {
		logger.Error("something went wrong migrating schema", zap.Error(err))
	}

	return
}

// RunUpMigrations uses db connection and locations of migration to run all the up migrations
func RunUpMigrations(db *sql.DB, logger *zap.Logger) (err error) {

	driver, _ := migrateMysql.WithInstance(db, &migrateMysql.Config{})
	fileLocation := fmt.Sprintf("file://%s", config.CurrentConfigs.MigrationsDir)
	m, err := migrate.NewWithDatabaseInstance(fileLocation, "sql", driver)
	if err != nil {
		logger.Error("❌ failed to get migration instance",
			zap.Error(err),
			zap.String("fileLocation", fileLocation),
		)
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		logger.Info("✅ no migrations were ran on db, all up to date!")
	} else if err != nil {
		logger.Error("❌ failed to run migrations", zap.Error(err))
		return err
	}

	return nil
}
