package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"testing"
)

func TestNewKlondikeGame(t *testing.T) {
	k := NewKlondikeGame()
	if !k.Stock.IsShuffled {
		t.Error("Stock should be shuffled")
	}
	if len(k.Stock.Cards) != 24 {
		t.Error("There should be 24 cards remaining in the stock")
	}
	if len(k.Waste) != 0 {
		t.Error("There should be 0 cards in the waste")
	}
	if k.Score != 0 {
		t.Error("The score should be 0")
	}
	if len(k.Foundation.Piles) != 4 {
		t.Error("There should be 4 piles in the foundation")
	}
	if len(k.UndoStack) != 0 {
		t.Error("The undo stack should be empty")
	}
}

func TestKlondikeGame_Deal(t *testing.T) {
	k := NewKlondikeGame()
	// first deal
	k.Deal()
	if len(k.Stock.Cards) != 23 {
		t.Error("There should be 23 cards remaining in the stock")
	}
	if len(k.Waste) != 1 {
		t.Error("There should be 1 card in the waste")
	}
	if len(k.UndoStack) != 1 {
		t.Error("The undo stack should have 1 action")
	}
	//exhaust the stock
	for i := 0; i < 23; i++ {
		k.Deal()
	}
	if len(k.Stock.Cards) != 0 {
		t.Error("There should be 23 cards remaining in the stock")
	}
	if len(k.Waste) != 24 {
		t.Error("There should be 24 cards in the waste")
	}
	if len(k.UndoStack) != 24 {
		t.Error("The undo stack should have 24 actions")
	}

	//force the stock to be replenished
	k.Deal()
	if len(k.Stock.Cards) != 23 {
		t.Error("There should be 23 cards remaining in the stock")
	}
	if len(k.Waste) != 1 {
		t.Error("There should be 1 card in the waste")
	}
	if len(k.UndoStack) != 25 {
		t.Error("The undo stack should have 1 action")
	}

	//exhaust the deck and the waste and force an error
	k.Stock.Cards = []cards.Card{}
	k.Waste = []cards.Card{}
	if k.Deal() == nil {
		t.Error("Empty deck and waste should return an error")
	}
}

func TestKlondikeGame_undoDeal(t *testing.T) {
	k := NewKlondikeGame()
	k.Deal()
	k.Undo()
	if len(k.Stock.Cards) != 24 {
		t.Error("There should be 24 cards remaining in the stock")
	}
	if len(k.Waste) != 0 {
		t.Error("There should be 0 cards in the waste")
	}
	if len(k.UndoStack) != 0 {
		t.Error("The undo stack should have 0 actions")
	}
	for i := 0; i < 25; i++ {
		k.Deal()
	}
	k.Undo()
	if len(k.Stock.Cards) != 0 {
		t.Error("There should be 0 cards remaining in the stock")
	}
	if len(k.Waste) != 24 {
		t.Error("There should be 24 cards in the waste")
	}
	if len(k.UndoStack) != 24 {
		t.Error("The undo stack should have 0 actions")
	}
}

func TestKlondikeGame_adjustScore(t *testing.T) {
	k := NewKlondikeGame()
	k.adjustScore(100)
	if k.Score != 100 {
		t.Error("The score should be 100")
	}
	k.adjustScore(-113)
	if k.Score != -13 {
		t.Error("The score should be -13")
	}
}

func TestKlondikeGame_SelectFoundation(t *testing.T) {
	klondike := NewKlondikeGame()
	// set up the tableau top cards
	for pileNum, cardString := range []string{"10♦", "9♠", "J♦", "6♣", "3♦", "9♥", "2♦"} {
		klondike.Tableau.Piles[pileNum][len(klondike.Tableau.Piles[pileNum])-1], _ = cards.ParseCard(cardString)
	}
	// empty pile - return error
	err := klondike.SelectFoundation(suit.Hearts, 3)
	if err == nil {
		t.Error("should return error when foundation pile is empty")
	}
	// invalid pile - return error
	err = klondike.SelectFoundation("", 3)
	if err == nil {
		t.Error("should return error when suit is nonexistent")
	}
	// valid pile, invalid destination - return error
	card, _ := cards.ParseCard("5♥")
	klondike.Foundation.Piles[suit.Hearts] = append(klondike.Foundation.Piles[suit.Hearts], *card)
	err = klondike.SelectFoundation(suit.Hearts, 7)
	if err == nil {
		t.Error("Should return error when destination pile doesn't exist")
	}
	// valid pile, valid destination, no fit - return error
	err = klondike.SelectFoundation(suit.Hearts, 4)
	if err == nil {
		t.Error("Should return error when card doesn't fit the destination")
	}
	// valid pile, valid destination with fit
	err = klondike.SelectFoundation(suit.Hearts, 3)
	if err != nil {
		t.Error("Should not have returned an error when there's a valid fit.")
	}
	klondike.Undo()
	// valid pile, seek pile with fit
	err = klondike.SelectFoundation(suit.Hearts)
	if err != nil {
		t.Error("Should not have returned an error when there's a valid fit.")
	}

}
