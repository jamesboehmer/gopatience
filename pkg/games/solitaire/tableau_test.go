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

func TestTableau_PutConcealedOnEmptyTableau(t *testing.T) {
	tableau := NewTableau(7, nil)
	for card := range cards.NewDeck(1, 0).Shuffle().DealAll() {
		if tableau.Put([]*cards.Card{card}, 0) == nil {
			t.Error("Putting a concealed card on the tableau should return an error")
		}
	}
}

func TestTableau_PutConcealedOnBuiltTableau(t *testing.T) {
	tableau := NewTableau(7, nil)
	for pileNum, suit := range []suit.Suit{suit.Hearts, suit.Diamonds, suit.Spades, suit.Clubs} {
		tableau.Put([]*cards.Card{{Pip: pip.King, Suit: suit, Revealed: true}}, pileNum)
	}
	for pileNum := 0; pileNum < len(suit.Suits); pileNum++ {
		if tableau.Put([]*cards.Card{new(cards.Card)}, pileNum) == nil {
			t.Error("Putting a concealed card on the tableau should return an error")
		}
	}
}

func TestTableau_PutNonKingsOnEmptyTableau(t *testing.T) {
	tableau := NewTableau(7, nil)
	for card := range cards.NewDeck(1, 0).Shuffle().DealAll() {
		if card.Pip != pip.King {
			if tableau.Put([]*cards.Card{card.Reveal()}, 0) == nil {
				t.Error("Putting a non-king on an empty tableau pile should return an error")
			}
		}
	}
}

func TestTableau_PutBuildOrder(t *testing.T) {
	tableau := NewTableau(7, nil)
	for pileNum, suit := range []suit.Suit{suit.Spades, suit.Diamonds, suit.Clubs, suit.Hearts} { // B/R/B/R
		tableau.Put([]*cards.Card{{Pip: pip.King, Suit: suit, Revealed: true}}, pileNum)
	}
	// Test invalid color with valid value
	for pileNum, suit := range []suit.Suit{suit.Spades, suit.Diamonds, suit.Clubs, suit.Hearts} { // B/R/B/R
		if tableau.Put([]*cards.Card{{Pip: pip.Queen, Suit: suit, Revealed: true}}, pileNum) == nil {
			t.Error("Tableau should return an error when building like colors")
		}
	}
	// Test valid color with valid value
	for pileNum, suit := range []suit.Suit{suit.Diamonds, suit.Clubs, suit.Hearts, suit.Spades} { // R/B/R/B
		if tableau.Put([]*cards.Card{{Pip: pip.Queen, Suit: suit, Revealed: true}}, pileNum) != nil {
			t.Error("Tableau should not return an error when building correct alternate colors")
		}
	}
	// Test valid color with invalid value
	for pileNum, suit := range []suit.Suit{suit.Clubs, suit.Hearts, suit.Spades, suit.Diamonds} { // B/R/B/R
		if tableau.Put([]*cards.Card{{Pip: pip.Ten, Suit: suit, Revealed: true}}, pileNum) == nil {
			t.Error("Tableau should return an error when building incorrect descending pips")
		}
	}
}
func TestTableau_PutInvalidPileNum(t *testing.T) {
	tableau := NewTableau(7, nil)
	if tableau.Put([]*cards.Card{{Pip: pip.King, Suit: suit.Hearts, Revealed: true}}, -1) == nil {
		t.Error("Negative pileNum should return an error")
	}
	if tableau.Put([]*cards.Card{{Pip: pip.King, Suit: suit.Hearts, Revealed: true}}, 7) == nil {
		t.Error("Invalid pileNum should return an error")
	}
}
func TestTableau_undoPut(t *testing.T) {
	tableau := NewTableau(7, nil)
	for pileNum, suit := range []suit.Suit{suit.Spades, suit.Diamonds, suit.Clubs, suit.Hearts} { // B/R/B/R
		tableau.Put([]*cards.Card{{Pip: pip.King, Suit: suit, Revealed: true}}, pileNum)
	}
	if len(tableau.UndoStack) != 4 {
		t.Error("There should be exactly 4 undo actions")
	}
	for i := 3; i >= 0; i-- {
		if len(tableau.Piles[i]) != 1 {
			t.Error("Pile %n should have 1 card", i)
		}
		tableau.Undo()
		if len(tableau.Piles[i]) != 0 {
			t.Error("Pile %n should have 0 cards", i)
		}
	}
}
