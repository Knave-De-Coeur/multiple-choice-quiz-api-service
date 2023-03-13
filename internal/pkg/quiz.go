package pkg

import (
	"database/sql"

	"gorm.io/gorm"
)

// User is generic user that launches and signs up
type User struct {
	gorm.Model
	FirstName          string       `json:"first_name"`
	LastName           string       `json:"last_name"`
	Email              string       `json:"email"`
	Age                int8         `json:"age"`
	Username           string       `json:"username"`
	Password           string       `json:"password"`
	LastLoginTimeStamp sql.NullTime `json:"-"`
}

// Game is a single game played with all the stats
type Game struct {
	gorm.Model
	Questions        []Question
	User             User   `gorm:"-"`
	SubmittedAnswers []rune `gorm:"-" json:"submitted_answers"`
	Score            int    `json:"score"`
}

// Question presented and evaluated
type Question struct {
	gorm.Model
	GameID      uint
	Description string // the literal question to be displayed
	Answers     []Answer
}

// Answer is the answer to a specific question
type Answer struct {
	gorm.Model
	QuestionID uint
	Text       string `gorm:"longtext"`
	IsCorrect  bool
}

// UserGames is user_games table in quiz db.
type UserGames struct {
	UserID uint
	GameID uint
}

// UserAnswers is the user_answers table in quiz db.
type UserAnswers struct {
	UserID   uint
	AnswerID uint
}
