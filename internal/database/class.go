package database

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Class struct {
	Name   string
	Sprite int
	Stats  Stats
}

var Classes []Class

func init() {
	Classes = loadClasses()
}

func loadClasses() []Class {
	file, err := os.OpenFile("data/classes.json", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var classes []Class

	err = json.Unmarshal(bytes, &classes)
	if err != nil {
		log.Fatal(err)
	}

	return classes
}

func (class *Class) GetMaxVital(Vital VitalType, stat int) int {
	switch Vital {
	case VitalHP:
		return (1 + (stat / 2) + class.Stats.Strength) * 2
	case VitalMP:
		return (1 + (stat / 2) + class.Stats.Magic) * 2
	case VitalSP:
		return (1 + (stat / 2) + class.Stats.Speed) * 2
	}
	return 0
}
