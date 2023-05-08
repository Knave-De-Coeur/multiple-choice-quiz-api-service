package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"

	"quiz-api-service/internal/api"
	"quiz-api-service/internal/services"
)

type UserHandler struct {
	Nats        *nats.Conn
	UserService services.UserServices
	Validator   *validator.Validate
}

func NewUserHandler(service *services.UserService, nc *nats.Conn) *UserHandler {
	return &UserHandler{
		Nats:        nc,
		UserService: service,
		Validator:   validator.New(),
	}
}

// UserRoutes sets up user routes with accompanying methods for processing
func (handler *UserHandler) UserRoutes(r *gin.RouterGroup) {

	r.POST("login", handler.login)

	r.Group("user").
		GET("", handler.getUsers).
		GET("id/:uID", handler.getUserByID).
		POST("new", handler.newUser).
		PUT("/:uID", handler.updateUser)

}

func (handler *UserHandler) getUsers(c *gin.Context) {

	users, err := handler.UserService.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to get users", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully grabbed all users", users, nil))

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

}

func (handler *UserHandler) newUser(c *gin.Context) {

	var newUserReq api.NewUserRequest

	if err := c.ShouldBindJSON(&newUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to parse new user request", nil, err))
		return
	}

	if err := handler.Validator.Struct(newUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("missing or incorrect data received", nil, err))
		return
	}

	res, err := handler.UserService.InsertUser(newUserReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to add user", nil, err))
		return
	}

	c.JSON(http.StatusCreated, api.GenerateMessageResponse("successfully inserted user", res, nil))

}
func (handler *UserHandler) updateUser(c *gin.Context) {

	userID := c.Param("uID")
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("wrong id format in url", nil, err))
		return
	}

	var updateUserReq api.UpdateUserRequest

	if err = c.ShouldBindJSON(&updateUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to parse new user request", nil, err))
		return
	}

	updateUserReq.ID = uint(userIDint)

	if err = handler.Validator.Struct(updateUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("missing or incorrect data received", nil, err))
		return
	}

	err = handler.UserService.UpdateUser(updateUserReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to add user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully updated user", nil, nil))

}

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

}
