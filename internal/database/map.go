package database

import (
	"fmt"
	"log"
)

const (
	MaxMapWidth  = 15
	MaxMapHeight = 11
	MaxMapItems  = 20
	MaxMapNpcs   = 5
)

type MapMoral int

const (
	MapMoralNone MapMoral = iota
	MapMoralSafe
	MapMoralInn
	MapMoralArena
)

type TileType int

const (
	TileTypeWalkable TileType = iota
	TileTypeBlocked
	TileTypeWarp
	TileTypeItem
	TileTypeNpcAvoid
	TileTypeKey
	TileTypeKeyOpen
	TileTypeHeal
	TileTypeKill
	TileTypeDoor
	TileTypeSign
	TileTypeMsg
	TileTypeSprite
	TileTypeNpcSpawn
	TileTypeNudge
)

type Tile struct {
	Num   [8]int
	Type  TileType
	Data1 int
	Data2 int
	Data3 int
}

type Map struct {
	Name     string
	Revision int
	Moral    MapMoral
	TileSet  int
	Up       int
	Down     int
	Left     int
	Right    int
	Music    int
	BootMap  int
	BootX    int
	BootY    int
	Shop     int
	Tiles    [MaxMapWidth * MaxMapHeight]Tile
	Npcs     [MaxMapNpcs]int
}

var Maps [MaxMaps]Map

func init() {
	loadMaps()
}

func getMapFilename(mapId int) string {
	return fmt.Sprintf("data/maps/map%d.gob", mapId+1)
}

func loadMap(mapId int) {
	Maps[mapId].Clear()

	fileName := getMapFilename(mapId)

	err := loadFromFile(fileName, &Maps[mapId])
	if err != nil {
		log.Printf("Error loading map '%s': %s\n", fileName, err)
	}
}

func loadMaps() {
	createFolderIfNotExists("data/maps")

	for i := 0; i < MaxMaps; i++ {
		loadMap(i)
	}
}

func (m *Map) Clear() {
	m.Name = ""
	m.Revision = 0
	m.Moral = MapMoralNone
	m.TileSet = 0
	m.Up = -1
	m.Down = -1
	m.Left = -1
	m.Right = -1
	m.Music = 0
	m.BootMap = -1
	m.BootX = 0
	m.BootY = 0
	m.Shop = -1

	m.ClearTiles()
	m.ClearNpcs()
}

func (m *Map) ClearTiles() {
	for i := 0; i < len(m.Tiles); i++ {
		for j := 0; j < len(m.Tiles[i].Num); j++ {
			m.Tiles[i].Num[j] = 0
		}
		m.Tiles[i].Type = TileTypeWalkable
		m.Tiles[i].Data1 = 0
		m.Tiles[i].Data2 = 0
		m.Tiles[i].Data3 = 0
	}
}

func (m *Map) ClearNpcs() {
	for i := 0; i < len(m.Npcs); i++ {
		m.Npcs[i] = -1
	}
}

func SaveMap(mapId int) {
	if mapId < 0 || mapId >= MaxMaps {
		return
	}

	fileName := getMapFilename(mapId)

	err := saveToFile(fileName, &Maps[mapId])
	if err != nil {
		log.Printf("Error saving map '%s': %s\n", fileName, err)
	}
}

func SaveMaps() {
	for i := 0; i < MaxMaps; i++ {
		SaveMap(i)
	}
}
