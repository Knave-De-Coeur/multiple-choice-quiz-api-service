package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/knave-de-coeur/user-api-service/internal/api"
	"github.com/knave-de-coeur/user-api-service/internal/middleware"
	"github.com/knave-de-coeur/user-api-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type IUserHandler interface {
	SetUpRoutes(r *gin.RouterGroup)
	getUsers(c *gin.Context)
	newUser(c *gin.Context)
	updateUser(c *gin.Context)
	deleteUser(c *gin.Context)
	login(c *gin.Context)
}

type UserHandler struct {
	UserService services.IUserService
	Middleware  middleware.IAuthMiddleware
	Validator   *validator.Validate
	RedisClient *redis.Client
	Nats        *nats.Conn
}

func NewUserHandler(service services.IUserService, auth middleware.IAuthMiddleware, redisClient *redis.Client,
	nc *nats.Conn) *UserHandler {

	return &UserHandler{
		Nats:        nc,
		UserService: service,
		Validator:   validator.New(),
		Middleware:  auth,
		RedisClient: redisClient,
	}
}

// SetUpRoutes sets up user routes with accompanying methods for processing
func (h *UserHandler) SetUpRoutes(r *gin.RouterGroup) {

	r.POST("login", h.login)
	// TODO: finalize logout
	//r.POST("logout", h.logout)

	r.GET("users", h.getUsers)

	r.Group("user").
		POST("", h.newUser).
		GET("/:uID", h.Middleware.RequireAuth(), h.getUserByID).
		PUT("/:uID", h.Middleware.RequireAuth(), h.updateUser).
		DELETE("/:uID", h.Middleware.RequireAuth(), h.deleteUser)

}

func (h *UserHandler) getUsers(c *gin.Context) {

	users, err := h.UserService.GetBasicUserDataList()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to get users", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully grabbed all users", users, nil))

}

func (h *UserHandler) getUserByID(c *gin.Context) {

	userID := c.Param("uID")
	userIDint, err := strconv.Atoi(userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to get userID", nil, err))
		return
	} else if userIDint < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to get user", nil, fmt.Errorf("invalid user id")))
		return
	}

	user, err := h.UserService.GetUserByID(uint(userIDint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to get user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully got user", user, nil))

}

func (h *UserHandler) newUser(c *gin.Context) {

	var newUserReq api.NewUserRequest

	if err := c.ShouldBindJSON(&newUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to parse new user request", nil, err))
		return
	}

	if err := h.Validator.Struct(newUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("missing or incorrect data received", nil, err))
		return
	}

	res, err := h.UserService.InsertUser(newUserReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to add user", nil, err))
		return
	}

	c.JSON(http.StatusCreated, api.GenerateMessageResponse("successfully inserted user", res, nil))

}

func (h *UserHandler) updateUser(c *gin.Context) {

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

	if err = h.Validator.Struct(updateUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("missing or incorrect data received", nil, err))
		return
	}

	err = h.UserService.UpdateUser(updateUserReq)
	if err != nil {
		var status int
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotModified
		} else {
			status = http.StatusInternalServerError
		}
		c.AbortWithStatusJSON(status, api.GenerateMessageResponse("failed to update user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully updated user", nil, nil))

}

func (h *UserHandler) deleteUser(c *gin.Context) {

	userID := c.Param("uID")
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("wrong id format in url", nil, err))
		return
	}

	var deleteUserReq api.DeleteUserRequest

	if err = c.ShouldBindJSON(&deleteUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to parse new user request", nil, err))
		return
	}

	deleteUserReq.ID = uint(userIDint)

	if err = h.Validator.Struct(deleteUserReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("missing or incorrect data received", nil, err))
		return
	}

	err = h.UserService.DeleteUser(deleteUserReq)
	if err != nil {
		var status int
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotModified
		} else {
			status = http.StatusInternalServerError
		}
		c.AbortWithStatusJSON(status, api.GenerateMessageResponse("failed to delete user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully deleted user", nil, nil))
}

// Login endpoint function that checks username and password and sets user appropriately
func (h *UserHandler) login(c *gin.Context) {
	var loginReq api.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to login", nil, err))
		return
	}

	user, err := h.UserService.Login(loginReq)
	if err != nil && err == gorm.ErrRecordNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, api.GenerateMessageResponse("failed to login requested user", nil, err))
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to login requested user", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("login successful", user, nil))

}
