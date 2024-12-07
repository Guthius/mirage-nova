package main

const (
	MaxMapWidth  = 15
	MaxMapHeight = 11
	MaxMapItems  = 20
	MaxMapNpcs   = 5
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

type MapMoral int

const (
	MapMoralNone MapMoral = iota
	MapMoralSafe
	MapMoralInn
	MapMoralArena
)

type LevelData struct {
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
