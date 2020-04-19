package domain

import (
	"github.com/leandrotula/golangmicroservice/util"
	"net/http"
)

var (
	userData = map[int64]User{
		1: {
			Id:        1,
			FirstName: "TestFirstName",
			LastName:  "TestLastName",
			Email:     "test@domain.com",
		},
	}

	UserDao userDaoInterface
)

func init() {

	UserDao = &userDaoImpl{}
}

type userDaoInterface interface {

	GetUser(userId int64)(*User, *util.ResponseError)
}

type userDaoImpl struct {
}

func(u *userDaoImpl) GetUser(userId int64)(*User, *util.ResponseError) {

	user, present := userData[userId]

	if !present {
		return nil, &util.ResponseError{
			Message: "No user found",
			Code:    http.StatusNotFound,
		}
	}

	return &user, nil
}
