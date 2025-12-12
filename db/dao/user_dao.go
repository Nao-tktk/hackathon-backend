package dao

import (
	"database/sql"
	"db/model"
)

type UserDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{db: db}
}

func (d *UserDao) FindByName(name string) ([]model.User, error) {
	rows, err := d.db.Query("SELECT id, name, password FROM users WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (d *UserDao) Insert(user *model.User) (int, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	
	result, err := tx.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(id64), nil
}
