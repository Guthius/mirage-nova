package database

import (
	"fmt"
	"log"
)

type ItemType int

const (
	ITEM_TYPE_NONE ItemType = iota
	ITEM_TYPE_WEAPON
	ITEM_TYPE_ARMOR
	ITEM_TYPE_HELMET
	ITEM_TYPE_SHIELD
	ITEM_TYPE_POTIONADDHP
	ITEM_TYPE_POTIONADDMP
	ITEM_TYPE_POTIONADDSP
	ITEM_TYPE_POTIONSUBHP
	ITEM_TYPE_POTIONSUBMP
	ITEM_TYPE_POTIONSUBSP
	ITEM_TYPE_KEY
	ITEM_TYPE_CURRENCY
	ITEM_TYPE_SPELL
)

type Item struct {
	Name  string
	Pic   int
	Type  ItemType
	Data1 int
	Data2 int
	Data3 int
}

var Items [MaxItems]Item

func init() {
	loadItems()
}

func getItemFilename(itemId int) string {
	return fmt.Sprintf("data/items/item%d.gob", itemId+1)
}

func loadItem(itemId int) {
	Items[itemId].Clear()

	fileName := getItemFilename(itemId)

	err := loadFromFile(fileName, &Items[itemId])
	if err != nil {
		log.Printf("error loading item '%s' (%s)\n", fileName, err)
	}
}

func loadItems() {
	createFolderIfNotExists("data/items")

	for i := 0; i < MaxItems; i++ {
		loadItem(i)
	}
}

func SaveItem(itemId int) {
	if itemId < 0 || itemId >= MaxNpcs {
		return
	}

	fileName := getItemFilename(itemId)

	err := saveToFile(fileName, &Items[itemId])
	if err != nil {
		log.Printf("error saving item '%s' (%s)\n", fileName, err)
	}
}

func SaveItems() {
	for i := 0; i < MaxItems; i++ {
		SaveItem(i)
	}
}

func GetItem(itemId int) *Item {
	if itemId < 0 || itemId >= MaxItems {
		return nil
	}
	return &Items[itemId]
}

func (item *Item) Clear() {
	item.Name = ""
	item.Pic = 0
	item.Type = ITEM_TYPE_NONE
	item.Data1 = 0
	item.Data2 = 0
	item.Data3 = 0
}
