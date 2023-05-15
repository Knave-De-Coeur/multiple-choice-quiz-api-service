package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/knave-de-coeur/user-api-service/internal/api"
	"github.com/knave-de-coeur/user-api-service/internal/pkg"
	"github.com/knave-de-coeur/user-api-service/internal/utils"
)

type UserService struct {
	DBConn   *gorm.DB
	Nats     *nats.Conn
	logger   *zap.Logger
	settings UserServiceSettings
}

// UserServiceSettings used to affect code flow
type UserServiceSettings struct {
	Port      int
	Hostname  string
	JWTSecret string
}

type UserServices interface {
	InsertUser(user api.NewUserRequest) (*api.User, error)
	UpdateUser(req api.UpdateUserRequest) error
	DeleteUser(req api.DeleteUserRequest) error
	checkDuplicatePasswords(currentPass string) error
	GetBasicUserDataList() ([]api.User, error)
	GetUserByUsername(username string) (*pkg.User, error)
	GetUserByID(uID uint) (*api.User, error)
	Login(request api.LoginRequest) (*api.LoginResponse, error)
	getDBUserByID(uID uint) (*pkg.User, error)
}

func NewUserService(dbConn *gorm.DB, nc *nats.Conn, logger *zap.Logger, settings UserServiceSettings) *UserService {
	return &UserService{
		Nats:     nc,
		DBConn:   dbConn,
		logger:   logger,
		settings: settings,
	}
}

// InsertUser inserts new user in users table from data passed in arg.
func (service *UserService) InsertUser(req api.NewUserRequest) (*api.User, error) {

	users, err := service.GetBasicUserDataList()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == req.Email {
			return nil, errors.New("email already registered")
		}
		if user.Username == req.Username {
			return nil, errors.New("username taken")
		}
	}

	var encryptedPass string

	if service.Nats != nil {
		jsonGPR, err := json.Marshal(&api.GeneratePasswordRequest{
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Age:       req.Age,
			Email:     req.Email,
			Password:  req.Password,
		})
		if err != nil {
			service.logger.Error("something went wrong marshalling request", zap.Any("req", req), zap.Error(err))
			return nil, err
		}

		msg, err := service.Nats.Request(pkg.AuthGeneratePass, jsonGPR, 10*time.Second)
		if err != nil {
			service.logger.Error("couldn't get a response from auth service", zap.Any("req", jsonGPR), zap.Error(err))
			return nil, err
		}

		var gpResponse api.GeneratePasswordResponse

		if err = json.Unmarshal(msg.Data, &gpResponse); err != nil {
			service.logger.Error("something went wrong unmarshalling response from auth service", zap.Any("msg", msg), zap.Error(err))
			return nil, err
		}

		if gpResponse.Password == "" {
			err = errors.New("empty pass")
			service.logger.Error("missing pass from response", zap.Any("res", gpResponse), zap.Error(err))
			return nil, err
		}
	} else {
		encryptedPass, err = utils.HashAndSalt([]byte(req.Password))
		if err != nil {
			service.logger.Error("failed to encrypt pass", zap.Any("request", req), zap.Error(err))
			return nil, err
		}
	}

	user := &pkg.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Age:       req.Age,
		Password:  encryptedPass,
	}

	res := service.DBConn.
		Select("first_name", "last_name", "email", "age", "username", "password").
		Create(user)
	if res.Error != nil {
		service.logger.Error("something went wrong inserting user", zap.Any("user", user), zap.Error(res.Error))
		return nil, res.Error
	}

	service.logger.Debug("rows inserted", zap.Int64("rowsAffected", res.RowsAffected))

	req.User.ID = strconv.Itoa(int(user.ID))

	return req.User, nil
}

func (service *UserService) GetBasicUserDataList() ([]api.User, error) {

	var users []api.User

	// Get all records
	res := service.DBConn.
		Select("id", "first_name", "last_name", "email", "age", "username").
		Find(&users)
	if res.Error != nil {
		service.logger.Error("something went wrong getting all users", zap.Error(res.Error))
		return nil, res.Error
	}

	service.logger.Debug("users grabbed", zap.Int64("number", res.RowsAffected))

	return users, nil
}

// GetUserByUsername attempts to retrieve a single row from the users table.
func (service *UserService) GetUserByUsername(username string) (*pkg.User, error) {

	var user pkg.User
	// Get all records
	res := service.DBConn.
		Select("id", "first_name", "last_name", "email", "age", "username", "password", "created_at", "updated_at", "last_login_time_stamp").
		Where("username = ?", username).
		First(&user)
	if res.Error != nil {
		service.logger.Error("something went wrong getting user by username", zap.Error(res.Error), zap.String("username", username))
		return nil, res.Error
	}

	service.logger.Debug("user grabbed", zap.Any("user", user))

	return &user, nil
}

