package controller

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"pmj-server/dao"
	"pmj-server/entity"
	"pmj-server/service"
)

type ConfigController struct {
	base *BaseController
}

func (action ConfigController) Execute(w http.ResponseWriter, r *http.Request) {
	action.base = PreExecute(w, r, action)
	Execute(w, r, action)
}

func (action ConfigController) Get(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	query := r.URL.Query()
	body := make(map[string]interface{})

	db, err := dao.Driver()
	if err != nil {
		return body, err
	}

	configService := service.ConfigService{
		Logger: action.base.logger,
	}
	key := query.Get("key")
	configuration, err := configService.GetByKey(key, db)
	if err != nil {
		action.base.logger.Error("ユーザ情報の登録に失敗しました。", "error info", err)
		return body, err
	}
	body["result"] = configuration
	action.base.logger.Info("return", "param", body)
	return body, err
}

func (action ConfigController) Post(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	var requestBody entity.Configuration
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return nil, err
	}

	db, err := dao.Driver()
	configService := service.ConfigService{}
	configuration := configService.BindConfigurationEntity(requestBody.Key, requestBody.Value)
	_, err = configService.Delete(requestBody.Key, db)
	if err != nil {
		return nil, err
	}
	newId, err := configService.Insert(configuration, db)
	configuration.Id = int(newId)
	if err != nil {
		action.base.logger.Error("ユーザ情報の登録に失敗しました。", "error info", err)
		return nil, err
	}
	action.base.logger.Info("return", "param", configuration)
	return convertResponse(configuration), nil
}
