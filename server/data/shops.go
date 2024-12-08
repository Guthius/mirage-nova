package data

import (
	"log"

	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/storage"
)

type TradeItemData struct {
	GiveItemId int
	GiveValue  int
	GetItemId  int
	GetValue   int
}

type ShopData struct {
	Name       string
	JoinSay    string
	LeaveSay   string
	FixesItems bool
	TradeItems [config.MaxTrades]TradeItemData
}

var shopStore = storage.NewFileStore("data/shops", "shop", resetShopData)
var shops [config.MaxShops]*ShopData

func init() {
	for i := 0; i < config.MaxShops; i++ {
		shop, err := shopStore.Load(i)
		if err != nil {
			log.Printf("Error loading shop %03d (%s)\n", i, err)
		}
		shops[i] = shop
	}
}

// resetShopData resets the fields of the specified ShopData back to their default values.
func resetShopData(shop *ShopData) {
	shop.Name = ""
	shop.JoinSay = ""
	shop.LeaveSay = ""
	shop.FixesItems = false
	for i := range shop.TradeItems {
		shop.TradeItems[i].reset()
	}
}

// reset resets the fields of the specified TradeItemData back to their default values.
func (item *TradeItemData) reset() {
	item.GiveItemId = -1
	item.GiveValue = 0
	item.GetItemId = -1
	item.GetValue = 0
}

// SaveShop saves the data of the shop with the specified ID to the backing file store.
func SaveShop(id int) {
	if id < 0 || id >= config.MaxShops {
		return
	}

	err := shopStore.Save(id, shops[id])
	if err != nil {
		log.Printf("error saving shop %03d (%s)\n", id, err)
	}
}

// SaveAllShops saves the data of all shops to the backing file store.
func SaveAllShops() {
	for i := 0; i < config.MaxShops; i++ {
		SaveShop(i)
	}
}

// GetShop returns the data of the shop with the specified id
func GetShop(id int) *ShopData {
	if id < 0 || id >= config.MaxShops {
		return nil
	}
	return shops[id]
}
