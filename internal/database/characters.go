package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type CharacterGender int

const (
	GenderMale CharacterGender = iota
	GenderFemale
)

type AccessLevel int

const (
	AccessNone AccessLevel = iota
	AccessMonitor
	AccessMapper
	AccessDeveloper
	AccessCreator
)

type InventorySlot struct {
	Item  int
	Value int
	Dur   int
}

type Character struct {
	Id          int64
	AccountId   int64
	Name        string
	Gender      CharacterGender
	Class       int
	Sprite      int
	Level       int
	Exp         int
	Access      AccessLevel
	PK          bool
	Guild       string
	GuildAccess int
	Vitals      Vitals
	Stats       Stats
	Points      int
	Equipment   Equipment
	Inv         [MaxInventory]InventorySlot
	Spells      [MaxCharacterSpells]int
	Map         int
	X           int
	Y           int
	Dir         Direction
}

func CharacterExists(characterName string) bool {
	if !IsValidName(characterName) {
		return false
	}

	db, err := openDatabase()
	if err != nil {
		return false
	}

	stmt, err := db.Prepare("SELECT COUNT(id) FROM characters WHERE name = ?")
	if err != nil {
		return false
	}

	defer stmt.Close()

	row := stmt.QueryRow(characterName)

	var count int64

	err = row.Scan(&count)
	if err != nil {
		return false
	}

	return count == 1
}

func decodeInventoryFromJson(inventoryJson string) [MaxInventory]InventorySlot {
	var slots [MaxInventory]InventorySlot

	for i := 0; i < MaxInventory; i++ {
		slots[i].Item = -1
		slots[i].Value = 0
		slots[i].Dur = 0
	}

	err := json.Unmarshal([]byte(inventoryJson), &slots)
	if err != nil {
		log.Printf("error decoding inventory (%s)\n", err)
		return slots
	}

	return slots
}

func decodeSpellsFromJson(spellsJson string) [MaxCharacterSpells]int {
	var spells [MaxCharacterSpells]int

	err := json.Unmarshal([]byte(spellsJson), &spells)
	if err != nil {
		log.Printf("error decoding spells (%s)\n", err)
		return spells
	}

	return spells
}

func LoadCharactersForAccount(accountId int64) []Character {
	characters := make([]Character, 0)

	db, err := openDatabase()
	if err != nil {
		log.Printf("error loading characters for account %d (%s)\n", accountId, err)
		return characters
	}

	stmt, err := db.Prepare("SELECT * FROM characters WHERE account_id = ?")
	if err != nil {
		log.Printf("error loading characters for account %d (%s)\n", accountId, err)
		return characters
	}

	defer stmt.Close()

	rows, err := stmt.Query(accountId)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error loading characters for account %d (%s)\n", accountId, err)
		}
		return characters
	}

	defer rows.Close()

	var characterInventory string
	var characterSpells string

	for rows.Next() {
		var character Character

		err := rows.Scan(
			&character.Id,
			&character.AccountId,
			&character.Name,
			&character.Gender,
			&character.Class,
			&character.Sprite,
			&character.Level,
			&character.Exp,
			&character.Access,
			&character.PK,
			&character.Guild,
			&character.GuildAccess,
			&character.Vitals.HP,
			&character.Vitals.MP,
			&character.Vitals.SP,
			&character.Stats.Strength,
			&character.Stats.Defense,
			&character.Stats.Speed,
			&character.Stats.Magic,
			&character.Equipment.Weapon,
			&character.Equipment.Armor,
			&character.Equipment.Helmet,
			&character.Equipment.Shield,
			&characterInventory,
			&characterSpells,
			&character.Map,
			&character.X,
			&character.Y,
			&character.Dir)

		character.Inv = decodeInventoryFromJson(characterInventory)
		character.Spells = decodeSpellsFromJson(characterSpells)

		if err != nil {
			log.Printf("error loading characters for account %d (%s)\n", accountId, err)
			continue
		}

		characters = append(characters, character)
	}

	return characters
}

func (c *Character) Clear() {
	c.Id = 0
	c.Name = ""
	c.Gender = GenderMale
	c.Class = 0
	c.Sprite = 0
	c.Level = 1
	c.Exp = 0
	c.Access = AccessNone
	c.PK = false
	c.Guild = ""
	c.GuildAccess = 0
	c.Vitals = Vitals{}
	c.Stats = Stats{}
	c.Equipment.Weapon = -1
	c.Equipment.Armor = -1
	c.Equipment.Helmet = -1
	c.Equipment.Shield = -1
	c.Map = 0
	c.X = 0
	c.Y = 0
	c.Dir = Down

	c.ClearInventory()
	c.ClearSpells()
}

func (c *Character) ClearInventory() {
	for i := 0; i < MaxInventory; i++ {
		c.Inv[i].Item = -1
		c.Inv[i].Value = 0
		c.Inv[i].Dur = 0
	}
}

func (c *Character) ClearSpells() {
	for i := 0; i < MaxCharacterSpells; i++ {
		c.Spells[i] = -1
	}
}

func encodeInventoryAsJson(inventory [MaxInventory]InventorySlot) string {
	bytes, err := json.Marshal(inventory)
	if err != nil {
		log.Printf("error encoding inventory (%s)\n", err)
		return ""
	}

	return string(bytes)
}

