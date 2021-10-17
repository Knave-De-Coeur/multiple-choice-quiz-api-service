package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"quiz-api-service/internal/api"
	"quiz-api-service/internal/config"
	"quiz-api-service/internal/pkg"
)

// ListOfQuestions will hold a list of Question structs to be displayed and evaluated
var ListOfQuestions []pkg.Question

// ListOfUsers all the players that have registered to the quiz
var ListOfUsers []pkg.User

// CurrentUserID The logged-in user id interacting with the application
var CurrentUserID int

// Port is the port number the server will run on, defined as an arg in the app launch
var Port string

// FullHostname is the host and port concatenated
var FullHostname string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multiple-choice-quiz-api-service",
	Short: "My test for Fast Track",
	Long: `This is simple quiz where the user is presses ted with a couple questions
			and they have to select one from three to get the right answer.`,
	Args: func(cmd *cobra.Command, args []string) error {
		Port = ""
		if len(args) < 1 {
			Port = config.DefaultPort
		} else {
			Port = args[0]
		}
		checkAndAssignPort(Port)
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		// sets the questions up
		generateQuestions()
		// sets users
		generateDummyUsers()

		wg := new(sync.WaitGroup)

		wg.Add(2)

		// creates server and sets endpoints
		go handleRequests()

		// runs game
		go runGame()

		wg.Wait()
	},
}

// handleRequests sets up the server and the endpoints
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/new-player", addNewUser).Methods("POST")
	myRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/logout", logout).Methods("POST")
	myRouter.HandleFunc("/submit-answer", submitAnswersAndGetResults).Methods("POST")
	myRouter.HandleFunc("/compare-your-score", compareUserScores).Methods("POST")
	myRouter.HandleFunc("/players", showPlayers)
	log.Fatal(http.ListenAndServe(":"+Port, myRouter))
}

// This is one of the main goroutines of the application that runs the user interface part
func runGame() {
	fmt.Println("Welcome to Alex's quiz! Press enter to begin.")

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
				case 'b':
					loginPrompt(reader)
				case 'c':
					os.Exit(0)
				default:
					fmt.Println("Invalid try again")
				}
			} else {
				fmt.Println("a: Play")
				fmt.Println("b: Logout")
				fmt.Println("c: Compare")
				fmt.Println("d: Exit")

				optionStr, _ := reader.ReadString('\n')

				option := []rune(optionStr)

				switch option[0] {
				case 'a':
					play(reader)
				case 'b':
					logoutPrompt()
				case 'c':
					compare()
				case 'd':
					os.Exit(0)
				default:
					fmt.Println("Invalid try again")
				}
			}

		}
	}
}

