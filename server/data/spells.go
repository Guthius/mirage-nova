package data

import (
	"log"

	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/storage"
)

type SpellType int

const (
	SpellAddHP SpellType = iota
	SpellAddMP
	SpellAddSP
	SpellSubHP
	SpellSubMP
	SpellSubSP
	SpellGiveItem
)

type SpellData struct {
	Name     string
	Pic      int
	MPReq    int
	ClassReq int
	LevelReq int
	Type     SpellType
	Data1    int
	Data2    int
	Data3    int
}

var spellStore = storage.NewFileStore("spell", "data/spells", resetSpellData)
var spells [config.MaxSpells]*SpellData

func init() {
	for i := 0; i < config.MaxSpells; i++ {
		spell, err := spellStore.Load(i)
		if err != nil {
			log.Printf("Error loading spell %03d: %s\n", i, err)
		}
		spells[i] = spell
	}
}

// resetSpellData resets the fields of the specified SpellData back to their default values.
func resetSpellData(s *SpellData) {
	s.Name = ""
	s.Pic = 0
	s.MPReq = 0
	s.ClassReq = 0
	s.LevelReq = 0
	s.Type = SpellAddHP
	s.Data1 = 0
	s.Data1 = 0
	s.Data3 = 0
}

// SaveSpell saves the data of the spell with the specified ID to the backing file store.
func SaveSpell(id int) {
	if id < 0 || id >= config.MaxSpells {
		return
	}

	err := spellStore.Save(id, spells[id])
	if err != nil {
		log.Printf("Error saving spell %03d: %3\n", id, err)
	}
}

// SaveAllSpells saves the data of all spells to the backing file store.
func SaveAllSpells() {
	for i := 0; i < config.MaxSpells; i++ {
		SaveSpell(i)
	}
}

// GetSpell returns the data of the spell with the specified id
func GetSpell(id int) *SpellData {
	if id < 0 || id >= config.MaxSpells {
		return nil
	}
	return spells[id]
}
