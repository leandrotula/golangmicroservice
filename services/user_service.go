package services

import (
	"github.com/golangmicroservice/domain"
	"github.com/golangmicroservice/util"
)

func GetUser(id int64) (*domain.User, *util.ResponseError) {

	return domain.GetUser(id)
}
