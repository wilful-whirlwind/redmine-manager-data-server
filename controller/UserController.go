package controller

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"pmj-server/dao"
	"pmj-server/entity"
	"pmj-server/service"
	"strconv"
)

type UserController struct {
	base *BaseController
}

func (action UserController) Execute(w http.ResponseWriter, r *http.Request) {
	action.base = PreExecute(w, r, action)
	Execute(w, r, action)
}

func (action UserController) Get(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	query := r.URL.Query()
	body := make(map[string]interface{})

	db, err := dao.Driver()
	if err != nil {
		return body, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			action.base.logger.Error("close db err", err)
		}
	}(db)

	userService := service.UserService{
		Logger: action.base.logger,
	}
	requestIdRaw := query.Get("id")
	if requestIdRaw == "all" {
		selectedUsers, err := userService.GetAll(db)
		if err != nil {
			action.base.logger.Error("ユーザ情報の登録に失敗しました。", "error info", err)
			return body, err
		}
		body["userList"] = selectedUsers
		action.base.logger.Info("return", "param", body)
		return body, err
	} else if requestIdRaw == "auth" {
		passwordRequest := query.Get("password")

		selectedUser, err := userService.GetByHash(passwordRequest, db)
		if err != nil {
			action.base.logger.Error("ユーザ情報の登録に失敗しました。", "error info", err)
			return body, err
		}
		body["result"] = selectedUser
		action.base.logger.Info("return", "param", body)
		return body, err
	} else {
		requestId, err := strconv.Atoi(requestIdRaw)
		if err != nil {
			action.base.logger.Error("リクエストされたIDが不正です。", "error info", err)
			return body, err
		}

		selectedUser, err := userService.GetById(requestId, db)
		if err != nil {
			action.base.logger.Error("ユーザ情報の登録に失敗しました。", "error info", err)
			return body, err
		}
		body["result"] = selectedUser
		action.base.logger.Info("return", "param", body)
		return body, err
	}
}

func (action UserController) Post(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	var requestBody entity.User
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return nil, err
	}

	db, err := dao.Driver()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			action.base.logger.Error("close db err", err)
		}
	}(db)
	userService := service.UserService{}
	user := userService.BindUserEntity(requestBody.MailAddress, requestBody.Name, requestBody.PasswordHash, -1)
	newUserId, err := userService.Insert(user, db)
	user.Id = int(newUserId)
	if err != nil {
		action.base.logger.Error("ユーザ情報の登録に失敗しました。", "error info", err)
		return nil, err
	}
	action.base.logger.Info("return", "param", user)
	return convertResponse(user), nil
}

func (action UserController) Patch(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	var requestBody entity.User
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return nil, err
	}

	db, err := dao.Driver()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			action.base.logger.Error("close db err", err)
		}
	}(db)
	userService := service.UserService{}
	user := userService.BindUserEntity(requestBody.MailAddress, requestBody.Name, "", requestBody.Id)
	newUserId, err := userService.Update(user, db)
	user.Id = int(newUserId)
	if err != nil {
		action.base.logger.Error("ユーザ情報の登録に失敗しました。", "error info", err)
		return nil, err
	}
	action.base.logger.Info("return", "param", user)
	return convertResponse(user), nil
}
