package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"github.com/jamesboehmer/gopatience/pkg/games/util"
)

type KlondikeGame struct {
	util.Undoable
	Score      int
	Errors     []error
	Stock      cards.Deck
	Waste      []cards.Card
	Foundation Foundation
	Tableau    Tableau
}

func (k *KlondikeGame) Deal() error {
	return nil
}

func (k *KlondikeGame) undoDeal() error {
	return nil
}

func (k *KlondikeGame) adjustScore(points int) {
	k.Score += points
}

func (k *KlondikeGame) SelectFoundation(suit suit.Suit, pileNum ...int) error {
	return nil
}

func (k *KlondikeGame) undoSelectFoundation() error {
	return nil
}

func (k *KlondikeGame) SelectWaste(pileNum ...int) error {
	return nil
}

func (k *KlondikeGame) undoSelectWaste() error {
	return nil
}

func (k *KlondikeGame) seekTableauToFoundation() error {
	return nil
}

func (k *KlondikeGame) undoSeekTableauToFoundation() error {
	return nil
}

func (k *KlondikeGame) SelectTableau(pileCardDestination ...int) error {
	//TODO: pileNum int, cardNum int, destination int from pileCardDestination
	return nil
}

func (k *KlondikeGame) undoSelectTableau() error {
	return nil
}

func (k *KlondikeGame) IsSolvable() bool {
	return false
}

func (k *KlondikeGame) IsSolved() bool {
	return false
}

func (k *KlondikeGame) Solve() error {
	return nil
}

func NewKlondikeGame() *KlondikeGame {
	game := new(KlondikeGame)
	game.Stock = *cards.NewDeck(1, 0)
	game.Foundation = *NewFoundation([]suit.Suit{suit.Hearts, suit.Diamonds, suit.Clubs, suit.Spades})
	game.Tableau = *NewTableau(7, &game.Stock)
	return game
}
