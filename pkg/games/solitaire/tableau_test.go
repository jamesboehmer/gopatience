package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/pip"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
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
		if len(pile) != pileNum+1 {
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

func TestTableau_GetEmptyPile(t *testing.T) {
	tableau := NewTableau(7, nil)
	_, err := tableau.Get(0, 0)
	if err == nil {
		t.Error("Empty tableau get should return an error")
	}
}

func TestTableau_GetInvalidPile(t *testing.T) {
	tableau := NewTableau(7, nil)
	_, err := tableau.Get(-1, 0)
	if err == nil {
		t.Error("Negative tableau pileNum get should return an error")
	}
	_, err = tableau.Get(7, 0)
	if err == nil {
		t.Error("Invalid tableau pileNum get should return an error")
	}
}

func TestTableau_GetInvalidCardNum(t *testing.T) {
	tableau := NewTableau(7, cards.NewDeck(1, 0).Shuffle())
	_, err := tableau.Get(0, 1)
	if err == nil {
		t.Error("Invalid cardNum get should return an error")
	}
}

func TestTableau_GetConcealedCard(t *testing.T) {
	tableau := NewTableau(7, cards.NewDeck(1, 0).Shuffle())
	_, err := tableau.Get(1, 0)
	if err == nil {
		t.Error("Concealed tableau get should return an error")
	}
}

func TestTableau_GetValidSliceWithReveal(t *testing.T) {
	tableau := NewTableau(7, cards.NewDeck(1, 0).Shuffle())
	if tableau.Piles[6][5].Revealed {
		t.Error("Tableau.Piles[6][5] should start concealed")
	}
	cards, _ := tableau.Get(6, 6)
	if !tableau.Piles[6][5].Revealed {
		t.Error("Tableau.Piles[6][5] should now be revealed")
	}
	if len(tableau.Piles[6]) != 6 {
		t.Error("Tableau pile 6 should have 1 fewer card")
	}
	if len(cards) != 1 {
		t.Error("Gotten card slice should have 1 element")
	}
	tableau.Undo()
	if tableau.Piles[6][5].Revealed {
		t.Error("Tableau.Piles[6][5] should be concealed again")
	}
	if len(tableau.Piles[6]) != 7 {
		t.Error("Tableau pile 6 should have 7 cards after undoing")
	}
	if tableau.Piles[6][6] != cards[0] {
		t.Error("Tableau undo should have put the gotten card back")
	}
}

func TestTableau_GetValidSliceWithoutReveal(t *testing.T) {
	tableau := NewTableau(7, cards.NewDeck(1, 0).Shuffle())
	tableau.Piles[6][5] = &cards.Card{Pip: pip.King, Suit: suit.Hearts, Revealed: true}
	tableau.Piles[6][6] = &cards.Card{Pip: pip.Queen, Suit: suit.Hearts, Revealed: true}
	if !tableau.Piles[6][5].Revealed {
		t.Error("Tableau.Piles[6][5] should be revealed")
	}
	tableau.Get(6, 6)
	if !tableau.Piles[6][5].Revealed {
		t.Error("Tableau.Piles[6][5] should still be revealed")
	}
	tableau.Undo()
	if !tableau.Piles[6][5].Revealed {
		t.Error("Tableau.Piles[6][5] should still be revealed")
	}
}

func TestTableau_reveal(t *testing.T) {
	tableau := NewTableau(7, cards.NewDeck(1, 0).Shuffle())
	tableau.Piles[6][6].Conceal()
	if tableau.Piles[6][6].Revealed {
		t.Error("Tableau.Piles[6][6] should be concealed")
	}
	tableau.reveal(6)
	if !tableau.Piles[6][6].Revealed {
		t.Error("Tableau.Piles[6][6] should be revealed")
	}
}

func TestTableau_conceal(t *testing.T) {
	tableau := NewTableau(7, cards.NewDeck(1, 0).Shuffle())
	if !tableau.Piles[6][6].Revealed {
		t.Error("Tableau.Piles[6][6] should be revealed")
	}
	tableau.Piles[6][6].Conceal()
	tableau.conceal(6)
	if tableau.Piles[6][6].Revealed {
		t.Error("Tableau.Piles[6][6] should be concealed")
	}
}