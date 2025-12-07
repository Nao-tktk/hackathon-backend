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
	rows, err := d.db.Query("SELECT id, name, age FROM user WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (d *UserDao) Insert(user *model.User) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", user.ID, user.Name, user.Age)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
