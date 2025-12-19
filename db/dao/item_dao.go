package dao

import (
	"database/sql"
	"db/model"
	"fmt"
)

type ItemDao struct {
	db *sql.DB
}

func NewItemDao(db *sql.DB) *ItemDao {
	return &ItemDao{db: db}
}

func (dao *ItemDao) GetItems() ([]model.Item, error) {
	query := `
		SELECT 
			i.id, i.name, c.name as category_name, i.price, i.description, i.status, i.seller_id, i.image_name
		FROM items i
		JOIN categories c ON i.category_id = c.id
	`
	rows, err := dao.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var i model.Item
		// Scanの順番に注意
		if err := rows.Scan(&i.ID, &i.Name, &i.CategoryName, &i.Price, &i.Description, &i.Status, &i.SellerID, &i.ImageName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

// Insert: 商品出品
func (dao *ItemDao) Insert(item *model.Item) (int, error) {
	tx, err := dao.db.Begin()
	if err != nil {
		return 0, err
	}

	// statusはデフォルト(ON_SALE)任せ
	query := `INSERT INTO items (seller_id, category_id, name, price, description, image_name) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, item.SellerID, item.CategoryID, item.Name, item.Price, item.Description, item.ImageName)
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
