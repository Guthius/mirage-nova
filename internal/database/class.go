package database

type Class struct {
	Name   string
	Sprite int
	Stats  Stats
}

var Classes = loadClasses()

func loadClasses() []Class {
	// TODO: Implement me
	classes := make([]Class, 0)
	return classes
}

func (class *Class) GetMaxVital(Vital VitalType) int {
	switch Vital {
	case VitalHP:
		return (1 + (class.Stats.Strength / 2) + class.Stats.Strength) * 2
	case VitalMP:
		return (1 + (class.Stats.Magic / 2) + class.Stats.Magic) * 2
	case VitalSP:
		return (1 + (class.Stats.Speed / 2) + class.Stats.Speed) * 2
	}
	return 0
}
