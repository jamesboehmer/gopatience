package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/pip"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"testing"
)

func TestFoundation_Put(t *testing.T) {
	f := NewFoundation([]suit.Suit{suit.Hearts, suit.Diamonds, suit.Clubs, suit.Spades})
	err := f.Put(cards.Card{Pip: pip.Ace, Suit: suit.Hearts, Revealed: false})
	if err == nil {
		t.Error("Concealed cards should be rejected")
	}
	err = f.Put(cards.Card{Pip: pip.Ace, Suit: suit.Hearts, Revealed: true})
	if err != nil {
		t.Error("Ace should be able to be put on an empty foundation pile.")
	}
	err = f.Put(cards.Card{Pip: pip.Two, Suit: suit.Hearts, Revealed: true})
	if err != nil {
		t.Error("Two of Hearts should be able to be put on an Ace of Hearts in the foundation pile.")
	}
	err = f.Put(cards.Card{Pip: pip.King, Suit: suit.Hearts, Revealed: true})
	if err == nil {
		t.Error("Non-sequential puts should be rejected")
	}
	err = f.Put(cards.Card{Pip: pip.King, Suit: suit.Diamonds, Revealed: true})
	if err == nil {
		t.Error("King of Diamonds shouldn't be able to fit in an empty foundation")
	}
}

func TestFoundation_undoPut(t *testing.T) {
	f := NewFoundation([]suit.Suit{suit.Hearts, suit.Diamonds, suit.Clubs, suit.Spades})
	if len(f.Piles[suit.Hearts]) != 0 {
		t.Error("Hearts foundation pile should have 0 cards")
	}
	f.Put(cards.Card{Pip: pip.Ace, Suit: suit.Hearts, Revealed: true})
	if len(f.Piles[suit.Hearts]) != 1 {
		t.Error("Hearts foundation pile should have 1 card")
	}
	f.Undo()
	if len(f.Piles[suit.Hearts]) != 0 {
		t.Error("Hearts foundation pile should have 0 cards left")
	}
}

func TestFoundation_Get(t *testing.T) {
	f := NewFoundation([]suit.Suit{suit.Hearts, suit.Diamonds, suit.Clubs, suit.Spades})
	_, err := f.Get("")
	if err == nil {
		t.Error("Foundation.Get should have returned an error for a nonexistent suit")
	}
	_, err = f.Get(suit.Hearts)
	if err == nil {
		t.Error("Foundation.Get should have returned an error for an empty pile")
	}

	card := cards.Card{Pip: pip.Ace, Suit: suit.Hearts, Revealed: true}
	f.Piles[suit.Hearts] = append(f.Piles[suit.Hearts], card)
	if len(f.Piles[suit.Hearts]) != 1 {
		t.Error("Hearts foundation pile should have 1 card")
	}
	gottenCard, err := f.Get(suit.Hearts)
	if err != nil {
		t.Errorf("Foundation.Get shouldn't have returned an error :%s", err)
	}
	if len(f.Piles[suit.Hearts]) != 0 {
		t.Error("Hearts foundation pile should have 0 cards left")
	}
	if *gottenCard != card {
		t.Error("The card from the foundation should be the same as the card we put there.")
	}
	if len(f.UndoStack) != 1 {
		t.Error("The foundation's UndoStack should have 1 action in it.")
	}
}

func TestFoundation_undoGet(t *testing.T) {
	f := NewFoundation([]suit.Suit{suit.Hearts, suit.Diamonds, suit.Clubs, suit.Spades})
	f.Piles[suit.Hearts] = append(f.Piles[suit.Hearts], cards.Card{Pip: pip.Ace, Suit: suit.Hearts, Revealed: true})
	if len(f.Piles[suit.Hearts]) != 1 {
		t.Error("Hearts foundation pile should have 1 card")
	}
	f.Get(suit.Hearts)
	if len(f.Piles[suit.Hearts]) != 0 {
		t.Error("Hearts foundation pile should have 0 cards left")
	}
	f.Undo()
	if len(f.Piles[suit.Hearts]) != 1 {
		t.Error("Hearts foundation pile should have 1 card again")
	}
}

func TestFoundation_IsFull(t *testing.T) {
	f := NewFoundation([]suit.Suit{suit.Hearts, suit.Diamonds, suit.Clubs, suit.Spades})
	for _, suit := range suit.Suits {
		for _, pip := range []pip.Pip{pip.Ace, pip.Two, pip.Three, pip.Four, pip.Five, pip.Six, pip.Seven,
			pip.Eight, pip.Nine, pip.Ten, pip.Jack, pip.Queen, pip.King,
		} {
			if f.IsFull() {
				t.Error("Foundation should not be full yet.")
			}
			f.Put(cards.Card{Pip: pip, Suit: suit, Revealed: true})
		}
	}
	if !f.IsFull() {
		t.Error("Foundation should be full.")
	}
}
