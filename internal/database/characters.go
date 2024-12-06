package database

import (
	"database/sql"
	"errors"
	"log"
)

type CharacterGender int

const (
	GenderMale CharacterGender = iota
	GenderFemale
)

type AccessLevel int

const (
	AccessNone AccessLevel = iota
)

type Character struct {
	// General
	Id        int64
	AccountId int64
	Name      string
	Gender    CharacterGender
	Class     int
	Sprite    int
	Level     int
	Exp       int
	Access    AccessLevel
	PK        bool
}

/*

Public Type PlayerRec
    ' General
    Name As String * NAME_LENGTH
    Sex As Byte
    Class As Byte
    Sprite As Integer
    Level As Byte
    Exp As Long
    Access As Byte
    PK As Byte

    Guild As String
    GuildAccess As Long

    ' Vitals
    Vital(1 To Vitals.Vital_Count - 1) As Long

    ' Stats
    Stat(1 To Stats.Stat_Count - 1) As Byte
    POINTS As Byte

    ' Worn equipment
    Equipment(1 To Equipment.Equipment_Count - 1) As Byte

    ' Inventory
    Inv(1 To MAX_INV) As PlayerInvRec
    Spell(1 To MAX_PLAYER_SPELLS) As Byte

    ' Position
    Map As Integer
    X As Byte
    Y As Byte
    Dir As Byte
End Type
*/

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
			&character.PK)

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
		    pk = ?
		WHERE id = ?`)

	if err != nil {
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(c.Gender, c.Class, c.Sprite, c.Level, c.Exp, c.Access, c.PK, c.Name)
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

func CreateCharacter(accountId int64, name string, gender CharacterGender, class int) (*Character, bool) {
	if CharacterExists(name) {
		return nil, false
	}

	character := &Character{
		AccountId: accountId,
		Name:      name,
		Gender:    gender,
		Class:     class,
		Sprite:    Classes[class].Sprite,
		Level:     1,
		Exp:       0,
		Access:    AccessNone,
		PK:        false,
	}

	db, err := openDatabase()
	if err != nil {
		log.Printf("error creating character '%s' (%s)\n", name, err)
		return nil, false
	}

	stmt, err := db.Prepare(
		`INSERT INTO characters (account_id, name, gender, class, sprite, level, exp, access, pk) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Printf("error creating character '%s' (%s)\n", name, err)
		return nil, false
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		character.AccountId,
		character.Name,
		character.Gender,
		character.Class,
		character.Sprite,
		character.Level,
		character.Exp,
		character.Access,
		character.PK)

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
