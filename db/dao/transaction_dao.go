package dao

import (
	"database/sql"
	"fmt"
)

type TransactionDao struct {
	db *sql.DB
}

func NewTransactionDao(db *sql.DB) *TransactionDao {
	return &TransactionDao{db: db}
}

// Purchase はトランザクションを使って「購入履歴保存」と「商品ステータス更新」を一気に行います
func (dao *TransactionDao) Purchase(itemID int, buyerID int) error {
	// 1. トランザクション開始 (失敗したら全部なかったことにする機能)
	tx, err := dao.db.Begin()
	if err != nil {
		return err
	}

	// 2. 商品がまだ売れていないかチェック (ON_SALEじゃなければエラー)
	// FOR UPDATE をつけることで、同時に誰かが買おうとしてもロックできる
	var status string
	err = tx.QueryRow("SELECT status FROM items WHERE id = ? FOR UPDATE", itemID).Scan(&status)
	if err != nil {
		tx.Rollback()
		return err
	}
	if status != "ON_SALE" {
		tx.Rollback()
		return fmt.Errorf("item is already sold out")
	}

	// 3. itemsテーブルのステータスを SOLD_OUT に更新
	_, err = tx.Exec("UPDATE items SET status = 'SOLD_OUT' WHERE id = ?", itemID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 4. transactionsテーブルに購入記録を追加
	_, err = tx.Exec("INSERT INTO transactions (item_id, buyer_id) VALUES (?, ?)", itemID, buyerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 5. 全部成功したので確定！
	return tx.Commit()
}
