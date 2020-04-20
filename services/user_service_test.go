package services

import (
	"github.com/leandrotula/golangmicroservice/domain"
	"github.com/leandrotula/golangmicroservice/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

var generateMockData func (id int64) (*domain.User, *util.ResponseError)

type mockDaoImpl struct {
}

func init() {
	domain.UserDao = &mockDaoImpl{}
}

func(m *mockDaoImpl) GetUser(userId int64)(*domain.User, *util.ResponseError) {
	return generateMockData(userId)
}

func TestGetUserNotFound(t *testing.T) {

	generateMockData = func(id int64) (*domain.User, *util.ResponseError) {
		return nil, &util.ResponseError{
			Message: "User not found",
			Code:    404,
		}
	}

	user, err := GetUser(11)
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "User not found", err.Message)
	assert.Equal(t, 404, err.Code)

}

func TestGetUserFound(t *testing.T) {

	generateMockData = func(id int64) (*domain.User, *util.ResponseError) {
		return &domain.User{
			Id:        uint64(2),
			FirstName: "User2",
			LastName:  "LastNameUser2",
			Email:     "test2@domain.com",
		}, nil
	}

	user, err := GetUser(2)
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), user.Id)
	assert.Equal(t, "User2", user.FirstName)
	assert.Equal(t, "LastNameUser2", user.LastName)
	assert.Equal(t, "test2@domain.com", user.Email)


}
