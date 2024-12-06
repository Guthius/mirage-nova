package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/accounts.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Create() {
	db, err := openDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	Exec := func(query string) {
		_, err := db.Exec(query)
		if err != nil {
			log.Panic(err)
		}
	}

	Exec(`CREATE TABLE IF NOT EXISTS accounts (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		name TEXT UNIQUE COLLATE NOCASE,
    		password_hash TEXT
    	)`)

	Exec(`CREATE TABLE IF NOT EXISTS characters (
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
		    map INTEGER NOT NULL DEFAULT 0, 
		    x INTEGER NOT NULL DEFAULT 0, 
		    y INTEGER NOT NULL DEFAULT 0,
		    dir INTEGER NOT NULL DEFAULT 0
		)`)
}

/*
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
	Dir         Direction*/
