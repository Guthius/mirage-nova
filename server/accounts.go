package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

type Account struct {
	Name         string
	PasswordHash string
}

func IsValidAccountName(accountName string) bool {
	bytes := []byte(accountName)
	for _, b := range bytes {
		if b < 48 || b > 122 {
			return false
		}
	}
	return true
}

func GetAccountFileName(accountName string) string {
	return fmt.Sprintf("accounts/%s.json", strings.ToLower(accountName))
}

func AccountExists(accountName string) bool {
	if !IsValidAccountName(accountName) {
		return false
	}

	fileName := GetAccountFileName(accountName)

	info, err := os.Stat(fileName)
	if err == nil {
		return !info.IsDir()
	}

	return false
}

func LoadAccount(accountName string) *Account {
	if !IsValidAccountName(accountName) {
		return nil
	}

	fileName := GetAccountFileName(accountName)
	file, err := os.Open(fileName)

	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error opening '%s': %s\n", fileName, err)
		}
		return nil
	}

	defer file.Close()

	var account Account

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&account)
	if err != nil {
		log.Printf("Error loaded '%s': %s\n", fileName, err)
	}

	return nil
}

func SaveAccount(account *Account) {
	if account == nil || len(account.Name) == 0 {
		return
	}

	fileName := GetAccountFileName(account.Name)
	file, err := os.Create(fileName)

	if err != nil {
		log.Printf("Error creating '%s': %s\n", fileName, err)
		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(account)
	if err != nil {
		log.Printf("Error encoding '%s': %s\n", fileName, err)
	}
}

func CreateAccount(accountName string, password string) (*Account, error) {
	if AccountExists(accountName) {
		return nil, fmt.Errorf("account '%s' already exists", accountName)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Name:         accountName,
		PasswordHash: string(passwordHash),
	}

	SaveAccount(account)

	return account, nil
}

func (account *Account) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(account.PasswordHash),
		[]byte(password),
	)
	return err == nil
}