func encodeSpellsAsJson(spells [MaxCharacterSpells]int) string {
	bytes, err := json.Marshal(spells)
	if err != nil {
		log.Printf("error encoding spells (%s)\n", err)
		return ""
	}

	return string(bytes)
}

func (c *Character) Save() bool {
	if c == nil || c.Id == 0 {
		return false
	}

	db, err := openDatabase()
	if err != nil {
		log.Fatal(err)
		return false
	}

	stmt, err := db.Prepare(`
		UPDATE characters 
		SET 
		    gender = ?, 
		    class = ?, 
		    sprite = ?, 
		    level = ?,
		    exp = ?,
		    access = ?, 
		    pk = ?,
		    guild = ?,
		    guild_access = ?,
		    vital_hp = ?,
		    vital_mp = ?,
		    vital_sp = ?,
		    stat_strength = ?,
		    stat_defense = ?,
		    stat_speed = ?,
		    stat_magic = ?,
		    equip_weapon = ?,
		    equip_armor = ?,
		    equip_helmet = ?,
		    equip_shield = ?,
		    inventory = ?,
		    spells = ?,
		    map = ?,
		    x = ?,
		    y = ?,
		    dir = ?
		WHERE id = ?`)

	if err != nil {
		return false
	}

	defer stmt.Close()

	characterInventory := encodeInventoryAsJson(c.Inv)
	characterSpells := encodeSpellsAsJson(c.Spells)

	_, err = stmt.Exec(
		c.Gender,
		c.Class,
		c.Sprite,
		c.Level,
		c.Exp,
		c.Access,
		c.PK,
		c.Name,
		c.Guild,
		c.GuildAccess,
		c.Vitals.HP,
		c.Vitals.MP,
		c.Vitals.SP,
		c.Stats.Strength,
		c.Stats.Defense,
		c.Stats.Speed,
		c.Stats.Magic,
		characterInventory,
		characterSpells,
		c.Map,
		c.X,
		c.Y,
		c.Dir)

	if err != nil {
		return false
	}

	return true
}

func (c *Character) Delete() bool {
	if c == nil || c.Id == 0 {
		return false
	}

	db, err := openDatabase()
	if err != nil {
		log.Printf("error deleting character %d (%s)\n", c.Id, err)
		return false
	}

	stmt, err := db.Prepare("DELETE FROM characters WHERE id = ?")
	if err != nil {
		log.Printf("error deleting character %d (%s)\n", c.Id, err)
		return false
	}

	_, err = stmt.Exec(c.Id)
	if err != nil {
		log.Printf("error deleting character %d (%s)\n", c.Id, err)
		return false
	}

	c.Clear()

	return true
}

func CreateCharacter(accountId int64, name string, gender CharacterGender, classId int) (*Character, bool) {
	if CharacterExists(name) {
		return nil, false
	}

	class := &Classes[classId]

	character := &Character{
		AccountId: accountId,
		Name:      name,
		Gender:    gender,
		Class:     classId,
		Sprite:    class.Sprite,
		Level:     1,
		Exp:       0,
		Access:    AccessNone,
		PK:        false,

		Vitals: Vitals{
			HP: class.GetMaxVital(VitalHP, class.Stats.Strength),
			MP: class.GetMaxVital(VitalMP, class.Stats.Magic),
			SP: class.GetMaxVital(VitalSP, class.Stats.Speed),
		},

		Stats: Stats{
			Strength: class.Stats.Strength,
			Defense:  class.Stats.Defense,
			Speed:    class.Stats.Speed,
			Magic:    class.Stats.Magic,
		},
	}
	
	character.ClearInventory()
	character.ClearSpells()

	db, err := openDatabase()
	if err != nil {
		log.Printf("error creating character '%s' (%s)\n", name, err)
		return nil, false
	}

	stmt, err := db.Prepare(
		`INSERT INTO characters 
    	(account_id, name, gender, class, sprite, level, exp, access, pk,
    	 vital_hp, vital_mp, vital_sp, 
    	 stat_strength, stat_defense, stat_speed, stat_magic,
    	 inventory, spells, map, x, y, dir) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf("error creating character '%s' (%s)\n", name, err)
		return nil, false
	}

	defer stmt.Close()

	characterInventory := encodeInventoryAsJson(character.Inv)
	characterSpells := encodeSpellsAsJson(character.Spells)

	result, err := stmt.Exec(
		character.AccountId,
		character.Name,
		character.Gender,
		character.Class,
		character.Sprite,
		character.Level,
		character.Exp,
		character.Access,
		character.PK,
		character.Vitals.HP,
		character.Vitals.MP,
		character.Vitals.SP,
		character.Stats.Strength,
		character.Stats.Defense,
		character.Stats.Speed,
		character.Stats.Magic,
		characterInventory,
		characterSpells,
		character.Map,
		character.X,
		character.Y,
		character.Dir)

	if err != nil {
		log.Printf("error creating character '%s' (%s)\n", name, err)
		return nil, false
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("error creating character '%s' (%s)\n", name, err)
		return nil, false
	}

	character.Id = id

	return character, true
}