// This populates the ListOfQuestions with dummy data
func generateQuestions() {
	ListOfQuestions = []pkg.Question{
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

// This simply populates the ListOfUsers with dummy data
func generateDummyUsers() {
	ListOfUsers = []pkg.User{
		{
			ID:       1,
			Name:     "David Smith",
			Age:      54,
			Username: "david54",
			Password: "pass765",
			SubmittedAnswers: []rune{
				'a', 'c', 'b', 'b', 'c',
			},
			Score: 2,
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
			Score: 3,
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
			Score: 2,
		},
	}
}

// CONSOLE FUNCTIONS

// This is the gameplay section, loops through each question, outputting the questions and possible answers
// Will wait for user input and go to the next question, once they're all answered it posts to get the result
func play(reader *bufio.Reader) bool {
	fmt.Printf("Press any key followed by enter to start the game you have %v questions \n", len(ListOfQuestions))

	currentUser := searchUsersByID(CurrentUserID)

	key, err := reader.ReadString('\n')
	check(err)

	if len(key) > 0 {
		var listOfSubmittedRunes []rune
		for _, q := range ListOfQuestions {
			fmt.Printf("Question: %d %s", q.ID, q.Description)
			questionKeyStrings := make([]string, 0, len(q.AnswerSelection))
			for k := range q.AnswerSelection {
				questionKeyStrings = append(questionKeyStrings, string(k))
			}
			sort.Strings(questionKeyStrings)
			for {
				for _, option := range questionKeyStrings {
					runeKey := []rune(option)
					fmt.Printf("%v: %v \n", option, q.AnswerSelection[runeKey[0]])
				}
				submittedAnswer, err := reader.ReadString('\n')
				check(err)
				if !isAnswerValid(questionKeyStrings, strings.TrimSpace(submittedAnswer)) {
					fmt.Println("Please enter only one of the following options: ")
					continue
				}
				answerRune := []rune(submittedAnswer)[0]
				fmt.Println("answer submitted: " + string(answerRune))
				listOfSubmittedRunes = append(listOfSubmittedRunes, answerRune)
				break
			}
		}
		fmt.Println("Evaluating answers...")
		time.Sleep(2 * time.Second) // for dramatic suspense

		submitAnswerRequest := api.SubmitAnswersRequest{
			UserID:           currentUser.ID,
			SubmittedAnswers: listOfSubmittedRunes,
		}

		postToEndpoint(submitAnswerRequest, "submit-answer")

	}

	return true
}

// Console input and logic to set user to post to endpoint
func createUser(reader *bufio.Reader) bool {
	uid := len(ListOfUsers) + 1
	newUser := pkg.User{
		ID: uid,
	}

	fmt.Println("Start creating your profile.")

	for {
		fmt.Println("Enter Full Name: ")
		newUser.Name, _ = reader.ReadString('\n')

		fmt.Println("Enter your age: ")
		rawAge, _ := reader.ReadString('\n')
		rawAge = strings.TrimSpace(rawAge)
		intAge, _ := strconv.Atoi(rawAge)
		newUser.Age = int8(intAge)

		fmt.Println("Enter a Username: ")
		newUser.Username, _ = reader.ReadString('\n')
		fmt.Println("Enter a password: ")
		newUser.Password, _ = reader.ReadString('\n')

		newUser.Name = strings.TrimSpace(newUser.Name)
		newUser.Username = strings.TrimSpace(newUser.Username)
		newUser.Password = strings.TrimSpace(newUser.Password)

		if postToEndpoint(newUser, "new-player") {
			break
		}

	}

	return true
}

// This will check the users input and set the CurrentUser
func loginPrompt(reader *bufio.Reader) bool {
	fmt.Println("Please enter Username:")
	inputtedUsername, _ := reader.ReadString('\n')

	fmt.Println("Please enter Password:")
	inputtedPassword, _ := reader.ReadString('\n')

	loginRequest := api.LoginRequest{
		Username: inputtedUsername,
		Password: inputtedPassword,
	}

	return postToEndpoint(loginRequest, "login")
}

// Console function to post to log-out endpoint by taking the CurrentUserID
func logoutPrompt() bool {
	logoutRequest := api.LogoutRequest{
		UserID: CurrentUserID,
	}

	return postToEndpoint(logoutRequest, "logout")
}

// Compare console func that simply posts to the endpoint and displays the message
func compare() bool {
	currentUser := searchUsersByID(CurrentUserID)

	requestData := api.CompareUsersRequest{
		UserID:    currentUser.ID,
		UserScore: currentUser.Score,
	}

	return postToEndpoint(requestData, "compare-your-score")
}

// Grouped logic that posts to endpoint and receives message to be outputted to the console
func postToEndpoint(data interface{}, endpoint string) bool {
	requestJSON, _ := json.Marshal(data)

	res, err := http.Post(FullHostname+endpoint, "application/json", bytes.NewBuffer(requestJSON))
	check(err)

	defer res.Body.Close()

	var responseMessage api.Response

	_ = json.NewDecoder(res.Body).Decode(&responseMessage)

	fmt.Println(responseMessage.Message)

	return !responseMessage.Error
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

	currentUser := searchUsersForUsername(loginReq.Username)

	foundError := false
	message := ""
	if currentUser.ID == 0 {
		message = "User not found: " + loginReq.Username
	} else if currentUser.Password != strings.TrimSpace(loginReq.Password) {
		message = "Password is wrong try again."
	} else {
		CurrentUserID = currentUser.ID
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

	if CurrentUserID != logoutRequest.UserID {
		message = "Id doesn't match id to eb logged out."
		foundError = true
	}

	CurrentUserID = 0

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
	for _, u := range ListOfUsers {
		if newPlayer.Username == u.Username {
			message = "Username " + u.Username + " is taken please try again."
			errorFound = true
		} else if newPlayer.Password == u.Password {
			message = "Please provide a different password"
			errorFound = true
		}
	}

	newPlayer.SubmittedAnswers = []rune{}

	if !errorFound {
		ListOfUsers = append(ListOfUsers, newPlayer)
		message = "User by username: " + newPlayer.Username + " has been successfully added!"
	}

	responseMessage := api.Response{
		Message: message,
		Error:   errorFound,
	}

	_ = json.NewEncoder(res).Encode(responseMessage)
}

// GET rest endpoint func that simply displays list of users in json
func showPlayers(res http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(res).Encode(ListOfUsers)
}

// This is one of the endpoint functions that stores the users submitted answers, then
// Sets the users score and response with a Response struct parsed into json
func submitAnswersAndGetResults(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var submitAnswerRequest api.SubmitAnswersRequest
	_ = json.Unmarshal(reqBody, &submitAnswerRequest)

	currentUser := searchUsersByID(submitAnswerRequest.UserID)
	currentUser.SubmittedAnswers = submitAnswerRequest.SubmittedAnswers
	currentUser.Score = 0

	for i := range ListOfQuestions {
		if ListOfQuestions[i].CorrectAnswer == currentUser.SubmittedAnswers[i] {
			currentUser.Score++
		}
	}

	updateUser(currentUser)

	message := "You have answered " + strconv.Itoa(currentUser.Score) + " out of " + strconv.Itoa(len(ListOfQuestions)) + " questions correctly!"

	responseMessage := api.Response{
		Message: message,
	}

	_ = json.NewEncoder(res).Encode(responseMessage)
}

// This updates the user answers in memory
func updateUser(user pkg.User) {
	for i, u := range ListOfUsers {
		if u.ID == user.ID {
			ListOfUsers[i] = user
		}
	}
}

// Compare stats endpoint func that returns the message with how the user did compare to others
func compareUserScores(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)
	var compareUsersRequest api.CompareUsersRequest
	_ = json.Unmarshal(reqBody, &compareUsersRequest)

	user := searchUsersByID(compareUsersRequest.UserID)

	message := ""
	errorFound := false
	if len(user.SubmittedAnswers) < 1 {
		message = "Start playing to compare results!"
		errorFound = true
	} else {
		x := getUserComparisonScore(user)

		negative := math.Signbit(x)

		userScoreComparison := strconv.FormatFloat(x, 'f', 0, 64)

		if negative {
			message = "You did " + userScoreComparison + "% worse than everyone!"
		} else {
			message = "You did " + userScoreComparison + "% better than everyone!"
		}
	}

	responseMessage := api.Response{
		Message: message,
		Error:   errorFound,
	}

	_ = json.NewEncoder(res).Encode(responseMessage)
}

// HELPER FUNCTIONS

// Simply check that the answer the user inputted exits
func isAnswerValid(answers []string, submittedAnswer string) bool {
	for _, item := range answers {
		if item == submittedAnswer {
			return true
		}
	}
	return false
}

// This simply checks if the port is available and assigns it to the global variables
func checkAndAssignPort(port string) {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Printf("Can't listen on port %q: %s \n", port, err)
		os.Exit(1)
	}

	_ = ln.Close()

	Port = port
	FullHostname = config.Host + ":" + Port + "/"
}

// Simply checks if error exists and panics accordingly
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Goes through lists of users and returns the user with the correct username or nothing
func searchUsersForUsername(userName string) pkg.User {
	user := pkg.User{}

	for _, userFromList := range ListOfUsers {
		if userFromList.Username == strings.TrimSpace(userName) {
			user = userFromList
			break
		}
	}

	return user
}

// Goes through list of users and returns user with correct ID or nothing
func searchUsersByID(ID int) pkg.User {
	user := pkg.User{}

	for _, userFromList := range ListOfUsers {
		if userFromList.ID == ID {
			user = userFromList
			break
		}
	}

	return user
}

// This calculates the comparison percentage the user has from other users
func getUserComparisonScore(currentUser pkg.User) float64 {

	var listOfScores []int

	var sumPercentages int

	for i := range ListOfUsers {
		if ListOfUsers[i].ID != currentUser.ID {
			scorePercentage := ListOfUsers[i].Score * 20
			sumPercentages += scorePercentage
			listOfScores = append(listOfScores, scorePercentage)
		}
	}

	averagePercentage := float64(sumPercentages) / (float64(len(listOfScores)))

	scorePercentage := float64(currentUser.Score * 20)

	x := scorePercentage - averagePercentage

	return x
}

// COBRA FUNCS

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("something went wrong: %+v", err)
	}
}

func init() {
	cobra.OnInitialize()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.multiple-choice-quiz-api-service.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
