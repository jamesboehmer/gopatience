package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
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
