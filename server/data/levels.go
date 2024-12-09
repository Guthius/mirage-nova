package data

import (
	"log"

	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/storage"
)

const (
	maxWidth  = 16
	maxHeight = 12
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
	Num   [9]int
	Type  TileType
	Data1 int
	Data2 int
	Data3 int
}

type LevelType int

const (
	LevelDefault LevelType = iota
	LevelSafe
	LevelInn
	LevelArena
)

type LevelData struct {
	Name     string
	Revision int
	Type     LevelType
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
	Width    int
	Height   int
	Tiles    [maxWidth * maxHeight]Tile
	Npcs     [config.MaxMapNpcs]int
}

var levelStore = storage.NewFileStore("data/levels", "level", resetLevelData)
var levels [config.MaxMaps]*LevelData

func init() {
	for i := 0; i < config.MaxMaps; i++ {
		level, err := levelStore.Load(i)
		if err != nil {
			log.Printf("error loading level %03d (%s)\n", i, err)
		}
		levels[i] = level
	}
}

// resetLevelData resets the fields of the specified LevelData back to their default values.
func resetLevelData(m *LevelData) {
	m.Name = ""
	m.Revision = 1
	m.Type = LevelDefault
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
	m.Width = maxWidth
	m.Height = maxHeight

	m.resetTiles()
	m.resetNpcs()
}

// resetTiles resets the fields of all tiles of the level back to their default values.
func (level *LevelData) resetTiles() {
	for i := 0; i < len(level.Tiles); i++ {
		for j := 0; j < len(level.Tiles[i].Num); j++ {
			level.Tiles[i].Num[j] = 0
		}
		level.Tiles[i].Type = TileTypeWalkable
		level.Tiles[i].Data1 = 0
		level.Tiles[i].Data2 = 0
		level.Tiles[i].Data3 = 0
	}
}

// resetNpcs resets the fields of all NPC's on the level back to their default values.
func (level *LevelData) resetNpcs() {
	for i := 0; i < len(level.Npcs); i++ {
		level.Npcs[i] = -1
	}
}

// SaveLevel saves the data of the level with the specified ID to the backing file store.
func SaveLevel(id int) {
	if id < 0 || id >= config.MaxMaps {
		return
	}

	err := levelStore.Save(id, levels[id])
	if err != nil {
		log.Printf("error saving level %03d (%s)\n", id, err)
	}
}

// SaveAllLevels saves the data of all levels to the backing file store.
func SaveAllLevels() {
	for i := 0; i < config.MaxMaps; i++ {
		SaveLevel(i)
	}
}

// GetLevel returns the data of the level with the specified id
func GetLevel(id int) *LevelData {
	if id < 0 || id >= config.MaxMaps {
		return nil
	}
	return levels[id]
}

// Contains return true if the specified position is within the boundaries of the level; otherwise, returns false.
func (level *LevelData) Contains(x int, y int) bool {
	return !(x < 0 || y < 0 || x >= level.Width || y >= level.Height)
}

// GetTile returns the tile at the specified position.
func (level *LevelData) GetTile(x int, y int) *Tile {
	if !level.Contains(x, y) {
		return nil
	}
	tid := y*level.Width + x
	return &level.Tiles[tid]
}

// GetTileType returns the type of the tile at the specified position.
// Returns TileTypeWalkable if the position is out of bounds.
func (level *LevelData) GetTileType(x int, y int) TileType {
	if !level.Contains(x, y) {
		return TileTypeWalkable
	}
	tid := y*level.Width + x
	return level.Tiles[tid].Type
}
