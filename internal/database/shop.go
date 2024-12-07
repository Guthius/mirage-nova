package database

import (
	"fmt"
	"log"
)

type TradeItem struct {
	GiveItemId int
	GiveValue  int
	GetItemId  int
	GetValue   int
}

type Shop struct {
	Name       string
	JoinSay    string
	LeaveSay   string
	FixesItems bool
	TradeItems [MaxTrades]TradeItem
}

var Shops [MaxShops]Shop

func init() {
	loadShops()
}

func getShopFilename(shopId int) string {
	return fmt.Sprintf("data/shops/shop%d.gob", shopId+1)
}

func loadShop(shopId int) {
	Shops[shopId].Clear()

	fileName := getShopFilename(shopId)

	err := loadFromFile(fileName, &Shops[shopId])
	if err != nil {
		log.Printf("Error loading shop '%s' (%s)\n", fileName, err)
	}
}

func loadShops() {
	createFolderIfNotExists("data/shops")

	for i := 0; i < MaxShops; i++ {
		loadShop(i)
	}
}

func SaveShop(shopId int) {
	if shopId < 0 || shopId >= MaxShops {
		return
	}

	fileName := getShopFilename(shopId)

	err := saveToFile(fileName, &Shops[shopId])
	if err != nil {
		log.Printf("Error saving shop '%s' (%s)\n", fileName, err)
	}
}

func SaveShops() {
	for i := 0; i < MaxShops; i++ {
		SaveShop(i)
	}
}

func GetShop(shopId int) *Shop {
	if shopId < 0 || shopId >= MaxShops {
		return nil
	}
	return &Shops[shopId]
}

func (item *TradeItem) Clear() {
	item.GiveItemId = -1
	item.GiveValue = 0
	item.GetItemId = -1
	item.GetValue = 0
}

func (shop *Shop) Clear() {
	shop.Name = ""
	shop.JoinSay = ""
	shop.LeaveSay = ""
	shop.FixesItems = false
	for i := range shop.TradeItems {
		shop.TradeItems[i].Clear()
	}
}
