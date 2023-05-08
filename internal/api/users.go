package api

type User struct {
	ID                 string `json:"ID,omitempty"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	Age                int8   `json:"age" validate:"required"`
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password,omitempty" validate:"required"`
	CreatedAT          string `json:"created_at,omitempty"`
	UpdatedAT          string `json:"updated_at,omitempty"`
	LastLoginTimeStamp string `json:"last_login_time_stamp,omitempty"`
}

type NewUserRequest struct {
	*User
}

type UpdateUserRequest struct {
	ID          uint   `json:"ID" validate:"gt=0"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Age         int8   `json:"age" validate:"required"`
	Username    string `json:"username" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type DeleteUserRequest struct {
	ID uint `json:"ID" validate:"gt=0"`
}
