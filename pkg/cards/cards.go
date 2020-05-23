package cards

import (
	"github.com/jamesboehmer/gopatience/pkg/cards/pip"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
)

type Deck struct {
	NumDecks  int
	NumJokers int
	Cards     []Card
}

type Card struct {
	Pip      pip.Pip
	Suit     suit.Suit
	Revealed bool
}

func (deck *Deck) Shuffle() *Deck {
	// TODO: shuffle
	return deck
}

func (deck *Deck) Remaining() int {
	return len(deck.Cards)
}

func (card *Card) Reveal() *Card {
	card.Revealed = true
	return card
}

func (card *Card) Conceal() *Card {
	card.Revealed = false
	return card
}

func NewDeck(numDecks int, numJokers int) *Deck {
	deck := new(Deck)
	// TODO: create cards
	return deck
}

func ParseCard(cardString string) (*Card, error) {
	return nil, nil
}
