package equipment

type Slot int

const (
	Weapon Slot = iota
	Armor
	Helmet
	Shield
)

type Data struct {
	Weapon int
	Armor  int
	Helmet int
	Shield int
}

// Reset resets the equipment to default values.
func (e *Data) Reset() {
	e.Weapon = -1
	e.Armor = -1
	e.Helmet = -1
	e.Shield = -1
}
