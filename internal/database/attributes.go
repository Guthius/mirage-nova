package database

type VitalType int

const (
	VitalHP VitalType = iota
	VitalMP
	VitalSP
)

type Vitals struct {
	HP int
	MP int
	SP int
}

type Stats struct {
	Strength int
	Defense  int
	Speed    int
	Magic    int
}

type EquipmentSlot int

const (
	EquipWeapon EquipmentSlot = iota
	EquipArmor
	EquipHelmet
	EquipShield
)

type Equipment struct {
	Weapon int
	Armor  int
	Helmet int
	Shield int
}
