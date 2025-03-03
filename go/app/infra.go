package app

import (
	"context"
	"database/sql"
	"errors"
	"os"

	// STEP 5-1: uncomment this line
	_ "github.com/mattn/go-sqlite3" // SQLite ドライバを import
)

var errImageNotFound = errors.New("image not found")
var errItemNotFound = errors.New("item not found")

type Item struct {
	ID        int    `db:"id" json:"-"`
	Name      string `db:"name" json:"name"`
	Category  string `db:"category" json:"category"`
	ImageName string `db:"image_name" json:"image_name"`
}

// Items 構造体（JSON全体を表す）
type Items struct {
	Items []Item `json:"items"`
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
	GetItems(ctx context.Context) (*Items, error)
	GetItem(ctx context.Context, id string) (*Item, error)
	SearchItems(ctx context.Context, keyword string) (*Items, error)
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
	// db is the database connection.
	db *sql.DB
}

// NewItemRepository creates a new itemRepository.
func NewItemRepository(dbPath string) ItemRepository {
	// データベースに接続
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil
	}
	return &itemRepository{fileName: "items.json", db: db}
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item
	// sqlite3にデータを挿入する
	tx, err := i.db.Begin()
	if err != nil {
		return err
	}
	var categoryID int
	err = tx.QueryRow("SELECT id FROM categories WHERE name = ?", item.Category).Scan(&categoryID)
	if err == sql.ErrNoRows {
		// カテゴリが存在しない場合、新しく追加
		result, err := tx.Exec("INSERT INTO categories (name) VALUES (?)", item.Category)
		if err != nil {
			tx.Rollback()
			return err
		}
		// 新しく作成した `category_id` を取得
		categoryID64, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
		categoryID = int(categoryID64)
	} else if err != nil {
		tx.Rollback()
		return err
	}

	// `items` テーブルにデータを追加（カテゴリIDが確定）
	_, err = tx.Exec("INSERT INTO items (name, category_id, image_name) VALUES (?, ?, ?)", item.Name, categoryID, item.ImageName)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetItems returns a list of items from the repository.
func (i *itemRepository) GetItems(ctx context.Context) (*Items, error) {
	// 5-3
	query := `
	SELECT items.id, items.name, categories.name AS category, items.image_name
	FROM items 
	INNER JOIN categories ON items.category_id = categories.id
`
	rows, err := i.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items Items
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.ImageName)
		if err != nil {
			return nil, err
		}
		items.Items = append(items.Items, item)
	}
	return &items, nil
}

// GetItem returns an item from the repository.
func (i *itemRepository) GetItem(ctx context.Context, id string) (*Item, error) {
	// STEP 4-1: add an implementation to get an item
	// SQLクエリの作成
	query := `
	SELECT items.id, items.name, categories.name AS category, items.image_name
	FROM items
	INNER JOIN categories ON items.category_id = categories.id
	WHERE items.id = ?
	`
	// データ取得用の `Item` 構造体
	var item Item

	// `QueryRowContext` を使って 1 件のデータを取得
	err := i.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name, &item.Category, &item.ImageName)
	if err != nil {
		if err == sql.ErrNoRows {
			// 指定された `id` の商品が見つからなかった場合
			return nil, errItemNotFound
		}
		// それ以外のエラー
		return nil, err
	}

	// `Item` を返す
	return &item, nil
}

// SearchItems returns a list of items that match the query from the repository.
func (i *itemRepository) SearchItems(ctx context.Context, keyword string) (*Items, error) {

	// SQLクエリの作成
	// query := "SELECT name, category, image_name FROM items WHERE name LIKE ?"
	query := `
	SELECT items.id, items.name, categories.name AS category, items.image_name
	FROM items
	INNER JOIN categories ON items.category_id = categories.id
	WHERE items.name LIKE ?
	`

	// ワイルドカード検索のため `%` を付ける
	rows, err := i.db.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items Items
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.ImageName)
		if err != nil {
			return nil, err
		}
		items.Items = append(items.Items, item)
	}

	return &items, nil
}

// StoreImage stores an image and returns an error if any.
// This package doesn't have a related interface for simplicity.
func StoreImage(fileName string, image []byte) error {
	// STEP 4-4: add an implementation to store an image
	if err := os.WriteFile(fileName, image, 0644); err != nil {
		return err
	}

	return nil
}
