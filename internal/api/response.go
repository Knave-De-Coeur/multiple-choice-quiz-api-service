package api

type GeneratePasswordResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password" validate:"required"`
}

// MessageResponse is a generic response struct that'll be marshalled to json and sent to the requester
type MessageResponse struct {
	Message string `json:"message"`
	Result  any    `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}

func GenerateMessageResponse(message string, res interface{}, err error) *MessageResponse {

	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}

	return &MessageResponse{
		Message: message,
		Result:  res,
		Error:   errorMessage,
	}
}
