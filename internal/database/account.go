package database

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Account struct {
	Id           int64
	Name         string
	PasswordHash string
}

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
		    gender INTEGER,
		    class INTEGER,
		    sprite INTEGER,
		    level INTEGER,
		    exp INTEGER,
		    access INTEGER,
		    pk INTEGER
		)`)
}

func AccountExists(accountName string) bool {
	if !IsValidName(accountName) {
		return false
	}

	db, err := openDatabase()
	if err != nil {
		return false
	}

	stmt, err := db.Prepare("SELECT COUNT(id) FROM accounts WHERE name = ?")
	if err != nil {
		return false
	}

	defer stmt.Close()

	row := stmt.QueryRow(accountName)

	var count int64

	err = row.Scan(&count)
	if err != nil {
		return false
	}

	return count == 1
}

func LoadAccount(accountName string) *Account {
	if !IsValidName(accountName) {
		return nil
	}

	db, err := openDatabase()
	if err != nil {
		log.Printf("error loading account '%s' (%s)\n", accountName, err)
		return nil
	}

	stmt, err := db.Prepare("SELECT * FROM accounts WHERE name = ?")
	if err != nil {
		log.Printf("error loading account '%s' (%s)\n", accountName, err)
		return nil
	}

	defer stmt.Close()

	row := stmt.QueryRow(accountName)

	var account Account

	err = row.Scan(&account.Id, &account.Name, &account.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		log.Printf("error loading account '%s' (%s)\n", accountName, err)
		return nil
	}

	return &account
}

func CreateAccount(accountName string, password string) (*Account, bool) {
	if AccountExists(accountName) {
		return nil, false
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, false
	}

	account := &Account{
		Name:         accountName,
		PasswordHash: string(passwordHash),
	}

	db, err := openDatabase()
	if err != nil {
		log.Printf("error creating account '%s' (%s)\n", account.Name, err)
		return nil, false
	}

	stmt, err := db.Prepare("INSERT INTO accounts (name, password_hash) VALUES (?, ?)")
	if err != nil {
		log.Printf("error creating account '%s' (%s)\n", account.Name, err)
		return nil, false
	}

	defer stmt.Close()

	r, err := stmt.Exec(account.Name, account.PasswordHash)
	if err != nil {
		log.Printf("error creating account '%s' (%s)\n", account.Name, err)
		return nil, false
	}

	id, err := r.LastInsertId()
	if err != nil {
		log.Printf("error creating account '%s' (%s)\n", account.Name, err)
		return nil, false
	}

	account.Id = id

	return account, true
}

func (account *Account) Save() bool {
	if account == nil || len(account.Name) == 0 {
		return false
	}

	db, err := openDatabase()
	if err != nil {
		log.Printf("error saving account '%s': %s\n", account.Name, err)
		return false
	}

	stmt, err := db.Prepare("UPDATE accounts SET password_hash = ? WHERE id = ?")
	if err != nil {
		log.Printf("error saving account '%s': %s\n", account.Name, err)
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(account.PasswordHash, account.Id)
	if err != nil {
		log.Printf("error saving account '%s': %s\n", account.Name, err)
		return false
	}

	return true
}

func (account *Account) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(password),
	)
	return err == nil
}
