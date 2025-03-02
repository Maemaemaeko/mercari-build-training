package app

import (
	"context"
	"errors"
	// "os"
	"fmt"
	"io/ioutil"
	// "log"
	"encoding/json"

	// STEP 5-1: uncomment this line
	// _ "github.com/mattn/go-sqlite3"
)

var errImageNotFound = errors.New("image not found")

type Item struct {
	ID   int    `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
	Category string `db:"category" json:"category"`
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
	GetItems(ctx context.Context) ([]byte, error)
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

// Items 構造体（JSON全体を表す）
type Items struct {
	Items []Item `json:"items"`
}


// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item
	// fmt.Printf("Inserting item: %+v\n", item)
	// fmt.Printf("Inserting item_name: %+v\n", item.Name)
	

    // JSONファイル読み込み
    bytes, err := ioutil.ReadFile(i.fileName)
    if err != nil {
		fmt.Println("Error reading file:", err)
    }	
	// fmt.Printf("読み込んだデータ: %s\n", bytes)
	var data Items
	if err := json.Unmarshal(bytes, &data); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	// fmt.Printf("デコード後のデータ: %+v\n", data)

	// デコードしたデータを表示
	// for _, item := range data.Items {
	// 	fmt.Printf("%s %s\n",item.Name, item.Category)
	// }

	// **ID を要素数に応じて設定**
	item.ID = len(data.Items)
	fmt.Printf("ID: %d\n", item.ID)

	// itemをdataに追加
	data.Items = append(data.Items, *item)
	// fmt.Printf("追加後のデータ: %+v\n", data)

	// JSONにエンコード
	// bytes, err = json.Marshal(data)
	bytes, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// JSONファイルに書き込み
	if err := ioutil.WriteFile(i.fileName, bytes, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Println("Item successfully inserted")

	return nil
}

// GetItems returns a list of items from the repository.
func (i *itemRepository) GetItems(ctx context.Context) ([]byte, error) {
	// STEP 4-1: add an implementation to get items

	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile(i.fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	return bytes, nil
}

// StoreImage stores an image and returns an error if any.
// This package doesn't have a related interface for simplicity.
func StoreImage(fileName string, image []byte) error {
	// STEP 4-4: add an implementation to store an image

	return nil
}
