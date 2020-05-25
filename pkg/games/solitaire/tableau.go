package solitaire

import (
	"errors"
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/games/util"
)

type Tableau struct {
	util.Undoable
	Piles [][]*cards.Card
}

func NewTableau(size int, deck *cards.Deck) *Tableau {
	tableau := new(Tableau)
	tableau.Piles = make([][]*cards.Card, size)
	if deck != nil {
		for pileNum, _ := range tableau.Piles {
			tableau.Piles[pileNum] = make([]*cards.Card, 0, 13)
		}
		for startingPileNum, _ := range tableau.Piles {
			for pileNum := startingPileNum; pileNum < len(tableau.Piles); pileNum++ {
				card := deck.Deal()
				if pileNum == startingPileNum {
					card.Reveal()
				}
				tableau.Piles[pileNum] = append(tableau.Piles[pileNum], card)
			}
		}
	}
	return tableau
}

func (t *Tableau) Put(cards []cards.Card, pileNum int) error {
	return nil
}

func (t *Tableau) undoPut(pileNum int, numCards int) error {
	return nil
}

func (t *Tableau) Get(pileNum int, cardNum int) ([]*cards.Card, error) {
	if pileNum < 0 || pileNum > len(t.Piles)-1 {
		return nil, errors.New("invalid pile number")
	}
	if cardNum < 0 || cardNum > len(t.Piles[pileNum])-1 {
		return nil, errors.New("invalid card number")
	}
	if !t.Piles[pileNum][cardNum].Revealed {
		return nil, errors.New("card is concealed")
	}

	cards := t.Piles[pileNum][cardNum:]
	t.Piles[pileNum] = t.Piles[pileNum][:cardNum]
	revealed := t.reveal(pileNum)
	t.UndoStack = append(t.UndoStack, util.UndoAction{
		Function: t.undoGet, Args: []interface{}{pileNum, cards, revealed},
	})
	return cards, nil
}

func (t *Tableau) undoGet(args ...interface{}) error {
	pileNum := args[0].(int)
	cards := args[1].([]*cards.Card)
	reConceal := args[2].(bool)
	if reConceal {
		t.conceal(pileNum)
	}
	t.Piles[pileNum] = append(t.Piles[pileNum], cards...)
	return nil
}

func (t *Tableau) reveal(pileNum int) bool {
	pile := t.Piles[pileNum]
	if len(pile) > 0 && !pile[len(pile)-1].Revealed {
		pile[len(pile)-1].Reveal()
		return true
	}
	return false

}

func (t *Tableau) conceal(pileNum int) bool {
	pile := t.Piles[pileNum]
	if len(pile) > 0 && pile[len(pile)-1].Revealed {
		pile[len(pile)-1].Conceal()
		return true
	}
	return false

}
