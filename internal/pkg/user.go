package pkg

// User is generic user that launches and signs up
type User struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Age              int8   `json:"age"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	SubmittedAnswers []rune `json:"submitted_answers"`
	Score            int    `json:"score"`
}
