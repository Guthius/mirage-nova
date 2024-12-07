package vitals

type Type int

const (
	HP Type = iota
	MP
	SP
)

type Data struct {
	HP int
	MP int
	SP int
}
