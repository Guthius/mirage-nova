package data

import (
	"log"

	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data/stats"
	"github.com/guthius/mirage-nova/storage"
)

type NpcBehaviour int

const (
	NpcBehaviourAttackOnSight NpcBehaviour = iota
	NpcBehaviourAttackWhenAttacked
	NpcBehaviourFriendly
	NpcBehaviourShopKeeper
	NpcBehaviourGuard
)

type NpcData struct {
	Name          string
	AttackSay     string
	Sprite        int
	SpawnSecs     int64
	Behaviour     NpcBehaviour
	Range         int
	DropChance    int
	DropItemId    int
	DropItemValue int
	Stats         stats.Data
}

var npcStore = storage.NewFileStore("npc", "data/npcs", resetNpcData)
var npcs [config.MaxNpcs]*NpcData

func init() {
	for i := 0; i < config.MaxNpcs; i++ {
		npc, err := npcStore.Load(i)
		if err != nil {
			log.Printf("Error loading npc %03d: %s\n", i, err)
		}
		npcs[i] = npc
	}
}

// resetNpcData resets the fields of the specified NpcData back to their default values.
func resetNpcData(n *NpcData) {
	n.Name = ""
	n.AttackSay = ""
	n.Sprite = 0
	n.SpawnSecs = 0
	n.Behaviour = NpcBehaviourAttackOnSight
	n.Range = 0
	n.DropChance = 0
	n.DropItemId = -1
	n.DropItemValue = 0
	n.Stats.Reset()
}

// SaveNpc saves the data of the NPC with the specified ID to the backing file store.
func SaveNpc(id int) {
	if id < 0 || id >= config.MaxNpcs {
		return
	}

	err := npcStore.Save(id, npcs[id])
	if err != nil {
		log.Printf("Error saving npc %03d: %3\n", id, err)
	}
}

// SaveAllNpcs saves the data of all NPC's to the backing file store.
func SaveAllNpcs() {
	for i := 0; i < config.MaxNpcs; i++ {
		SaveNpc(i)
	}
}

// GetNpc returns the data of the NPC with the specified id
func GetNpc(id int) *NpcData {
	if id < 0 || id >= config.MaxNpcs {
		return nil
	}
	return npcs[id]
}
