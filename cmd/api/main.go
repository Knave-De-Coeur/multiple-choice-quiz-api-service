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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"quiz-api-service/internal/api"
	"quiz-api-service/internal/config"
	"quiz-api-service/internal/pkg"
	"quiz-api-service/internal/services"
	"quiz-api-service/internal/utils"
)

var QuizService *services.QuizService

func main() {

	logger, err := utils.SetUpLogger()
	if err != nil {
		log.Fatalf("somethign went wrong setting up logger for api: %+v", err)
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("something went wrong deferring the close to the logger: %v", err)
		}
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

	portNum, err := strconv.Atoi(config.CurrentConfigs.Port)
	if err != nil {
		logger.Fatal(fmt.Sprintf("port config not int %d", err))
		return
	}

	QuizService = services.NewQuizService(quizDBConn, logger, services.QuizServiceSettings{
		Port:     portNum,
		Hostname: config.CurrentConfigs.Host,
	})
}

// REST FUNCTIONS

// One of the endpoints that shows the homepage
func homePage(res http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(res, "Welcome to the HomePage! Go to the console to start playing.")
}

// Login endpoint function that checks username and password and sets user appropriately
func login(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var loginReq api.LoginRequest
	_ = json.Unmarshal(reqBody, &loginReq)

	currentUser, err := QuizService.GetUserByUsername(loginReq.Username)
	if err != nil {
		messageResponse := api.Response{
			Message: err.Error(),
			Error:   true,
		}

		_ = json.NewEncoder(res).Encode(messageResponse)
		return
	}

	foundError := false
	message := ""
	if currentUser.ID == 0 {
		message = "User not found: " + loginReq.Username
	} else if currentUser.Password != strings.TrimSpace(loginReq.Password) {
		message = "Password is wrong try again."
	} else {
		foundError = true
		message = "Successfully logged in with user: " + currentUser.Username
	}

	messageResponse := api.Response{
		Message: message,
		Error:   foundError,
	}

	_ = json.NewEncoder(res).Encode(messageResponse)

}

// Logout endpoint function that simply removes the CurrentUserID
func logout(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var logoutRequest api.LogoutRequest
	_ = json.Unmarshal(reqBody, &logoutRequest)

	message := ""
	foundError := false

	resMessage := api.Response{
		Message: message,
		Error:   foundError,
	}

	_ = json.NewEncoder(res).Encode(resMessage)
}

// Add player rest endpoint simply adds the posted user details to the global param ListOfUsers
func addNewUser(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)
	var newPlayer pkg.User
	_ = json.Unmarshal(reqBody, &newPlayer)

	message := ""
	errorFound := false

	err := QuizService.InsertUser(&newPlayer)
	if err != nil {
		message = err.Error()
		errorFound = true
	}

	responseMessage := api.Response{
		Message: message,
		Error:   errorFound,
	}

	_ = json.NewEncoder(res).Encode(responseMessage)
}

// GET rest endpoint func that simply displays list of users in json
func showPlayers(res http.ResponseWriter, _ *http.Request) {
	users, err := QuizService.GetUsers()
	if err != nil {
		responseMessage := api.Response{
			Message: err.Error(),
			Error:   true,
		}
		_ = json.NewEncoder(res).Encode(responseMessage)
	}
	_ = json.NewEncoder(res).Encode(users)
}
