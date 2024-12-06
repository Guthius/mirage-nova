package database

import (
	"fmt"
	"log"
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

var Npcs [MaxNpcs]Npc

func init() {
	Npcs = loadNpcs()
}

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
		log.Printf("Error loading npc '%s': %s\n", fileName, err)
	}

	return npc
}

func loadNpcs() [MaxNpcs]Npc {
	var npcs [MaxNpcs]Npc

	createFolderIfNotExists("data/npcs")

	for i := 0; i < MaxNpcs; i++ {
		npcs[i] = loadNpc(i)
	}

	return npcs
}

func SaveNpc(npcId int) {
	if npcId < 0 || npcId >= MaxNpcs {
		return
	}

	fileName := getNpcFilename(npcId)

	err := saveToFile(fileName, &Npcs[npcId])
	if err != nil {
		log.Printf("Error saving npc '%s': %s\n", fileName, err)
	}
}

func SaveNpcs() {
	for i := 0; i < MaxNpcs; i++ {
		SaveNpc(i)
	}
}

func GetNpc(npcId int) *Npc {
	if npcId < 0 || npcId >= MaxNpcs {
		return nil
	}
	return &Npcs[npcId]
}
