package database

import (
	"fmt"
)

type NpcBehaviour int

const (
	NpcBehaviourAttackOnSight NpcBehaviour = iota
	NpcBehaviourAttackWhenAttacked
	NpcBehaviourFriendly
	NpcBehaviourShopKeeper
	NpcBehaviourGuard
)

type Npc struct {
	Name          string
	AttackSay     string
	Sprite        int
	SpawnSecs     int64
	Behaviour     NpcBehaviour
	Range         int
	DropChance    int
	DropItemId    int
	DropItemValue int
	Stats         Stats
}

var Npcs = loadNpcs()

func getNpcFilename(npcId int) string {
	return fmt.Sprintf("data/npcs/npc%d.gob", npcId+1)
}

func loadNpc(npcId int) Npc {
	npc := Npc{
		Behaviour:  NpcBehaviourAttackOnSight,
		DropItemId: -1,
	}

	fileName := getNpcFilename(npcId)

	err := loadFromFile(fileName, &npc)
	if err != nil {
		Log.Printf("Error loading npc '%s': %s\n", fileName, err)
	}

	return npc
}

func loadNpcs() [MAX_NPCS]Npc {
	var npcs [MAX_NPCS]Npc

	createFolderIfNotExists("data/npcs")

	for i := 0; i < MAX_NPCS; i++ {
		npcs[i] = loadNpc(i)
	}

	return npcs
}

func SaveNpc(npcId int) {
	if npcId < 0 || npcId >= MAX_NPCS {
		return
	}

	fileName := getNpcFilename(npcId)

	err := saveToFile(fileName, &Npcs[npcId])
	if err != nil {
		Log.Printf("Error saving npc '%s': %s\n", fileName, err)
	}
}

func SaveNpcs() {
	for i := 0; i < MAX_NPCS; i++ {
		SaveNpc(i)
	}
}

func GetNpc(npcId int) *Npc {
	if npcId < 0 || npcId >= MAX_NPCS {
		return nil
	}
	return &Npcs[npcId]
}
