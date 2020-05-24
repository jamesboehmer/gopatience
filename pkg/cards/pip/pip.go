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
	Jack  Pip = "J"
	Queen Pip = "Q"
	King  Pip = "K"
)

var Pips = map[string]Pip{
	"A": Ace, "2": Two, "3": Three, "4": Four, "5": Five, "6": Six, "7": Seven,
	"8": Eight, "9": Nine, "10": Ten, "J": Jack, "Q": Queen, "K": King,
}

func (p Pip) IsFace() bool {
	return p == Jack || p == Queen || p == King
}
