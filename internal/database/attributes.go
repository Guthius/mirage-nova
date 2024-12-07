package database

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
