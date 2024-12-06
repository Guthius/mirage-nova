package database

import (
	"encoding/gob"
	"log"
	"os"
)

func createFolderIfNotExists(folderName string) {
	info, err := os.Stat(folderName)
	if !os.IsNotExist(err) {
		if !info.IsDir() {
			log.Panicf("Unable to create directory '%s', a file with this name already exists", folderName)
		}
		return
	}
	err = os.MkdirAll(folderName, 0644)
	if err != nil {
		log.Panic(err)
	}
}

func loadFromFile(fileName string, v any) error {
	file, err := os.Open(fileName)

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	defer file.Close()

	decoder := gob.NewDecoder(file)

	return decoder.Decode(v)
}

func saveToFile(fileName string, v any) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := gob.NewEncoder(file)

	return encoder.Encode(v)
}

func IsValidName(name string) bool {
	bytes := []byte(name)
	for _, b := range bytes {
		if b < 48 || b > 122 {
			return false
		}
	}
	return true
}
