package suit

type Suit string
type Color int

const (
	Red   Color = iota
	Black Color = iota

	Spades   Suit = "♠"
	Hearts   Suit = "♥"
	Diamonds Suit = "♦"
	Clubs    Suit = "♣"
)

var Suits []Suit = []Suit{Spades, Hearts, Diamonds, Clubs}

func (suit Suit) Color() Color {
	if suit == Diamonds || suit == Hearts {
		return Red
	} else {
		return Black
	}
}
