package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"testing"
)

func TestNewTableau(t *testing.T) {
	tableau := NewTableau(8, nil)
	if len(tableau.Piles) != 8 {
		t.Error("Tableau size is wrong")
	}
	for pileNum, _ := range tableau.Piles {
		if len(tableau.Piles[pileNum]) != 0 {
			t.Error("Tableau without deck should have no cards")
		}
	}
	deck := cards.NewDeck(1, 0).Shuffle()
	tableau = NewTableau(7, deck)
	if len(tableau.Piles) != 7 {
		t.Error("Tableau should have 7 piles")

	}
	for pileNum, pile := range tableau.Piles {
		if len(pile) != pileNum + 1 {
			t.Error("Pile length is incorrect")
		}
		for cardNum, card := range pile {
			if cardNum == pileNum {
				if !card.Revealed {
					t.Error("Top card should be revealed.")
				}
			} else {
				if card.Revealed {
					t.Error("Non-top cards should be concealed")
				}
			}
		}
	}
}
