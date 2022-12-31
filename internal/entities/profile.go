package entities

const (
	MaleGender   = 1
	FemaleGender = 2
)

type Profile struct {
	ID        int64
	Name      string
	Surname   string
	City      string
	Interests []string
	Age       int8
	Gender    int8
}
