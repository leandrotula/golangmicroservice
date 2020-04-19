package services

import (
	"github.com/leandrotula/golangmicroservice/domain"
	"github.com/leandrotula/golangmicroservice/util"
)

func GetUser(id int64) (*domain.User, *util.ResponseError) {

	return domain.GetUser(id)
}
