package dao

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"log/slog"
	"pmj-server/entity"
)

type UserDao struct {
	Driver *sql.DB
	Logger *slog.Logger
}

func (dao UserDao) Insert(user entity.User) (int64, error) {
	sqlStr := `
		insert into 
		    users(
			  mail_address,
			  name,
			  password_hash,
			  created_at
			) values (?, ?, ?, ?)
	`
	result, err := dao.Driver.Exec(sqlStr, user.MailAddress, user.Name, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return -1, errors.New("登録に失敗しました。")
	}
	newUserId, errResult := result.LastInsertId()
	if errResult != nil {
		return -1, errors.New("登録に失敗しました。")
	}
	return newUserId, err
}

func (dao UserDao) GetById(id int) (entity.User, error) {
	sqlStr := `
		SELECT
		    id,
		    mail_address,
		    name,
		    created_at
		FROM
		    users
		WHERE
		    id = ?
	`
	rows, err := dao.Driver.Query(sqlStr, id)
	if err != nil {
		return entity.User{}, errors.New("取得に失敗しました。")
	}
	selectedUser := entity.User{}
	for rows.Next() {
		if err := rows.Scan(&selectedUser.Id, &selectedUser.MailAddress, &selectedUser.Name, &selectedUser.CreatedAt); err != nil {
			log.Fatalf("failed to scan row: %s", err)
		}
	}
	dao.Logger.Info("result", "entity", selectedUser)
	return selectedUser, err
}
