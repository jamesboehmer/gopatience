package cards

import (
	"errors"
	"github.com/jamesboehmer/gopatience/pkg/cards/pip"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"math/rand"
	"strings"
	"time"
)

type Deck struct {
	NumDecks   int
	NumJokers  int
	Cards      []Card
	IsShuffled bool
}

type Card struct {
	Pip      pip.Pip
	Suit     suit.Suit
	Revealed bool
}

func (deck *Deck) Shuffle() *Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck.Cards), func(i, j int) { deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i] })
	deck.IsShuffled = true
	return deck
}

func (deck *Deck) Deal() *Card {
	if len(deck.Cards) == 0 {
		return nil
	}
	card := deck.Cards[0]
	deck.Cards = deck.Cards[1:]
	return &card

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

func (card *Card) IsFace() bool {
	if card.Pip != "" {
		return card.Pip.IsFace()
	}
	return false
}

func (card *Card) String() string {
	buffer := strings.Builder{}
	if !card.Revealed {
		buffer.WriteString("|")
	}
	if card.Pip == "" {
		buffer.WriteString("*")
	} else {
		buffer.WriteString(string(card.Pip))
		buffer.WriteString(string(card.Suit))
	}
	return buffer.String()
}

func NewDeck(numDecks int, numJokers int) *Deck {
	deck := new(Deck)

	deck.NumDecks = numDecks
	deck.NumJokers = numJokers
	var cards []Card

	for deckNum := 0; deckNum < numDecks; deckNum++ {
		for _, suit := range suit.Suits {
			for _, pip := range pip.Pips {
				cards = append(cards, Card{pip, suit, false})
			}
		}
		for jokerNum := 0; jokerNum < numJokers; jokerNum++ {
			cards = append(cards, Card{})
		}
	}
	deck.IsShuffled = false
	deck.Cards = cards

	return deck
}

func ParseCard(cardString string) (*Card, error) {
	revealed := true
	runes := []rune(cardString)
	if runes[0] == '|' {
		revealed = false
		runes = runes[1:]
	}
	if runes[0] == '*' {
		return &Card{"", "", revealed}, nil
	}
	suit, found := suit.Suits[string(runes[len(runes)-1:][0])]
	if !found {
		return nil, errors.New("invalid suit")
	}
	runes = runes[:len(runes)-1]
	pip, found := pip.Pips[string(runes)]
	if !found {
		return nil, errors.New("invalid pip")
	}
	return &Card{pip, suit, revealed}, nil
}
