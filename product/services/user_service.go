package services

import (
	"imooc-product/datamodels"
	"imooc-product/repositories"

	"github.com/kataras/iris/v12/x/errors"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userID int64, err error)
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword) //这里user.HashPassword为空，检查一下
	if !isOk {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

// 自定义加密函数,明文密码转为密文，保护用户密码
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// 验证密码是否正确
func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误!")
	}
	return true, nil
}
