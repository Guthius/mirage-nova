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
	loadNpcs()
}

func getNpcFilename(npcId int) string {
	return fmt.Sprintf("data/npcs/npc%d.gob", npcId+1)
}

func loadNpc(npcId int) {
	Npcs[npcId].Clear()

	fileName := getNpcFilename(npcId)

	err := loadFromFile(fileName, &Npcs[npcId])
	if err != nil {
		log.Printf("Error loading npc '%s': %s\n", fileName, err)
	}
}

func loadNpcs() {
	createFolderIfNotExists("data/npcs")

	for i := 0; i < MaxNpcs; i++ {
		loadNpc(i)
	}
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

func (n *Npc) Clear() {
	n.Name = ""
	n.AttackSay = ""
	n.Sprite = 0
	n.SpawnSecs = 0
	n.Behaviour = NpcBehaviourAttackOnSight
	n.Range = 0
	n.DropChance = 0
	n.DropItemId = -1
	n.DropItemValue = 0
	n.Stats.Strength = 0
	n.Stats.Defense = 0
	n.Stats.Speed = 0
	n.Stats.Magic = 0
}
