package solitaire

import (
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/pip"
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

var PipValue = map[pip.Pip]int{
	pip.Ace:   1,
	pip.Two:   2,
	pip.Three: 3,
	pip.Four:  4,
	pip.Five:  5,
	pip.Six:   6,
	pip.Seven: 7,
	pip.Eight: 8,
	pip.Nine:  9,
	pip.Ten:   10,
	pip.Jack:  11,
	pip.Queen: 12,
	pip.King:  13,
}
//
//func PipValue(p pip.Pip) int {
//	return map[pip.Pip]int{
//		pip.Ace:   1,
//		pip.Two:   2,
//		pip.Three: 3,
//		pip.Four:  4,
//		pip.Five:  5,
//		pip.Six:   6,
//		pip.Seven: 7,
//		pip.Eight: 8,
//		pip.Nine:  9,
//		pip.Ten:   10,
//		pip.Jack:  11,
//		pip.Queen: 12,
//		pip.King:  13,
//	}[p]
//}
