package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	// STEP 5-1: uncomment this line
	_ "github.com/mattn/go-sqlite3"
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
	CloseDB() error
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new itemRepository.
func NewItemRepository(dbPath string) ItemRepository {
	// データベースに接続
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil
	}
	// テーブル作成
	cmd := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		category_id INTEGER NOT NULL,
		image_name TEXT NOT NULL,
		FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
	);`

	_, err = db.Exec(cmd)

	if err != nil {
		return nil
	}

	return &itemRepository{db: db}
}

func (i *itemRepository) getCategoryIDFromDB(ctx context.Context, category string) (int, error) {
	var categoryID int
	err := i.db.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&categoryID)
	if err == sql.ErrNoRows {
		return i.insertCategoryInDB(ctx, category)
	}
	return categoryID, nil
}

func (i *itemRepository) insertCategoryInDB(ctx context.Context, category string) (int, error) {
	result, err := i.db.Exec("INSERT INTO categories (name) VALUES (?)", category)
	if err != nil {
		return 0, err
	}
	categoryID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(categoryID64), nil
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 5-1 Insert an item into the database
	// Set up a transaction to ensure the consistency of the data
	var categoryID int
	categoryID, err := i.getCategoryIDFromDB(ctx, item.Category)
	if err != nil {
		return err
	}

	// `items` テーブルにデータを追加（カテゴリIDが確定）
	_, err = i.db.Exec("INSERT INTO items (name, category_id, image_name) VALUES (?, ?, ?)", item.Name, categoryID, item.ImageName)
	if err != nil {
		return err
	}

	return nil
}

// GetItems returns a list of items from the repository.
func (i *itemRepository) GetItems(ctx context.Context) (*Items, error) {
	// STEP 5-1, 5-3: Get items from the database
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
	// STEP 5-1, 5-3: (Optional) Get a single item from the database
	query := `
	SELECT items.id, items.name, categories.name AS category, items.image_name
	FROM items
	INNER JOIN categories ON items.category_id = categories.id
	WHERE items.id = ?
	`

	var item Item

	err := i.db.QueryRowContext(ctx, query, id).Scan(&item.ID, &item.Name, &item.Category, &item.ImageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errItemNotFound
		}
		return nil, err
	}

	return &item, nil
}

// SearchItems returns a list of items that match the query from the repository.
func (i *itemRepository) SearchItems(ctx context.Context, keyword string) (*Items, error) {
	// STEP 5-2: Search items from the database using a keyword
	query := `
	SELECT items.id, items.name, categories.name AS category, items.image_name
	FROM items
	INNER JOIN categories ON items.category_id = categories.id
	WHERE items.name LIKE ?
	`
	fmt.Println("keyword", keyword)
	// Add % to the keyword to search for partial matches
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

// CloseDB closes the database connection.
func (i *itemRepository) CloseDB() error {
	// STEP 5-1: Close the database connection
	return i.db.Close()
}
