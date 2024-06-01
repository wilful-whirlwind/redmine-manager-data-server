package service

import (
	"database/sql"
	"log/slog"
	"pmj-server/dao"
	"pmj-server/entity"
	"time"
)

type ConfigService struct {
	Logger *slog.Logger
}

func (service ConfigService) GetByKey(key string, driver *sql.DB) (entity.Configuration, error) {
	configDao := dao.ConfigDao{Driver: driver, Logger: service.Logger}
	configuration, err := configDao.GetByKey(key)
	return configuration, err
}

func (service ConfigService) Insert(config entity.Configuration, driver *sql.DB) (int64, error) {
	configDao := dao.ConfigDao{Driver: driver}
	createdUserId, err := configDao.Insert(config.Key, config.Value, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return -1, err
	}
	return createdUserId, nil
}

func (service ConfigService) Delete(key string, driver *sql.DB) (int64, error) {
	configDao := dao.ConfigDao{Driver: driver}
	err := configDao.Delete(key)
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (service ConfigService) BindConfigurationEntity(key string, value string) entity.Configuration {
	return entity.Configuration{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}
