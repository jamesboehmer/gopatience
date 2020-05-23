package pip

type Pip string

const (
	Ace   Pip = "A"
	Two   Pip = "2"
	Three Pip = "3"
	Four  Pip = "4"
	Five  Pip = "5"
	Six   Pip = "6"
	Seven Pip = "7"
	Eight Pip = "8"
	Nine  Pip = "9"
	Ten   Pip = "10"
	Jack  Pip = "Jack"
	Queen Pip = "Queen"
	King  Pip = "King"
)

func (p Pip) IsFace() bool {
	return p == Jack || p == Queen || p == King
}
