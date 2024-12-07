package data

import (
	"log"

	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/storage"
)

type ItemType int

const (
	ItemNone ItemType = iota
	ItemWeapon
	ItemArmor
	ItemHelmet
	ItemShield
	ItemPotionAddHP
	ItemPotionAddMP
	ItemPotionAddSP
	ItemPotionSubHP
	ItemPotionSubMP
	ItemPotionSubSP
	ItemKey
	ItemCurrency
	ItemSpell
)

type ItemData struct {
	Name  string
	Pic   int
	Type  ItemType
	Data1 int
	Data2 int
	Data3 int
}

var itemStore = storage.NewFileStore("", "data/items", resetItemData)
var items [config.MaxItems]*ItemData

func init() {
	for i := 0; i < config.MaxItems; i++ {
		item, err := itemStore.Load(i)
		if err != nil {
			log.Printf("Error loading item %03d: %s\n", i, err)
		}
		items[i] = item
	}
}

// resetItemData resets the fields of the specified ItemData back to their default values.
func resetItemData(item *ItemData) {
	item.Name = ""
	item.Pic = 0
	item.Type = ItemNone
	item.Data1 = 0
	item.Data2 = 0
	item.Data3 = 0
}

// SaveItem saves the data of the item with the specified ID to the backing file store.
func SaveItem(id int) {
	if id < 0 || id >= config.MaxItems {
		return
	}

	err := itemStore.Save(id, items[id])
	if err != nil {
		log.Printf("Error saving item %03d: %3\n", id, err)
	}
}

// SaveAllItems saves the data of all items to the backing file store.
func SaveAllItems() {
	for i := 0; i < config.MaxItems; i++ {
		SaveItem(i)
	}
}

// GetItem returns the data of the item with the specified id
func GetItem(id int) *ItemData {
	if id < 0 || id >= config.MaxItems {
		return nil
	}
	return items[id]
}

// IsEquipable returns true if the specified item type can be equipped; otherwise, returns false.
func IsEquipable(it ItemType) bool {
	return it == ItemWeapon || it == ItemArmor || it == ItemHelmet || it == ItemShield
}

// IsEquipable returns true if the item can be equipped; otherwise, returns false.
func (i *ItemData) IsEquipable() bool {
	return IsEquipable(i.Type)
}

// IsCurrency returns true if the item represents a currency; otherwise, returns false.
func (i *ItemData) IsCurrency() bool {
	return i.Type == ItemCurrency
}
