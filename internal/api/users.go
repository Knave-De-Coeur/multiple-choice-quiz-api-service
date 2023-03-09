package api

type User struct {
	ID                 string `json:"ID"`
	Name               string `json:"name" validate:"required"`
	Age                int8   `json:"age" validate:"required"`
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password,omitempty" validate:"required"`
	CreatedAT          string `json:"created_at,omitempty"`
	UpdatedAT          string `json:"updated_at,omitempty"`
	LastLoginTimeStamp string `json:"last_login_time_stamp,omitempty"`
}
