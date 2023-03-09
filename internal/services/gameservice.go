package services

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"quiz-api-service/internal/api"
)

type GameService struct {
	DBConn   *gorm.DB
	logger   *zap.Logger
	settings GameServiceSettings
}

type GameServiceSettings struct {
	Port     int
	Hostname string
}

type GameServices interface {
	StartGame(req *api.StartGameRequest) (*api.Game, error)
	// SubmitAnswer(answer *api.Answer) (err error)
	// FinishGame(req *api.FinishGameRequest) (err error)
}

func NewGameService(dbConn *gorm.DB, logger *zap.Logger, settings GameServiceSettings) *GameService {
	return &GameService{
		DBConn:   dbConn,
		logger:   logger,
		settings: settings,
	}
}

func (g GameService) StartGame(req *api.StartGameRequest) (*api.Game, error) {
	return nil, nil
}
