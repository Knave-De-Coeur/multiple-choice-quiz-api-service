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

// LogoutRequest is the parsed struct of the /logout endpoint
type LogoutRequest struct {
	UserID int `json:"user_id"`
}
