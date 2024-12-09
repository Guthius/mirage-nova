package character

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/guthius/mirage-nova/server/common"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/data/equipment"
	"github.com/guthius/mirage-nova/server/data/stats"
	"github.com/guthius/mirage-nova/server/data/vitals"
	"github.com/guthius/mirage-nova/server/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Gender int

const (
	GenderMale Gender = iota
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
	Gender      Gender
	Class       int
	Sprite      int
	Level       int
	Exp         int
	Access      AccessLevel
	PK          bool
	Guild       string
	GuildAccess int
	Vitals      vitals.Data
	Stats       stats.Data
	Points      int
	Equipment   equipment.Data
	Inv         [config.MaxInventory]InventorySlot
	Spells      [config.MaxCharacterSpells]int
	Room        int
	X           int
	Y           int
	Dir         common.Direction
}

func init() {
	db, err := openDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS characters (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    account_id INTEGER,
		    name TEXT UNIQUE COLLATE NOCASE,
		    gender INTEGER NOT NULL DEFAULT 0,
		    class INTEGER NOT NULL,
		    sprite INTEGER NOT NULL DEFAULT 0,
		    level INTEGER NOT NULL DEFAULT 1,
		    exp INTEGER NOT NULL DEFAULT 0,
		    access INTEGER NOT NULL DEFAULT 0,
		    pk INTEGER NOT NULL DEFAULT 0,
		    guild TEXT NOT NULL COLLATE NOCASE DEFAULT '',
		    guild_access INTEGER NOT NULL DEFAULT 0,
		    vital_hp INTEGER NOT NULL,
		    vital_mp INTEGER NOT NULL,
		    vital_sp INTEGER NOT NULL,
		    stat_strength INTEGER NOT NULL,
		    stat_defense INTEGER NOT NULL,
		    stat_speed INTEGER NOT NULL,
		    stat_magic INTEGER NOT NULL,
		    equip_weapon INTEGER NOT NULL DEFAULT -1,
		    equip_armor INTEGER NOT NULL DEFAULT -1,
		    equip_helmet INTEGER NOT NULL DEFAULT -1,
		    equip_shield INTEGER NOT NULL DEFAULT -1,
		    inventory TEXT NOT NULL DEFAULT '',
		    spells TEXT NOT NULL DEFAULT '',
		    room INTEGER NOT NULL DEFAULT 0, 
		    x INTEGER NOT NULL DEFAULT 0, 
		    y INTEGER NOT NULL DEFAULT 0,
		    dir INTEGER NOT NULL DEFAULT 0
		)`)

	if err != nil {
		log.Panic(err)
	}
}

func Exists(characterName string) bool {
	if !utils.IsValidName(characterName) {
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

func decodeInventoryFromJson(inventoryJson string) [config.MaxInventory]InventorySlot {
	var slots [config.MaxInventory]InventorySlot

	for i := 0; i < config.MaxInventory; i++ {
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

func decodeSpellsFromJson(spellsJson string) [config.MaxCharacterSpells]int {
	var spells [config.MaxCharacterSpells]int

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
			&character.Room,
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
	c.Vitals = vitals.Data{}
	c.Stats = stats.Data{}
	c.Equipment.Weapon = -1
	c.Equipment.Armor = -1
	c.Equipment.Helmet = -1
	c.Equipment.Shield = -1
	c.Room = 0
	c.X = 0
	c.Y = 0
	c.Dir = common.DirDown

	c.ClearInventory()
	c.ClearSpells()
}

func (c *Character) ClearInventory() {
	for i := 0; i < config.MaxInventory; i++ {
		c.Inv[i].Item = -1
		c.Inv[i].Value = 0
		c.Inv[i].Dur = 0
	}
}

func (c *Character) ClearSpells() {
	for i := 0; i < config.MaxCharacterSpells; i++ {
		c.Spells[i] = -1
	}
}

func encodeInventoryAsJson(inventory [config.MaxInventory]InventorySlot) string {
	bytes, err := json.Marshal(inventory)
	if err != nil {
		log.Printf("error encoding inventory (%s)\n", err)
		return ""
	}

	return string(bytes)
}

func encodeSpellsAsJson(spells [config.MaxCharacterSpells]int) string {
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
		    room = ?,
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
		c.Room,
		c.X,
		c.Y,
		c.Dir)

	return err == nil
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

func CreateCharacter(accountId int64, name string, gender Gender, classId int) (*Character, bool) {
	if Exists(name) {
		return nil, false
	}

	class := data.GetClass(classId)
	if class == nil {
		return nil, false
	}

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
		Room:      config.StartRoom,
		X:         config.StartX,
		Y:         config.StartY,
		Dir:       common.DirDown,

		Vitals: vitals.Data{
			HP: class.GetMaxVital(vitals.HP, class.Stats.Strength),
			MP: class.GetMaxVital(vitals.MP, class.Stats.Magic),
			SP: class.GetMaxVital(vitals.SP, class.Stats.Speed),
		},

		Stats: stats.Data{
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
		character.Room,
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

// openDatabase opens the SQLite database and returns a database handle.
func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/characters.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
