package services

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"quiz-api-service/internal/config"
	"quiz-api-service/internal/utils"
)

var (
	mockDB  *sql.DB
	sqlMock sqlmock.Sqlmock

	gormDB *gorm.DB
	log    *zap.Logger

	err error

	userService *UserService
)

func TestMain(m *testing.M) {
	log, err = utils.SetUpLogger()
	if err != nil {
		panic(err)
	}

	defer log.Sync()

	mockDB, sqlMock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		panic(fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	}
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	})

	gormDB, err = gorm.Open(dialector, &gorm.Config{
		NowFunc: func() time.Time {
			return time.Unix(time.Now().Unix(), 0).UTC()
		},
	})
	if err != nil {
		panic(err)
	}

	userService = NewUserService(gormDB, nil, log, UserServiceSettings{
		Port:     0,
		Hostname: config.CurrentConfigs.Host,
	})

	m.Run()
}