func (service *UserService) getDBUserByID(uID uint) (*pkg.User, error) {

	var user pkg.User
	// Get all records
	res := service.DBConn.
		Select("id", "first_name", "last_name", "email", "age", "username", "password", "created_at", "updated_at", "last_login_time_stamp").
		Where("id = ?", uID).
		First(&user)
	if res.Error != nil {
		service.logger.Error("something went wrong getting user by ID", zap.Error(res.Error))
		return nil, res.Error
	}

	service.logger.Debug("user grabbed", zap.Any("user", user))

	return &user, nil
}

// GetUserByID grabs from table by id
func (service *UserService) GetUserByID(uID uint) (*api.User, error) {

	user, err := service.getDBUserByID(uID)
	if err != nil {
		return nil, err
	}

	response := &api.User{
		ID:                 strconv.Itoa(int(user.ID)),
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Age:                user.Age,
		Email:              user.Email,
		Username:           user.Username,
		CreatedAT:          user.CreatedAt.Format(time.RFC3339),
		UpdatedAT:          user.UpdatedAt.Format(time.RFC3339),
		LastLoginTimeStamp: user.LastLoginTimeStamp.Time.Format(time.RFC3339),
	}

	return response, nil
}

// Login is a wrapper for the GetUserByUsername that also validates the password
func (service *UserService) Login(request api.LoginRequest) (*api.LoginResponse, error) {

	user, err := service.GetUserByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	isSame, err := utils.ComparePasswords(user.Password, []byte(request.Password))
	if err != nil {
		service.logger.Error("something went wrong comparing the passwords", zap.Error(err))
		return nil, err
	} else if !isSame {
		service.logger.Error("passwords don't match", zap.Any("req pass", request.Password))
		return nil, fmt.Errorf("passwords don't match")
	}

	// save userID in jwt token for requests
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(service.settings.JWTSecret))
	if err != nil {
		service.logger.Error("failed to create token", zap.Error(err))
		return nil, err
	}

	unixCT := service.DBConn.NowFunc()

	fieldsToUpdate := map[string]interface{}{"last_login_time_stamp": unixCT}

	// update record with login timestamp
	res := service.DBConn.
		Table("users").
		Where("id = ?", user.ID).
		Updates(fieldsToUpdate)
	if res.Error != nil {
		service.logger.Error("something went wrong updating a user", zap.Error(res.Error))
		return nil, res.Error
	}

	return &api.LoginResponse{
		Token: tokenString,
	}, nil
}

func (service *UserService) checkDuplicatePasswords(currentPass string) error {

	res := service.DBConn.Select("password").Where("password = ?", currentPass)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		service.logger.Error("something went wrong getting password", zap.Error(res.Error), zap.String("password", currentPass))
		return res.Error
	}

	return nil
}

func (service *UserService) UpdateUser(req api.UpdateUserRequest) error {

	user, err := service.getDBUserByID(req.ID)
	if err != nil {
		return err
	}

	fieldDataMap := map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"username":   req.Username,
		"email":      req.Email,
		"age":        req.Age,
	}

	if req.OldPassword != "" && req.NewPassword != "" {

		isSame, err := utils.ComparePasswords(user.Password, []byte(req.OldPassword))
		if err != nil {
			service.logger.Error("something went wrong comparing the passwords", zap.Error(err))
			return err
		} else if !isSame {
			service.logger.Error("passwords don't match", zap.Any("req pass", req.OldPassword))
			return fmt.Errorf("passwords don't match")
		}

		if err = service.checkDuplicatePasswords(req.NewPassword); err != nil {
			service.logger.Error("password is taken", zap.Error(err), zap.String("password", req.NewPassword))
			return err
		}

		encryptedPass, err := utils.HashAndSalt([]byte(req.NewPassword))
		if err != nil {
			service.logger.Error("something went wrong encrypting the new password", zap.Error(err))
			return err
		}

		fieldDataMap["password"] = encryptedPass
	}

	// update record with login timestamp
	res := service.DBConn.
		Table("users").
		Where("id = ?", user.ID).
		Updates(fieldDataMap)
	if res.Error != nil {
		service.logger.Error("something went wrong updating a user", zap.Error(res.Error))
		return res.Error
	}

	return nil
}

func (service *UserService) DeleteUser(req api.DeleteUserRequest) error {

	var tx *gorm.DB

	if req.HardDelete {
		tx = service.DBConn.Unscoped()
	} else {
		// updates record with deleted_at timestamp
		tx = service.DBConn
	}

	res := tx.Table("users").Where("id = ?", req.ID).Delete(&pkg.User{Model: gorm.Model{ID: req.ID}})

	if res.Error != nil {
		service.logger.Error("something went wrong deleting a user", zap.Error(res.Error))
		return res.Error
	}

	return nil
}
