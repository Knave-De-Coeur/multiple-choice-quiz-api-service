package pkg

import (
	"database/sql"

	"gorm.io/gorm"
)

// User is generic user that launches and signs up
type User struct {
	gorm.Model
	Name               string       `json:"name"`
	Age                int8         `json:"age"`
	Username           string       `json:"username"`
	Password           string       `json:"password"`
	LastLoginTimeStamp sql.NullTime `json:"-"`
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

// Game is a single game played with all the stats
type Game struct {
	gorm.Model
	Questions        []Question
	User             User   `gorm:"embedded"`
	SubmittedAnswers []rune `gorm:"-" json:"submitted_answers"`
	Score            int    `json:"score"`
}
