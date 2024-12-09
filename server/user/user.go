package user

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/guthius/mirage-nova/server/utils"

	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	Id           int64
	Name         string
	PasswordHash string
}

func init() {
	db, err := openDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		name TEXT UNIQUE COLLATE NOCASE,
    		password_hash TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_from_ip TEXT
    	)`)

	if err != nil {
		log.Panic(err)
	}
}

// Exists checks if an account with the specified name exists in the database.
func Exists(accountName string) bool {
	if !utils.IsValidName(accountName) {
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

// Load loads the account with the specified name from the database.
func Load(accountName string) *Account {
	if !utils.IsValidName(accountName) {
		return nil
	}

	db, err := openDatabase()
	if err != nil {
		log.Printf("error loading account '%s' (%s)\n", accountName, err)
		return nil
	}

	stmt, err := db.Prepare("SELECT id, name, password_hash FROM accounts WHERE name = ?")
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

// Create creates a new account with the specified name and password.
func Create(accountName string, password string, createdFromIp string) (*Account, bool) {
	if Exists(accountName) {
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

	stmt, err := db.Prepare("INSERT INTO accounts (name, password_hash, created_from_ip) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("error creating account '%s' (%s)\n", account.Name, err)
		return nil, false
	}

	defer stmt.Close()

	r, err := stmt.Exec(account.Name, account.PasswordHash, createdFromIp)
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

// Save saves the account to the database.
func (account *Account) Save() bool {
	if account == nil || len(account.Name) == 0 {
		return false
	}

	db, err := openDatabase()
	if err != nil {
		log.Printf("error saving account '%s' (%s)\n", account.Name, err)
		return false
	}

	stmt, err := db.Prepare("UPDATE accounts SET password_hash = ? WHERE id = ?")
	if err != nil {
		log.Printf("error saving account '%s' (%s)\n", account.Name, err)
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(account.PasswordHash, account.Id)
	if err != nil {
		log.Printf("error saving account '%s' (%s)\n", account.Name, err)
		return false
	}

	return true
}

// IsPasswordCorrect checks if the specified password is correct for the account.
func (account *Account) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(password),
	)
	return err == nil
}

// openDatabase opens the SQLite database and returns a database handle.
func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/accounts.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
