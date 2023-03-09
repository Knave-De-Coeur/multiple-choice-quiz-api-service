package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"quiz-api-service/internal/api"
	"quiz-api-service/internal/services"
)

type GameHandler struct {
	GameService services.GameServices
	Validator   *validator.Validate
}

func NewGameHandler(service *services.GameService) *GameHandler {
	return &GameHandler{
		GameService: service,
		Validator:   validator.New(),
	}
}

// UserRoutes sets up user routes with accompanying methods for processing
func (handler *GameHandler) GameRoutes(r *gin.RouterGroup) {

	r.Group("game").
		POST("start", handler.startGame)
	// POST("submit", handler.newUser).
	// POST("finish", handler.newUser)

	return
}

func (handler *GameHandler) startGame(c *gin.Context) {

	var gameReq *api.StartGameRequest

	if err := c.ShouldBindJSON(&gameReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.GenerateMessageResponse("failed to start game", nil, err))
		return
	}

	game, err := handler.GameService.StartGame(gameReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.GenerateMessageResponse("failed to get users", nil, err))
		return
	}

	c.JSON(http.StatusOK, api.GenerateMessageResponse("successfully grabbed all users", game, nil))
	return
}
