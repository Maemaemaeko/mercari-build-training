package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	// STEP 5-1: uncomment this line
	// _ "github.com/mattn/go-sqlite3"
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
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
}

// NewItemRepository creates a new itemRepository.
func NewItemRepository() ItemRepository {
	return &itemRepository{fileName: "items.json"}
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item

	// JSONファイル読み込み
	file, err := os.Open(i.fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err

	}
	defer file.Close()

	// JSONデコーダを作成
	decoder := json.NewDecoder(file)

	// JSONデータを構造体にデコード
	var data Items
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	// **ID を要素数に応じて設定**
	item.ID = len(data.Items)

	// itemをdataに追加
	data.Items = append(data.Items, *item)
	// JSONにエンコード
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// JSONファイルに書き込み
	err = os.WriteFile(i.fileName, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Println("Item successfully inserted")

	return nil
}

// GetItems returns a list of items from the repository.
func (i *itemRepository) GetItems(ctx context.Context) (*Items, error) {
	// STEP 4-1: add an implementation to get items

	// JSONファイル読み込み
	file, err := os.Open(i.fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err

	}
	defer file.Close()

	// JSONデコーダを作成
	decoder := json.NewDecoder(file)

	// JSONデータを構造体にデコード
	var data Items
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}

	return &data, nil
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

	return nil
}
