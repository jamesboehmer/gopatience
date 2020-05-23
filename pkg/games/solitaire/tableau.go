package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/games/util"
)

type Tableau struct {
	util.Undoable
	Piles [][]cards.Card
}

func NewTableau(size int, deck *cards.Deck) *Tableau {
	tableau := new(Tableau)
	tableau.Piles = make([][]cards.Card, size)
	// TODO: deal the deck to the tableau
	return tableau
}


func (t *Tableau) Put(cards []cards.Card, pileNum int) error {
	return nil
}

func (t *Tableau) undoPut(pileNum int, numCards int) error {
	return nil
}

func (t *Tableau) Get(pileNum int, cardNum int) ([]cards.Card, error) {
	return []cards.Card{}, nil
}

func (t *Tableau) undoGet(pileNum int, cardStrings []string, reConceal bool) error {
	return nil
}

