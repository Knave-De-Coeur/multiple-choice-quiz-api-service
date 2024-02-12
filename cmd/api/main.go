package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/knave-de-coeur/user-api-service/internal/config"
	"github.com/knave-de-coeur/user-api-service/internal/handlers"
	"github.com/knave-de-coeur/user-api-service/internal/middleware"
	"github.com/knave-de-coeur/user-api-service/internal/services"
	"github.com/knave-de-coeur/user-api-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {

	logger, err := utils.SetUpLogger()
	if err != nil {
		log.Fatalf("somethign went wrong setting up logger for api: %+v", err)
	}

	utils.Check(logger.Sync)

	logger.Info("🚀 connecting to db")

	dbConnection, err := utils.SetUpDBConnection(
		config.CurrentConfigs.DBUser,
		config.CurrentConfigs.DBPassword,
		config.CurrentConfigs.Host,
		config.CurrentConfigs.DBName,
		logger,
	)
	if err != nil {
		logger.Fatal("exiting application...", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("✅ Setup connection to %s db.", dbConnection.Migrator().CurrentDatabase()))

	logger.Info("🚀 Running migrations")

	if err = utils.SetUpSchema(dbConnection, logger); err != nil {
		logger.Fatal(err.Error())
	}

	db, err := dbConnection.DB()
	if err != nil {
		logger.Fatal("something went wrong getting the database conn from gorm", zap.Error(err))
	}

	if err = utils.RunUpMigrations(db, logger); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("✅ Applied migrations to %s db.", dbConnection.Migrator().CurrentDatabase()))

	var nc *nats.Conn

	if config.CurrentConfigs.NatsURL != "" {
		logger.Info("🚀 Setting up nats connection.")
		// Connect to a server
		// nc, err := nats.Connect("nats://127.0.0.1:4222")
		nc, err = nats.Connect(config.CurrentConfigs.NatsURL)
		if err != nil {
			logger.Fatal(fmt.Sprintf("❌ Failed to set up nats %s", err.Error()))
		}

		logger.Info("✅ Connected to nats!")

		utils.Check(nc.Drain)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.CurrentConfigs.RedisAddress,
		Password: config.CurrentConfigs.RedisPassword,
		DB:       config.CurrentConfigs.RedisDB,
	})

	routes, err := setUpRoutes(dbConnection, redisClient, nc, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err = routes.Run(); err != nil {
		logger.Fatal("something went wrong setting up router")
	}
}

// setUpRoutes adds routes and returns gin engine
func setUpRoutes(dbConn *gorm.DB, redisClient *redis.Client, nc *nats.Conn, logger *zap.Logger) (*gin.Engine, error) {

	portNum, err := strconv.Atoi(config.CurrentConfigs.Port)
	if err != nil {
		logger.Error(fmt.Sprintf("port config not int %v", err))
		return nil, err
	}

	userService := services.NewUserService(dbConn, nc, logger, services.UserServiceSettings{
		Port:      portNum,
		Hostname:  config.CurrentConfigs.Host,
		JWTSecret: config.CurrentConfigs.JWTSecret,
	})

	r := gin.New()

	r.Use(gin.Logger())

	// r.Use(gin.Middleware)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	authMiddleware := middleware.NewAuthMiddleware(config.CurrentConfigs.JWTSecret)

	handlers.NewUserHandler(userService, authMiddleware, redisClient, nc).SetUpRoutes(r.Group("/api/v1"))

	return r, nil
}
