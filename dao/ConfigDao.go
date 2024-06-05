package dao

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"log/slog"
	"pmj-server/entity"
)

type ConfigDao struct {
	Driver *sql.DB
	Logger *slog.Logger
}

func (dao ConfigDao) Insert(key string, value string, createdAt string) (int64, error) {
	sqlStr := `
		insert into 
		    configurations(
			  configuration_key,
			  configuration_value,
			  created_at,
			  updated_at
			) values (?, ?, ?, ?)
	`
	result, err := dao.Driver.Exec(sqlStr, key, value, createdAt, createdAt)
	if err != nil {
		return -1, errors.New("登録に失敗しました。")
	}
	newUserId, errResult := result.LastInsertId()
	if errResult != nil {
		return -1, errors.New("登録に失敗しました。")
	}
	return newUserId, err
}

func (dao ConfigDao) GetByKey(key string) (entity.Configuration, error) {
	sqlStr := `
		SELECT
		    id,
		    configuration_key,
		    configuration_value,
		    created_at,
		    updated_at
		FROM
		    configurations
		WHERE
		    configuration_key = ?
	`
	rows, err := dao.Driver.Query(sqlStr, key)
	configuration := entity.Configuration{}
	if err != nil {
		return configuration, errors.New("取得に失敗しました。")
	}
	for rows.Next() {
		if err := rows.Scan(&configuration.Id, &configuration.Key, &configuration.Value, &configuration.CreatedAt, &configuration.UpdatedAt); err != nil {
			log.Fatalf("failed to scan row: %s", err)
		}
	}
	dao.Logger.Info("result", "entity", configuration)
	return configuration, err
}

func (dao ConfigDao) Delete(key string) error {
	sqlStr := `
		delete from 
		    configurations
		where
			configuration_key = ?
	`
	_, err := dao.Driver.Exec(sqlStr, key)
	if err != nil {
		return errors.New("削除に失敗しました。")
	}
	return nil
}

func (dao ConfigDao) GetAll() ([]entity.Configuration, error) {
	sqlStr := `
		SELECT
		    id,
		    configuration_key,
		    configuration_value,
		    created_at,
		    updated_at
		FROM
		    configurations
	`
	rows, err := dao.Driver.Query(sqlStr)
	result := make([]entity.Configuration, 0)
	if err != nil {
		return result, errors.New("取得に失敗しました。")
	}
	for rows.Next() {
		configuration := entity.Configuration{}
		if err := rows.Scan(&configuration.Id, &configuration.Key, &configuration.Value, &configuration.CreatedAt, &configuration.UpdatedAt); err != nil {
			log.Fatalf("failed to scan row: %s", err)
		}
		result = append(result, configuration)
	}
	dao.Logger.Info("result", "entity", result)
	return result, err
}
