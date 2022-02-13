/*
Copyright © 2020 NAME HERE <alexanderm1496@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
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
		_ = logger.Sync()
		// if err != nil {
		// 	fmt.Printf("something went wrong deferring the close to the logger: %v", err)
		// }
	}(logger)

	logger.Info("🚀 connecting to db")

	quizDBConn, err := utils.SetUpDBConnection(
		config.CurrentConfigs.DBUser,
		config.CurrentConfigs.DBPassword,
		config.CurrentConfigs.Host,
		config.CurrentConfigs.DBName,
		logger,
	)
	if err != nil {
		logger.Fatal("exiting application...", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("✅ Setup connection to %s db.", quizDBConn.Migrator().CurrentDatabase()))

	logger.Info("🚀 Running migrations")

	if err = utils.SetUpSchema(quizDBConn, logger); err != nil {
		logger.Fatal(err.Error())
	}

	db, err := quizDBConn.DB()
	if err != nil {
		logger.Fatal("something went wrong getting the database conn from gorm", zap.Error(err))
	}

	if err = utils.RunUpMigrations(db, logger); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("✅ Applied migrations to %s db.", quizDBConn.Migrator().CurrentDatabase()))

	routes, err := setUpRoutes(quizDBConn, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = routes.Run()
	if err != nil {
		logger.Fatal("something went wrong setting up router")
	}
}

// setUpRoutes adds routes and returns gin engine
func setUpRoutes(quizDBConn *gorm.DB, logger *zap.Logger) (*gin.Engine, error) {

	portNum, err := strconv.Atoi(config.CurrentConfigs.Port)
	if err != nil {
		logger.Error(fmt.Sprintf("port config not int %d", err))
		return nil, err
	}

	userService := services.NewUserService(quizDBConn, logger, services.UserServiceSettings{
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

	handlers.NewUserHandler(userService).UserRoutes(r.Group("/"))

	return r, nil
}
