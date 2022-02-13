package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"quiz-api-service/internal/api"
	"quiz-api-service/internal/services"
)

type UserHandler struct {
	UserService services.UserServices
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: service,
	}
}

// UserRoutes sets up user routes with accompanying methods for processing
func (handler *UserHandler) UserRoutes(r *gin.RouterGroup) {

	r.GET("", handler.getUsers).
		GET("username/:username", handler.getUserByUsername).
		GET("id/:uID", handler.getUserByID)

	return
}

func (handler *UserHandler) getUsers(c *gin.Context) {

	users, err := handler.UserService.GetUsers()
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, users)
	return
}

func (handler *UserHandler) getUserByUsername(c *gin.Context) {

	username := c.Param("username")

	if username == "" {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("missing username in url"))
		return
	}

	users, err := handler.UserService.GetUserByUsername(username)
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, users)
	return
}

func (handler *UserHandler) getUserByID(c *gin.Context) {

	userID := c.Param("uID")
	userIDint, err := strconv.Atoi(userID)

	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	} else if userIDint < 1 {
		_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	users, err := handler.UserService.GetUserByID(uint(userIDint))
	if err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, users)
	return
}

// HandleRequests sets up the server and the endpoints
// func HandleRequests() {
// 	myRouter := mux.NewRouter().StrictSlash(true)
// 	myRouter.HandleFunc("/", homePage)
// 	myRouter.HandleFunc("/new-player", addNewUser).Methods("POST")
// 	myRouter.HandleFunc("/login", login).Methods("POST")
// 	myRouter.HandleFunc("/logout", logout).Methods("POST")
// 	// myRouter.HandleFunc("/submit-answer", submitAnswersAndGetResults).Methods("POST")
// 	// myRouter.HandleFunc("/compare-your-score", compareUserScores).Methods("POST")
// 	myRouter.HandleFunc("/players", showPlayers)
// 	log.Fatal(http.ListenAndServe(":"+config.CurrentConfigs.Port, myRouter))
// }

// One of the endpoints that shows the homepage
func homePage(res http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(res, "Welcome to the HomePage! Go to the console to start playing.")
}

// Login endpoint function that checks username and password and sets user appropriately
// func login(res http.ResponseWriter, req *http.Request) {
// 	reqBody, _ := ioutil.ReadAll(req.Body)
//
// 	var loginReq api.LoginRequest
// 	_ = json.Unmarshal(reqBody, &loginReq)
//
// 	currentUser, err := UserService.GetUserByUsername(loginReq.Username)
// 	if err != nil {
// 		messageResponse := api.Response{
// 			Message: err.Error(),
// 			Error:   true,
// 		}
//
// 		_ = json.NewEncoder(res).Encode(messageResponse)
// 		return
// 	}
//
// 	foundError := false
// 	message := ""
// 	if currentUser.ID == 0 {
// 		message = "User not found: " + loginReq.Username
// 	} else if currentUser.Password != strings.TrimSpace(loginReq.Password) {
// 		message = "Password is wrong try again."
// 	} else {
// 		foundError = true
// 		message = "Successfully logged in with user: " + currentUser.Username
// 	}
//
// 	messageResponse := api.Response{
// 		Message: message,
// 		Error:   foundError,
// 	}
//
// 	_ = json.NewEncoder(res).Encode(messageResponse)
//
// }

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

