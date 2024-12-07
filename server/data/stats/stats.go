package stats

type Type int

const (
	Strength Type = iota
	Defense
	Speed
	Magic
)

type Data struct {
	Strength int
	Defense  int
	Speed    int
	Magic    int
}

// Get returns the value of the specified stat.
func (d *Data) Get(stat Type) int {
	switch stat {
	case Strength:
		return d.Strength
	case Defense:
		return d.Defense
	case Speed:
		return d.Speed
	case Magic:
		return d.Magic
	}
	return 0
}

// Reset resets all stats back to zero.
func (d *Data) Reset() {
	d.Strength = 0
	d.Defense = 0
	d.Speed = 0
	d.Magic = 0
}
