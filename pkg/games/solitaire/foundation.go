package solitaire

import (
	"errors"
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/pip"
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

func (f *Foundation) undoPut(args ...interface{}) error {
	suit := args[0].(suit.Suit)
	f.Piles[suit] = f.Piles[suit][:len(f.Piles[suit])-1]
	return nil
}

func (f *Foundation) Put(card cards.Card) error {
	if !card.Revealed {
		return errors.New("foundation cards must be revealed")
	}
	pile := f.Piles[card.Suit]
	if card.Pip == pip.Ace {
		f.Piles[card.Suit] = append(pile, card)
		f.UndoStack = append(f.UndoStack, util.UndoAction{Function: f.undoPut, Args: []interface{}{card.Suit}})
		return nil
	} else if len(pile) == 0 {
		return errors.New("the first card on a foundation pile must be an ace")
	} else {
		topCard := pile[len(pile)-1]
		if PipValue[card.Pip] != PipValue[topCard.Pip]+1 {
			return errors.New("foundation cards must be built sequentially by suit")
		}
		f.Piles[card.Suit] = append(pile, card)
		f.UndoStack = append(f.UndoStack, util.UndoAction{Function: f.undoPut, Args: []interface{}{card.Suit}})
		return nil
	}
}

func (f *Foundation) undoGet(args ...interface{}) error {
	card := args[0].(cards.Card)
	f.Piles[card.Suit] = append(f.Piles[card.Suit], card)
	return nil
}

func (f *Foundation) Get(suit suit.Suit) (*cards.Card, error) {
	pile, found := f.Piles[suit]
	if !found {
		return nil, errors.New("no such suit")
	}
	if len(pile) == 0 {
		return nil, errors.New("pile is empty")
	}
	topCard := pile[len(pile)-1]
	f.Piles[suit] = pile[:len(pile)-1]
	f.UndoStack = append(f.UndoStack, util.UndoAction{Function: f.undoGet, Args: []interface{}{topCard}})
	return &topCard, nil
}

func (f *Foundation) IsFull() bool {
	return false
}
