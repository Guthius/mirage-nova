package database

import (
	"fmt"
	"log"
)

type SpellType int

const (
	SPELL_TYPE_ADDHP SpellType = iota
	SPELL_TYPE_ADDMP
	SPELL_TYPE_ADDSP
	SPELL_TYPE_SUBHP
	SPELL_TYPE_SUBMP
	SPELL_TYPE_SUBSP
	SPELL_TYPE_GIVEITEM
)

type Spell struct {
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

var Spells [MaxSpells]Spell

func init() {
	loadSpells()
}

func getSpellFilename(spellId int) string {
	return fmt.Sprintf("data/spells/spell%d.gob", spellId+1)
}

func loadSpell(spellId int) {
	Spells[spellId].Clear()

	fileName := getSpellFilename(spellId)

	err := loadFromFile(fileName, &Spells[spellId])
	if err != nil {
		log.Printf("error loading spell '%s' (%s)\n", fileName, err)
	}
}

func loadSpells() {
	createFolderIfNotExists("data/spells")

	for i := 0; i < MaxSpells; i++ {
		loadSpell(i)
	}
}

func SaveSpell(spellId int) {
	if spellId < 0 || spellId >= MaxSpells {
		return
	}

	fileName := getSpellFilename(spellId)

	err := saveToFile(fileName, &Spells[spellId])
	if err != nil {
		log.Printf("error saving spell '%s' (%s)\n", fileName, err)
	}
}

func SaveSpells() {
	for i := 0; i < MaxSpells; i++ {
		SaveSpell(i)
	}
}

func GetSpell(spellId int) *Spell {
	if spellId < 0 || spellId >= MaxSpells {
		return nil
	}
	return &Spells[spellId]
}

func (spell *Spell) Clear() {
	spell.Name = ""
	spell.Pic = 0
	spell.MPReq = 0
	spell.ClassReq = 0
	spell.LevelReq = 0
	spell.Type = SPELL_TYPE_ADDHP
	spell.Data1 = 0
	spell.Data1 = 0
	spell.Data3 = 0
}
