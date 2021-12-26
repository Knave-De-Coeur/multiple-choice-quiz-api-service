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
	"log"
	"os"

	"go.uber.org/zap"

	"quiz-api-service/internal/config"
	"quiz-api-service/internal/services"
	"quiz-api-service/internal/utils"
)

func main() {

	config.Host = os.Getenv("HOST")
	config.DefaultPort = os.Getenv("PORT")
	config.DBConnectionString = os.Getenv("DB_CONNECTION")

	logger, err := utils.SetUpLogger()
	if err != nil {
		log.Fatalf("somethign went wrong setting up logger for api: %+v", err)
	}

	defer logger.Sync()

	_, err = utils.SetUpDBConnection(config.DBConnectionString, logger)
	if err != nil {
		logger.Error("error encountered setting up db conn:", zap.Error(err))
		os.Exit(1)
	}

	services.Execute()
}
