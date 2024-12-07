package data

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/guthius/mirage-nova/server/data/stats"
	"github.com/guthius/mirage-nova/server/data/vitals"
)

type ClassData struct {
	Name   string
	Sprite int
	Stats  stats.Data
}

var classes []ClassData

func init() {
	file, err := os.OpenFile("data/classes.json", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &classes)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Loaded %d classes\n", len(classes))
}

// GetClassCount returns the number of classes available.
func GetClassCount() int {
	return len(classes)
}

// GetClass returns the class with the specified id.
func GetClass(id int) *ClassData {
	if id < 0 || id >= len(classes) {
		return nil
	}
	return &classes[id]
}

// GetMaxVital returns the maximum value of the specified vital type.
func (c *ClassData) GetMaxVital(vital vitals.Type, stat int) int {
	switch vital {
	case vitals.HP:
		return (1 + (stat / 2) + c.Stats.Strength) * 2
	case vitals.MP:
		return (1 + (stat / 2) + c.Stats.Magic) * 2
	case vitals.SP:
		return (1 + (stat / 2) + c.Stats.Speed) * 2
	}
	return 0
}
