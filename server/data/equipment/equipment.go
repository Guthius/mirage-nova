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
func (d *Data) Reset() {
	d.Weapon = -1
	d.Armor = -1
	d.Helmet = -1
	d.Shield = -1
}
