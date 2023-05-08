package api

type GeneratePasswordRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int8   `json:"age"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// LoginRequest is the parsed struct of the /login endpoint
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SubmitAnswersRequest is the struct used to decode and encode when using the /submit-answers endpoint
type SubmitAnswersRequest struct {
	UserID           int    `json:"user_id"`
	SubmittedAnswers []rune `json:"submitted_answers"`
}

// CompareUsersRequest is the struct we parse and send over to the /compare-your-score endpoint
type CompareUsersRequest struct {
	UserID    int `json:"user_id"`
	UserScore int `json:"user_score"`
}

// LogoutRequest is the parsed struct of the /logout endpoint
type LogoutRequest struct {
	UserID int `json:"user_id"`
}

type StartGameRequest struct {
	GameID int `json:"game_id" validate:"required"`
	UserID int `json:"user_id" validate:"required"`
}

type FinishGameRequest struct {
	GameID int `json:"game_id" validate:"required"`
	UserID int `json:"user_id" validate:"required"`
}
