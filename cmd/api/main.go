package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"quiz-api-service/internal/config"
	"quiz-api-service/internal/handlers"
	"quiz-api-service/internal/services"
	"quiz-api-service/internal/utils"
)

func main() {

	logger, err := utils.SetUpLogger()
	if err != nil {
		log.Fatalf("somethign went wrong setting up logger for api: %+v", err)
	}

	defer func(logger *zap.Logger) {
		err = logger.Sync()
		panic(fmt.Sprintf("something went wrong with logger %w", err))
	}(logger)

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

	// Connect to a server
	// nc, err := nats.Connect("nats://127.0.0.1:4222")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Fatal(fmt.Sprintf("❌ Failed to set up nats %s", err.Error()))
	}

	defer nc.Drain()

	routes, err := setUpRoutes(dbConnection, nc, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err = routes.Run(); err != nil {
		logger.Fatal("something went wrong setting up router")
	}
}

// setUpRoutes adds routes and returns gin engine
func setUpRoutes(dbConn *gorm.DB, nc *nats.Conn, logger *zap.Logger) (*gin.Engine, error) {

	portNum, err := strconv.Atoi(config.CurrentConfigs.Port)
	if err != nil {
		logger.Error(fmt.Sprintf("port config not int %d", err))
		return nil, err
	}

	userService := services.NewUserService(dbConn, nc, logger, services.UserServiceSettings{
		Port:     portNum,
		Hostname: config.CurrentConfigs.Host,
	})

	gameService := services.NewGameService(dbConn, logger, services.GameServiceSettings{
		Port:     portNum,
		Hostname: config.CurrentConfigs.Host,
	})

	r := gin.New()

	r.Use(gin.Logger())

	// r.Use(gin.Middleware)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handlers.NewUserHandler(userService, nc).UserRoutes(r.Group("/api/v1"))
	handlers.NewGameHandler(gameService).GameRoutes(r.Group("/api/v1"))

	return r, nil
}
