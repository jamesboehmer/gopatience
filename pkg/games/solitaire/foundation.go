package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"github.com/jamesboehmer/gopatience/pkg/games/util"
)

type Foundation struct {
	util.Undoable
	Piles map[suit.Suit][]cards.Card
}

func NewFoundation(suits []suit.Suit) *Foundation {
	foundation := new(Foundation)
	foundation.Piles = make(map[suit.Suit][]cards.Card, len(suits))
	for _, suit := range suits {
		foundation.Piles[suit] = make([]cards.Card, 0, 13)
	}
	return foundation
}

func (f *Foundation) undoPut(suit suit.Suit) error {
	return nil
}

func (f *Foundation) Put(card cards.Card) error {
	return nil
}

func (f *Foundation) undoGet(card cards.Card) error {
	return nil
}

func (f *Foundation) Get(suit suit.Suit) (*cards.Card, error) {
	return nil, nil
}

func (f *Foundation) IsFull() bool {
	return false
}
