package solitaire

import (
	"errors"
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

const (
	PointsWasteFoundation   int = 10
	PointsWasteTableau      int = 5
	PointsTableauFoundation int = 15
)

func (k *KlondikeGame) Deal() error {
	replenished := false
	if k.Stock.Remaining() == 0 {
		if len(k.Waste) > 0 {
			for _, card := range k.Waste {
				k.Stock.Cards = append(k.Stock.Cards, *card.Conceal())
			}
			k.Waste = []cards.Card{}
			replenished = true
		} else {
			return errors.New("no cards remaining")
		}
	}
	k.Waste = append(k.Waste, *k.Stock.Deal().Reveal())
	k.UndoStack = append(k.UndoStack, util.UndoAction{
		Function: k.undoDeal,
		Args:     []interface{}{replenished},
	})
	return nil
}

func (k *KlondikeGame) undoDeal(args ...interface{}) error {
	replenished := args[0].(bool)
	card := k.Waste[len(k.Waste)-1]
	k.Waste = k.Waste[:len(k.Waste)-1]
	k.Stock.Cards = append([]cards.Card{*card.Reveal()}, k.Stock.Cards...)
	if replenished {
		for card := range k.Stock.DealAll() {
			k.Waste = append(k.Waste, *card.Reveal())
		}
	}
	return nil
}

func (k *KlondikeGame) adjustScore(points int) {
	k.Score += points
}

func (k *KlondikeGame) SelectFoundation(suit suit.Suit, tableauDestinations ...int) error {
	card, err := k.Foundation.Get(suit)
	if err != nil {
		return err
	}
	if tableauDestinations == nil || len(tableauDestinations) == 0 {
		for i := 0; i < len(k.Tableau.Piles); i++ {
			tableauDestinations = append(tableauDestinations, i)
		}
	}
	for _, pileNum := range tableauDestinations {
		err := k.Tableau.Put([]*cards.Card{card}, pileNum)
		if err == nil {
			k.adjustScore(-PointsTableauFoundation)
			k.UndoStack = append(k.UndoStack, util.UndoAction{
				Function: k.undoSelectFoundation,
				Args:     nil,
			})
			return nil
		}
	}
	k.Foundation.Undo()
	return errors.New("no tableau fit")
}

func (k *KlondikeGame) undoSelectFoundation(...interface{}) error {
	k.adjustScore(PointsTableauFoundation)
	k.Tableau.Undo()
	k.Foundation.Undo()
	return nil
}

func (k *KlondikeGame) SelectWaste(tableauDestinations ...int) error {
	if len(k.Waste) == 0 {
		return errors.New("no cards left in the waste pile")
	}
	topCard := k.Waste[len(k.Waste)-1]
	k.Waste = k.Waste[:len(k.Waste)-1]

	// try moving from the waste to the foundation if there was no tableau pile specified
	if tableauDestinations == nil {
		err := k.Foundation.Put(topCard)
		if err == nil {
			k.adjustScore(PointsWasteFoundation)
			k.UndoStack = append(k.UndoStack, util.UndoAction{
				Function: k.undoSelectWaste,
				Args:     []interface{}{true, topCard},
			})
			return nil
		}
	}
	// if there was no fit, then try the tableau
	if tableauDestinations == nil || len(tableauDestinations) == 0 {
		for i := 0; i < len(k.Tableau.Piles); i++ {
			tableauDestinations = append(tableauDestinations, i)
		}
	}

	for _, pileNum := range tableauDestinations {
		err := k.Tableau.Put([]*cards.Card{&topCard}, pileNum)
		if err == nil {
			k.adjustScore(PointsWasteTableau)
			k.UndoStack = append(k.UndoStack, util.UndoAction{
				Function: k.undoSelectWaste,
				Args:     []interface{}{false, topCard},
			})
			return nil
		}
	}
	// if there was no fit, put the card back on the waste pile
	k.Waste = append(k.Waste, topCard)
	return errors.New("no tableau fit")
}

func (k *KlondikeGame) undoSelectWaste(args ...interface{}) error {
	undoFoundation, card := args[0].(bool), args[1].(cards.Card)
	if undoFoundation {
		k.adjustScore(-PointsWasteFoundation)
		k.Foundation.Undo()
	} else {
		k.adjustScore(-PointsWasteTableau)
		k.Tableau.Undo()
	}
	k.Waste = append(k.Waste, card)
	return nil
}

func (k *KlondikeGame) seekTableauToFoundation() error {
	return nil
}

func (k *KlondikeGame) undoSeekTableauToFoundation() error {
	return nil
}

func (k *KlondikeGame) SelectTableau(pileNum int, cardDestination ...int) error {
	if pileNum < 0 || pileNum > len(k.Tableau.Piles)-1 {
		return errors.New("invalid pileNum")
	}
	if len(k.Tableau.Piles[pileNum]) == 0 {
		return errors.New("empty pile")
	}
	cardNum := -1
	if len(cardDestination) > 0 {
		cardNum = cardDestination[0]
	}
	if len(cardDestination) > 1 {
		if cardDestination[1] == pileNum || cardDestination[1] < 0 { //|| cardDestination[1] > len(k.Tableau.Piles)-1{
			return errors.New("invalid destination")
		}
	}
	if cardNum < 0 {
		// it's valid to ask for a negative index, just convert it to the positive offset
		cardNum = len(k.Tableau.Piles[pileNum])+cardNum
		if cardNum < 0 {
			return errors.New("invalid cardNum")
		}
	}
	cards, err := k.Tableau.Get(pileNum, cardNum) //needs to be undone if we can't find a fit
	if err != nil {
		return err
	}
	// If there's only 1 card selected from the tableau, and no destination specified, try to fit it in the foundation
	if len(cards) == 1 && len(cardDestination) < 2 {
		err = k.Foundation.Put(*cards[0])
		if err == nil {
			k.adjustScore(PointsTableauFoundation)
			k.UndoStack = append(k.UndoStack, util.UndoAction{
				Function: k.undoSelectTableau,
				Args:     []interface{}{true},
			})
			return nil
		}
		// don't quit here just because we didn't find a foundation fit.
	}

	// if there was no fit, then try the tableau
	var tableauDestinations []int
	if len(cardDestination) > 1 {
		tableauDestinations = append(tableauDestinations, cardDestination[1])
	} else {
		for i := 0; i < len(k.Tableau.Piles); i++ {
			tableauDestinations = append(tableauDestinations, i)
		}
	}

	for _, pileNum := range tableauDestinations {
		err := k.Tableau.Put(cards, pileNum)
		if err == nil {
			k.adjustScore(PointsWasteTableau)
			k.UndoStack = append(k.UndoStack, util.UndoAction{
				Function: k.undoSelectTableau,
				Args:     []interface{}{true},
			})
			return nil
		}
	}
	// We couldn't find a card in the tableau that fit in the foundation
	// OR The chosen tableau card didn't fit in the foundation
	// OR The chosen tableau card didn't fit anywhere in the tableau
	// OR the chosen tableau card didn't fit in the chosen tableau pile
	k.Tableau.Undo()
	return errors.New("no fit for chosen card(s)")
}

func (k *KlondikeGame) undoSelectTableau(args ...interface{}) error {
	undoFoundation := args[0].(bool)
	if undoFoundation {
		k.adjustScore(-PointsTableauFoundation)
		k.Foundation.Undo() //undo put
	} else {
		k.Tableau.Undo() //undo put
		k.Tableau.Undo() //undo get
	}
	return nil
}

func (k *KlondikeGame) IsSolvable() bool {
	if k.Stock.Remaining() + len(k.Waste) > 0 {
		return false
	}
	for _, pile := range k.Tableau.Piles {
		if len(pile) > 0 && !pile[0].Revealed {
			return false
		}
	}
	return true
}

func (k *KlondikeGame) IsSolved() bool {
	return false
}

func (k *KlondikeGame) Solve() error {
	return nil
}

func NewKlondikeGame() *KlondikeGame {
	game := new(KlondikeGame)
	game.Stock = *cards.NewDeck(1, 0).Shuffle()
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
