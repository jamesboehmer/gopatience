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
