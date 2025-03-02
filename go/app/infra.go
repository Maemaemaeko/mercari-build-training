package app

import (
	"context"
	"errors"
	// "os"
	"fmt"
	"io/ioutil"
	// "log"
	"encoding/json"
	"strconv"
	"database/sql"
	"os"
	

	// STEP 5-1: uncomment this line
	_ "github.com/mattn/go-sqlite3" // SQLite ドライバを import
)

var errImageNotFound = errors.New("image not found")
var errItemNotFound = errors.New("item not found")

type Item struct {
	ID   int    `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
	ImagePath string `db:"image" json:"image"`
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
	GetItems(ctx context.Context) ([]byte, error)
	GetItem(ctx context.Context, id string) (*Item, error)
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
	// db is the database connection.
	db *sql.DB
}

// NewItemRepository creates a new itemRepository.
func NewItemRepository(dbPath string) (ItemRepository) {
	// データベースに接続
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil
	}

	return &itemRepository{fileName: "items.json", db: db}
}

// Items 構造体（JSON全体を表す）
type Items struct {
	Items []Item `json:"items"`
}


// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item
	// sqlite3にデータを挿入
	_, err := i.db.Exec("INSERT INTO items (name, category, image_name) VALUES (?, ?, ?)", item.Name, item.Category, item.ImagePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}



	fmt.Println("Item successfully inserted")

	return nil
}

// GetItems returns a list of items from the repository.
func (i *itemRepository) GetItems(ctx context.Context) ([]byte, error) {
	// 5-1 sqlite3からデータを取得
	rows, err := i.db.Query("SELECT * FROM items")
	if err != nil {
		fmt.Println("Error reading file:", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next(){
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.ImagePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
		}
		items = append(items, item)
	}

	// JSON にエンコード
	response := Items{Items: items}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}


// GetItem returns an item from the repository.
func (i *itemRepository) GetItem(ctx context.Context, id string) (*Item, error) {
	// STEP 4-1: add an implementation to get an item

	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile(i.fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	var data Items
	if err := json.Unmarshal(bytes, &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	// IDを数値に変換
	id_int, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting id to int:", err)
	}

	// IDが範囲内か確認
	if id_int < 0 || id_int >= len(data.Items) {
		fmt.Println("Error: ID out of range")
		return nil, errItemNotFound
	}

	return &data.Items[id_int], nil
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
