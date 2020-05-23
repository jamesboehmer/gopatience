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

func (suit Suit) Color() Color {
	if suit == Diamonds || suit == Hearts {
		return Red
	} else {
		return Black
	}
}
