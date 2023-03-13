package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"quiz-api-service/internal/api"
	"quiz-api-service/internal/services"
)

type UserHandler struct {
	UserService services.UserServices
	Validator   *validator.Validate
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: service,
		Validator:   validator.New(),
	}
}

// UserRoutes sets up user routes with accompanying methods for processing
func (handler *UserHandler) UserRoutes(r *gin.RouterGroup) {

	r.POST("login", handler.login)

	r.Group("users").
		GET("", handler.getUsers).
		GET("username/:username", handler.getUserByUsername).
		GET("id/:uID", handler.getUserByID).
		POST("new", handler.newUser)

	return
}

func (handler *UserHandler) getUsers(c *gin.Context) {

	users, err := handler.UserService.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to get users", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully grabbed all users", users, nil))
	return
}

func (handler *UserHandler) getUserByUsername(c *gin.Context) {

	username := c.Param("username")

	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to get username from url", nil, fmt.Errorf("missing url")))
		return
	}

	user, err := handler.UserService.GetUserByUsername(username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to get user by username", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully got user", user, nil))
	return
}

func (handler *UserHandler) getUserByID(c *gin.Context) {

	userID := c.Param("uID")
	userIDint, err := strconv.Atoi(userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to get userID", nil, err))
		return
	} else if userIDint < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to get user", nil, fmt.Errorf("invalid user id")))
		return
	}

	user, err := handler.UserService.GetUserByID(uint(userIDint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to get user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully got user", user, nil))
	return
}

func (handler *UserHandler) newUser(c *gin.Context) {

	var user api.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to parse new user request", nil, err))
		return
	}

	if err := handler.Validator.Struct(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("missing or incorrect data received", nil, err))
		return
	}

	res, err := handler.UserService.InsertUser(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to add user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully inserted user", res, nil))
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

// Login endpoint function that checks username and password and sets user appropriately
func (handler *UserHandler) login(c *gin.Context) {
	var loginReq api.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to login", nil, err))
		return
	}

	user, err := handler.UserService.Login(loginReq)
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, api.GenerateMessageResponse("failed to login requested user", nil, err))
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to login requested user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("login successful", user, nil))
	return
}

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

	_ = json.NewEncoder(res).Encode(api.GenerateMessageResponse("", nil, nil))
}
