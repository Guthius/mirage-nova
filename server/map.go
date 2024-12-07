package main

import "mirage/internal/database"

type TempTile struct {
	DoorOpen bool
}

type TempMap struct {
	Cache       []byte
	PlayerCount int
	DoorTimer   int64
	Tiles       [database.MaxMapWidth * database.MaxMapHeight]TempTile
}

var TempMaps [database.MaxMaps]TempMap

func init() {
	for i := 0; i < database.MaxMaps; i++ {
		TempMaps[i].Cache = nil
		TempMaps[i].PlayerCount = 0
		TempMaps[i].DoorTimer = 0

		for j := 0; j < len(TempMaps[i].Tiles); j++ {
			TempMaps[i].Tiles[j].DoorOpen = false
		}
	}
}
