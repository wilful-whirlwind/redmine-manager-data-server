package service

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/argon2"
	"log/slog"
	"pmj-server/dao"
	"pmj-server/entity"
	"pmj-server/util"
	"time"
)

type UserService struct {
	Logger *slog.Logger
}

func (service UserService) GetById(userId int, driver *sql.DB) (entity.User, error) {
	userDao := dao.UserDao{Driver: driver, Logger: service.Logger}
	selectedUser, err := userDao.GetById(userId)
	return selectedUser, err
}

func (service UserService) Insert(user entity.User, driver *sql.DB) (int64, error) {
	userDao := dao.UserDao{Driver: driver}
	createdUserId, err := userDao.Insert(user)
	if err != nil {
		return -1, err
	}
	return createdUserId, nil
}

func (service UserService) CalculatePasswordHash(password string) string {
	envUtil := util.EnvUtil{}
	salt := envUtil.GetPasswordSalt()
	pepper := envUtil.GetPasswordPepper()

	return fmt.Sprintf("%x", argon2.Key([]byte(password), []byte(salt+pepper), 3, 32*1024, 4, 32))
}

func (service UserService) BindUserEntity(mailAddress string, name string, rawPassword string, id int) entity.User {
	return entity.User{
		Id:           id,
		MailAddress:  mailAddress,
		Name:         name,
		PasswordHash: service.CalculatePasswordHash(rawPassword),
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (service UserService) GetAll(driver *sql.DB) ([]entity.User, error) {
	userDao := dao.UserDao{Driver: driver, Logger: service.Logger}
	selectedUser, err := userDao.GetAll()
	if err != nil {
		return make([]entity.User, 0), err
	}
	return selectedUser, nil
}

func (service UserService) GetByHash(request string, driver *sql.DB) (entity.User, error) {
	userDao := dao.UserDao{Driver: driver, Logger: service.Logger}
	passwordHash := service.CalculatePasswordHash(request)
	selectedUser, err := userDao.GetByHash(passwordHash)
	return selectedUser, err
}

func (service UserService) Update(user entity.User, driver *sql.DB) (int64, error) {
	userDao := dao.UserDao{Driver: driver}
	createdUserId, err := userDao.Update(user)
	if err != nil {
		return -1, err
	}
	return createdUserId, nil
}
