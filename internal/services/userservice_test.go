package services

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"

	"github.com/knave-de-coeur/user-api-service/internal/api"
)

type genericTestCase struct {
	Name        string
	ExpectedErr bool
}

type InsertUserTest struct {
	genericTestCase
	Input          api.NewUserRequest
	ExpectedResult *api.User
	SqlMock        func(test InsertUserTest) bool
	SqlMockRows    *sqlmock.Rows
}

func TestUserService_InsertUser(t *testing.T) {
	testCases := []InsertUserTest{
		{
			genericTestCase: genericTestCase{
				Name:        "Simple test",
				ExpectedErr: false,
			},
			Input: api.NewUserRequest{
				User: &api.User{},
			},
			ExpectedResult: &api.User{},
			SqlMockRows:    sqlMock.NewRows([]string{"id"}).AddRow(1),
			SqlMock: func(test InsertUserTest) bool {
				sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT FROM users WHERE")).WillReturnRows()
				sqlMock.ExpectQuery(
					regexp.QuoteMeta("INSERT INTO users `first_name`,`last_name`, `email`, `age`, `username`, `password` VALUES (?, ?, ?, ?, ?)"),
				).WithArgs(test.Input.FirstName, test.Input.LastName, test.Input.Email, test.Input.Age, test.Input.Username, test.Input.Password).
					WillReturnRows(test.SqlMockRows)
				return true
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			// set up the mock results
			sqlM := test.SqlMock(test)

			res, err := userService.InsertUser(test.Input)
			if sqlM {
				if sqlMockErr := sqlMock.ExpectationsWereMet(); sqlMockErr != nil {
					t.Errorf("there were unfulfilled expectations: %s", sqlMockErr)
				}
			}
			if test.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, res, test.ExpectedResult)
			}
		})
	}

}
