/*
Package cmd Copyright © 2020 NAME HERE <alexanderm.1496@gmail.com>

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
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Endpoint is the hostname that we'll eb perforing rest requests to.
const Endpoint = "http://localhost:9990"

// User is generic user that launches and signs up
type User struct {
	ID               int    `json:"Id"`
	Name             string `json:"Name"`
	Age              int8   `json:"Age"`
	Username         string `json:"Username"`
	Password         string `json:"Password"`
	SubmittedAnswers []rune `json:"SubmittedAnswers"`
}

// LoginRequest is the parsed struct of the /login endpoint
type LoginRequest struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// LogoutRequest is the parsed struct of the /logout endpoint
type LogoutRequest struct {
	UID int
}

// Response is what is returned from the endpoints should an error occurr
type Response struct {
	Message string `json:"Message"`
	Error   bool   `json:"Password"`
}

// Question presented and evaluated
type Question struct {
	ID              int8
	Description     string // the literal question to be displayed
	CorrectAnswer   rune
	AnswerSelection map[rune]string // up to 5 possible answers
}

// SubmitAnswersRequest is the struct used to decode and encode when using the /submit-answers endpoint
type SubmitAnswersRequest struct {
	UserID           int    `json:"UserID"`
	SubmittedAnswers []rune `json:"SubmittedAnswers"`
}

// ListOfQuestions will hold the a list of Question structs to be displayed and evaluated
var ListOfQuestions []Question

// ListOfUsers all the players that have registered to the quiz
var ListOfUsers []User

// CurrentUserID The logged in user id interacting with the application
var CurrentUserID int

// Reader is the gloabl input reader for the application
var Reader bufio.Reader

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ft-quiz",
	Short: "My test for Fast Track",
	Long: `This is simple quiz where the user is prese ted with a couple questions
			and they have to select one from three to get the right answer.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		generateQuestions()  // sets questions
		generateDummyUsers() // sets users

		// create new wait group
		wg := new(sync.WaitGroup)

		// add two goroutines
		wg.Add(2)

		go handleRequests() // sets endpoints
		go runGame()        // runs game

		wg.Wait()
	},
}

func runGame() {
	fmt.Println("Welcome to Alex's quiz! Press any key followed by enter to begin.")

	reader := bufio.NewReader(os.Stdin)
	key, _ := reader.ReadString('\n')

	if len(key) > 0 {
		for {
			fmt.Println("Enter option letter and press enter")

			if CurrentUserID == 0 {
				fmt.Println("a: Sign Up")
				fmt.Println("b: Login")
				fmt.Println("c: Exit")

				optionStr, _ := reader.ReadString('\n')

				option := []rune(optionStr)

				switch option[0] {
				case 'a':
					createUser(reader)
					break
				case 'b':
					loginPrompt(reader)
					break
				case 'c':
					os.Exit(0)
					break
				default:
					fmt.Println("Invalid try again")
					break
				}
			} else {
				fmt.Println("a: Play")
				fmt.Println("b: Logout")
				fmt.Println("c: Exit")

				optionStr, _ := reader.ReadString('\n')

				option := []rune(optionStr)

				switch option[0] {
				case 'a':
					play(reader)
					break
				case 'b':
					logoutPrompt(reader)
					break
				case 'c':
					os.Exit(0)
					break
				default:
					fmt.Println("Invalid try again")
					break
				}
			}

		}
	}
}

func generateQuestions() {
	ListOfQuestions = []Question{
		{
			ID:            1,
			Description:   "What is the result of 240 / 12?",
			CorrectAnswer: 'b',
			AnswerSelection: map[rune]string{
				'a': "100",
				'b': "20",
				'c': "24",
			},
		},
		{
			ID:            2,
			Description:   "Who was the first man on the moon?",
			CorrectAnswer: 'c',
			AnswerSelection: map[rune]string{
				'a': "Buzz Aldren",
				'b': "Buzz Lightyear",
				'c': "Neil Armstrong",
			},
		},
		{
			ID:            3,
			Description:   "Who wrote the hit single, Yellow Submarine?",
			CorrectAnswer: 'c',
			AnswerSelection: map[rune]string{
				'a': "Elvis Presley",
				'b': "The Rolling Stones",
				'c': "The Beatles",
			},
		},
		{
			ID:            4,
			Description:   "In what year did Malta gain it's independence?",
			CorrectAnswer: 'c',
			AnswerSelection: map[rune]string{
				'a': "1889",
				'b': "2004",
				'c': "1964",
			},
		},
		{
			ID:            5,
			Description:   "When was Go launched?",
			CorrectAnswer: 'a',
			AnswerSelection: map[rune]string{
				'a': "2009",
				'b': "2019",
				'c': "1999",
			},
		},
	}
}

func generateDummyUsers() {
	ListOfUsers = []User{
		{
			ID:       1,
			Name:     "David Smith",
			Age:      54,
			Username: "david54",
			Password: "pass765",
			SubmittedAnswers: []rune{
				'a', 'c', 'b', 'b', 'c',
			},
		},
		{
			ID:       2,
			Name:     "John Doe",
			Age:      14,
			Username: "johndoe14",
			Password: "pass14",
			SubmittedAnswers: []rune{
				'a', 'a', 'b', 'c', 'c',
			},
		},
		{
			ID:       3,
			Name:     "Steve Bord",
			Age:      28,
			Username: "seteveb321",
			Password: "qwerty098",
			SubmittedAnswers: []rune{
				'c', 'b', 'c', 'c', 'c',
			},
		},
	}
}

func play(reader *bufio.Reader) bool {
	fmt.Printf("Press any key followed by enter to start the game you have %v questions \n", len(ListOfQuestions))

	currentUser := searchUsersByID(CurrentUserID)

	key, err := reader.ReadString('\n')
	check(err)

	if len(key) > 0 {
		fmt.Println(ListOfQuestions)
		listOfSubmittedRunes := []rune{}
		for _, q := range ListOfQuestions {
			fmt.Println("Qustion: " + string(q.ID) + " " + q.Description)
			for key, question := range q.AnswerSelection {
				fmt.Println(string(key) + ": " + question)
			}
			submittedAnswer, err := reader.ReadString('\n')
			check(err)
			answerRune := []rune(submittedAnswer)[0]
			fmt.Println("answer submitted: " + string(answerRune))
			listOfSubmittedRunes = append(listOfSubmittedRunes, answerRune)
		}
		fmt.Println("Evaluating answers...")
		time.Sleep(2 * time.Second) // for dramatic suspense

		submitAnswerRequest := SubmitAnswersRequest{
			currentUser.ID,
			listOfSubmittedRunes,
		}

		requestJSON, _ := json.Marshal(submitAnswerRequest)

		res, err := http.Post(Endpoint+"/submit-answer", "application/json", bytes.NewBuffer(requestJSON))
		check(err)

		defer res.Body.Close()

		var responseMessage Response

		json.NewDecoder(res.Body).Decode(&responseMessage)

		fmt.Println(responseMessage.Message)

	}

	return true
}

func submitAnswersAndGetResults(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var submitAnswerRequest SubmitAnswersRequest
	json.Unmarshal(reqBody, &submitAnswerRequest)

	currentUser := User{}

	// set or re-set the users answers
	for _, u := range ListOfUsers {
		if u.ID == submitAnswerRequest.UserID {
			u.SubmittedAnswers = submitAnswerRequest.SubmittedAnswers
			currentUser = u
			break
		}
	}

	fmt.Printf("Current user in submitanswers endpoint: %v ", currentUser)

	correctAnswers := 0

	for _, q := range ListOfQuestions {
		found := find(currentUser.SubmittedAnswers, q.CorrectAnswer)
		if found {
			correctAnswers++
		}
	}

	message := "You have answered " + strconv.Itoa(correctAnswers) + " out of " + strconv.Itoa(len(ListOfQuestions)) + " questions correctly!"

	responseMessage := Response{
		message,
		false,
	}

	json.NewEncoder(res).Encode(responseMessage)
}

func find(source []rune, value rune) bool {
	fmt.Printf("source: %v \n", source)
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func homePage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func createUser(reader *bufio.Reader) {
	uid := len(ListOfUsers) + 1
	newUser := User{
		ID: uid,
	}

	fmt.Println("Start creating your profile.")
	fmt.Println("Enter Name: ")
	newUser.Name, _ = reader.ReadString('\n')

	fmt.Println("Enter your age: ")
	rawAge, _ := reader.ReadString('\n')
	intAge, _ := strconv.ParseInt(rawAge, 3, 8)
	newUser.Age = int8(intAge)

	fmt.Println("Enter a Username: ")
	newUser.Username, _ = reader.ReadString('\n')
	fmt.Println("Enter a password: ")
	newUser.Password, _ = reader.ReadString('\n')

	userJSON, _ := json.Marshal(newUser)

	res, err := http.Post(Endpoint+"/new-player", "application/json", bytes.NewBuffer(userJSON))
	check(err)

	defer res.Body.Close()
}

// This will check the users input and set the CurrentUser
func loginPrompt(reader *bufio.Reader) bool {
	fmt.Println("Please enter Username:")
	inputtedUsername, _ := reader.ReadString('\n')

	fmt.Println("Please enter Password:")
	inputtedPassword, _ := reader.ReadString('\n')

	loginRequest := LoginRequest{
		inputtedUsername,
		inputtedPassword,
	}

	loginRequestJSON, _ := json.Marshal(loginRequest)

	res, err := http.Post(Endpoint+"/login", "application/json", bytes.NewBuffer(loginRequestJSON))
	check(err)

	defer res.Body.Close()

	var responseMessage Response

	json.NewDecoder(res.Body).Decode(&responseMessage)

	fmt.Printf("response of login is: %v \n", responseMessage.Message)

	return !responseMessage.Error
}

func login(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var loginReq LoginRequest
	json.Unmarshal(reqBody, &loginReq)

	currentUser := searchUsersForUsername(loginReq.Username)

	found := false
	message := ""
	if currentUser.ID == 0 {
		message = "User not found: " + loginReq.Username
	} else if currentUser.Password != strings.TrimSpace(loginReq.Password) {
		message = "Password is wrong try again."
	} else {
		CurrentUserID = currentUser.ID
		found = true
		message = "Sucessfully logged in with user: " + currentUser.Username
	}

	messageResponse := Response{
		message,
		found,
	}

	json.NewEncoder(res).Encode(messageResponse)

}

func logoutPrompt(reader *bufio.Reader) bool {
	logout := LogoutRequest{
		UID: CurrentUserID,
	}

	logoutRequestJSON, _ := json.Marshal(logout)

	res, err := http.Post(Endpoint+"/logout", "application/json", bytes.NewBuffer(logoutRequestJSON))
	check(err)

	defer res.Body.Close()

	var responseMessage Response

	json.NewDecoder(res.Body).Decode(&responseMessage)

	fmt.Println(responseMessage.Message)

	return !responseMessage.Error
}

func logout(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var logoutRequest LogoutRequest
	json.Unmarshal(reqBody, &logoutRequest)

	message := ""
	foundError := false

	if CurrentUserID != logoutRequest.UID {
		message = "Id doesn't match id to eb logged out."
		foundError = true 
	}

	fmt.Printf("Removing logged in userId: %v \n", CurrentUserID)

	CurrentUserID = 0

	resMessage := Response{
		message,
		foundError,
	}

	json.NewEncoder(res).Encode(resMessage)
}

func searchUsersForUsername(userName string) User {
	user := User{}

	for _, userFromList := range ListOfUsers {
		if userFromList.Username == strings.TrimSpace(userName) {
			user = userFromList
			break
		}
	}

	return user
}

func searchUsersByID(ID int) User {
	user := User{}

	for _, userFromList := range ListOfUsers {
		if userFromList.ID == ID {
			user = userFromList
			break
		}
	}

	return user
}

// TODO: finish this
// func searchUsersByProp(property string, value interface{}) User {
// 	user := User{}
// 	for i := range ListOfUsers {
// 		rv := reflect.ValueOf(ListOfUsers[i])

// 		rv = rv.Elem()

// 		field := rv.FieldByName(property)

// 		if !field.IsValid() {
// 			fmt.Errorf("not a field name: %s", property)
// 		}

// 		if field == value {
// 			user = ListOfUsers[i]
// 			break
// 		}
// 	}

// 	return user
// }

func addNewUser(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)
	var newPlayer User
	json.Unmarshal(reqBody, &newPlayer)
	ListOfUsers = append(ListOfUsers, newPlayer)
	fmt.Println("Player with username: " + newPlayer.Username + " added!")
	json.NewEncoder(res).Encode(newPlayer)
}

func showPlayers(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(ListOfUsers)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/new-player", addNewUser).Methods("POST")
	myRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/logout", logout).Methods("POST")
	myRouter.HandleFunc("/submit-answer", submitAnswersAndGetResults).Methods("POST")
	myRouter.HandleFunc("/compare-your-results", login).Methods("POST")
	myRouter.HandleFunc("/players", showPlayers)
	log.Fatal(http.ListenAndServe(":9990", myRouter))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ft-quiz.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ft-quiz" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ft-quiz")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