// This is the gameplay section, loops through each question, outputting the questions and possible answers
// Will wait for user input and go to the next question, once they're all answered it posts to get the result
// func play(reader *bufio.Reader) bool {
//
// 	key, err := reader.ReadString('\n')
// 	if err != nil {
// 		return false
// 	}
//
// 	if len(key) > 0 {
// 		var listOfSubmittedRunes []rune
// 		// for _, q := range ListOfQuestions {
// 		// 	fmt.Printf("Question: %d %s", q.ID, q.Description)
// 		// 	questionKeyStrings := make([]string, 0, len(q.AnswerSelection))
// 		// 	for k := range q.AnswerSelection {
// 		// 		questionKeyStrings = append(questionKeyStrings, string(k))
// 		// 	}
// 		// 	sort.Strings(questionKeyStrings)
// 		// 	for {
// 		// 		for _, option := range questionKeyStrings {
// 		// 			runeKey := []rune(option)
// 		// 			fmt.Printf("%v: %v \n", option, q.AnswerSelection[runeKey[0]])
// 		// 		}
// 		// 		submittedAnswer, err := reader.ReadString('\n')
// 		// 		check(err)
// 		// 		if !isAnswerValid(questionKeyStrings, strings.TrimSpace(submittedAnswer)) {
// 		// 			fmt.Println("Please enter only one of the following options: ")
// 		// 			continue
// 		// 		}
// 		// 		answerRune := []rune(submittedAnswer)[0]
// 		// 		fmt.Println("answer submitted: " + string(answerRune))
// 		// 		listOfSubmittedRunes = append(listOfSubmittedRunes, answerRune)
// 		// 		break
// 		// 	}
// 		// }
// 		fmt.Println("Evaluating answers...")
// 		time.Sleep(2 * time.Second) // for dramatic suspense
//
// 		submitAnswerRequest := api.SubmitAnswersRequest{
// 			// UserID:           currentUser.ID,
// 			SubmittedAnswers: listOfSubmittedRunes,
// 		}
//
// 		PostToEndpoint(submitAnswerRequest, "submit-answer")
//
// 	}
//
// 	return true
// }

// Console input and logic to set user to post to endpoint
// func createUser(reader *bufio.Reader) bool {
// 	newUser := pkg.User{}
//
// 	fmt.Println("Start creating your profile.")
//
// 	for {
// 		fmt.Println("Enter Full Name: ")
// 		newUser.Name, _ = reader.ReadString('\n')
//
// 		fmt.Println("Enter your age: ")
// 		rawAge, _ := reader.ReadString('\n')
// 		rawAge = strings.TrimSpace(rawAge)
// 		intAge, _ := strconv.Atoi(rawAge)
// 		newUser.Age = int8(intAge)
//
// 		fmt.Println("Enter a Username: ")
// 		newUser.Username, _ = reader.ReadString('\n')
// 		fmt.Println("Enter a password: ")
// 		newUser.Password, _ = reader.ReadString('\n')
//
// 		newUser.Name = strings.TrimSpace(newUser.Name)
// 		newUser.Username = strings.TrimSpace(newUser.Username)
// 		newUser.Password = strings.TrimSpace(newUser.Password)
//
// 		if PostToEndpoint(newUser, "new-player") {
// 			break
// 		}
//
// 	}
//
// 	return true
// }

// This will check the users input and set the CurrentUser
// func loginPrompt(reader *bufio.Reader) bool {
// 	fmt.Println("Please enter Username:")
// 	inputtedUsername, _ := reader.ReadString('\n')
//
// 	fmt.Println("Please enter Password:")
// 	inputtedPassword, _ := reader.ReadString('\n')
//
// 	loginRequest := api.LoginRequest{
// 		Username: inputtedUsername,
// 		Password: inputtedPassword,
// 	}
//
// 	return PostToEndpoint(loginRequest, "login")
// }

// Console function to post to log-out endpoint by taking the CurrentUserID
// func logoutPrompt(uID uint) bool {
// 	logoutRequest := api.LogoutRequest{
// 		UserID: int(uID),
// 	}
//
// 	return PostToEndpoint(logoutRequest, "logout")
// }

// This is one of the endpoint functions that stores the users submitted answers, then
// Sets the users score and response with a Response struct parsed into json
func submitAnswersAndGetResults(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var submitAnswerRequest api.SubmitAnswersRequest
	_ = json.Unmarshal(reqBody, &submitAnswerRequest)

	// currentUser := searchUsersByID(submitAnswerRequest.UserID)
	// currentUser.SubmittedAnswers = submitAnswerRequest.SubmittedAnswers
	// currentUser.Score = 0

	// for i := range ListOfQuestions {
	// 	if ListOfQuestions[i].CorrectAnswer == currentUser.SubmittedAnswers[i] {
	// 		currentUser.Score++
	// 	}
	// }

	// updateUser(currentUser)

	// message := "You have answered " + strconv.Itoa(currentUser.Score) + " out of " + strconv.Itoa(len(ListOfQuestions)) + " questions correctly!"

	responseMessage := api.Response{
		// Message: message,
	}

	_ = json.NewEncoder(res).Encode(responseMessage)
}
