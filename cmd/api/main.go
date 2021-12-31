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

	"go.uber.org/zap"

	"quiz-api-service/internal/config"
	"quiz-api-service/internal/pkg"
	"quiz-api-service/internal/services"
	"quiz-api-service/internal/utils"
)

func main() {

	logger, err := utils.SetUpLogger()
	if err != nil {
		log.Fatalf("somethign went wrong setting up logger for api: %+v", err)
	}

	defer logger.Sync()

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

	logger.Info("🚀 Setting up migrations")

	// set up schema
	err = quizDBConn.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&pkg.User{},
		&pkg.UserGames{},
		&pkg.Game{},
		&pkg.Question{},
		&pkg.Answer{},
		&pkg.UserAnswers{},
	)
	if err != nil {
		logger.Fatal("something went wrong migrating schema", zap.Error(err))
	}

	// set up dummy data
	quizDBConn.CreateInBatches(
		[]pkg.User{
			{Name: "David Smith", Age: 54, Username: "david54", Password: "pass54"},
			{Name: "John Doe", Age: 14, Username: "johndoe14", Password: "pass14"},
			{Name: "Steve Borg", Age: 28, Username: "steve321", Password: "qwerty321"},
		},
		3) // TODO: solve issue of duplicate rows, include a migrations table.

	logger.Info(fmt.Sprintf("✅ Applied migrations to %s db.", quizDBConn.Migrator().CurrentDatabase()))

	_ = services.NewQuizService(quizDBConn)

	services.Execute()
}
