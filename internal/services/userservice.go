package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"quiz-api-service/internal/api"
	"quiz-api-service/internal/pkg"
)

type UserService struct {
	DBConn   *gorm.DB
	Nats     *nats.Conn
	logger   *zap.Logger
	settings UserServiceSettings
}

// UserServiceSettings used to affect code flow
type UserServiceSettings struct {
	Port     int
	Hostname string
}

type UserServices interface {
	InsertUser(user *api.User) (*api.User, error)
	GetUsers() ([]api.User, error)
	GetUserByUsername(username string) (*pkg.User, error)
	GetUserByID(uID uint) (*api.User, error)
	Login(request api.LoginRequest) (*api.User, error)
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
func (service *UserService) InsertUser(req *api.User) (*api.User, error) {

	users, err := service.GetUsers()
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

	user := &pkg.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Age:       req.Age,
		Password:  gpResponse.Password,
	}

	res := service.DBConn.
		Select("first_name", "last_name", "email", "age", "username", "password").
		Create(user)
	if res.Error != nil {
		service.logger.Error("something went wrong inserting user", zap.Any("user", user), zap.Error(res.Error))
		return nil, res.Error
	}

	service.logger.Debug("rows inserted", zap.Int64("rowsAffected", res.RowsAffected))

	return req, nil
}

// GetUsers returns list of users in db.
func (service *UserService) GetUsers() ([]api.User, error) {

	var users []api.User

	// Get all records
	res := service.DBConn.
		Select("first_name", "last_name", "email", "age", "username", "created_at", "updated_at", "id").
		Find(&users)
	if res.Error != nil {
		service.logger.Error("something went wrong getting all players", zap.Error(res.Error))
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
		service.logger.Error("something went wrong getting player by username", zap.Error(res.Error), zap.String("username", username))
		return nil, res.Error
	}

	service.logger.Debug("user grabbed", zap.Any("user", user))

	return &user, nil
}

// GetUserByID grabs from table by id
func (service *UserService) GetUserByID(uID uint) (*api.User, error) {

	var user api.User
	// Get all records
	res := service.DBConn.
		Select("first_name", "last_name", "email", "age", "username", "password", "last_login_time_stamp").
		Where("id = ?", uID).
		First(&user)
	if res.Error != nil {
		service.logger.Error("something went wrong getting player by ID", zap.Error(res.Error))
		return nil, res.Error
	}

	service.logger.Debug("user grabbed", zap.Any("user", user))

	return &user, nil
}

// Login is a wrapper for the GetUserByUsername that also validates the password
func (service *UserService) Login(request api.LoginRequest) (*api.User, error) {

	user, err := service.GetUserByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	if user.Password != request.Password {
		return nil, fmt.Errorf("invalid passord for user")
	}

	unixCT := service.DBConn.NowFunc()

	fieldsToUpdate := map[string]interface{}{"last_login_time_stamp": unixCT, "updated_at": unixCT}

	// update record with login timestamp
	res := service.DBConn.
		Table("users").
		Where("id = ?", user.ID).
		Updates(fieldsToUpdate)
	if res.Error != nil {
		service.logger.Error("something went wrong updating a player", zap.Error(res.Error))
		return nil, res.Error
	}

	return &api.User{
		ID:                 strconv.Itoa(int(user.ID)),
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Email:              user.Email,
		Username:           user.Username,
		Age:                user.Age,
		CreatedAT:          user.CreatedAt.Format(time.RFC3339),
		UpdatedAT:          user.UpdatedAt.Format(time.RFC3339),
		LastLoginTimeStamp: user.LastLoginTimeStamp.Time.Format(time.RFC3339),
	}, nil
}

// Compare stats endpoint func that returns the message with how the user did compare to others
// func compareUserScores(res http.ResponseWriter, req *http.Request) {
// 	reqBody, _ := ioutil.ReadAll(req.Body)
// 	var compareUsersRequest api.CompareUsersRequest
// 	_ = json.Unmarshal(reqBody, &compareUsersRequest)
//
// 	user := searchUsersByID(compareUsersRequest.UserID)
//
// 	message := ""
// 	errorFound := false
// 	if len(user.SubmittedAnswers) < 1 {
// 		message = "Start playing to compare results!"
// 		errorFound = true
// 	} else {
// 		x := getUserComparisonScore(user)
//
// 		negative := math.Signbit(x)
//
// 		userScoreComparison := strconv.FormatFloat(x, 'f', 0, 64)
//
// 		if negative {
// 			message = "You did " + userScoreComparison + "% worse than everyone!"
// 		} else {
// 			message = "You did " + userScoreComparison + "% better than everyone!"
// 		}
// 	}
//
// 	responseMessage := api.Response{
// 		Message: message,
// 		Error:   errorFound,
// 	}
//
// 	_ = json.NewEncoder(res).Encode(responseMessage)
// }

// This calculates the comparison percentage the user has from other users
// func getUserComparisonScore(currentUser pkg.User) float64 {
//
// 	var listOfScores []int
//
// 	var sumPercentages int
//
// 	for i := range ListOfUsers {
// 		if ListOfUsers[i].ID != currentUser.ID {
// 			scorePercentage := ListOfUsers[i].Score * 20
// 			sumPercentages += scorePercentage
// 			listOfScores = append(listOfScores, scorePercentage)
// 		}
// 	}
//
// 	averagePercentage := float64(sumPercentages) / (float64(len(listOfScores)))
//
// 	scorePercentage := float64(currentUser.Score * 20)
//
// 	x := scorePercentage - averagePercentage
//
// 	return x
// }
