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

func TestKlondikeGame_SelectWasteWithDestination(t *testing.T) {
	klondike := NewKlondikeGame()
	for pileNum, cardString := range []string{"10♦", "9♠", "J♦", "6♣", "3♦", "9♥", "2♦"} {
		card, _ := cards.ParseCard(cardString)
		klondike.Tableau.Piles[pileNum][len(klondike.Tableau.Piles[pileNum])-1] = card
	}
	err := klondike.SelectWaste(0)
	if err == nil {
		t.Error("Should return error when the waste is empty")
	}
	// card in waste doesn't fit in destination - return error
	card, _ := cards.ParseCard("K♠")
	klondike.Waste = append(klondike.Waste, *card)
	err = klondike.SelectWaste(0)
	if err == nil {
		t.Error("Should return error when the waste card doesn't fit in the specified tableau destination")
	}

	// card in waste fits in destination
	klondike.Tableau.Piles[0] = make([]*cards.Card, 0, 13)
	err = klondike.SelectWaste(0)
	if err != nil {
		t.Error("Should not have returned an error when there was a tableau fit")
	}
}

func TestKlondikeGame_SelectWasteWithoutDestination(t *testing.T) {
	klondike := NewKlondikeGame()
	// If the waste is empty, return an error
	err := klondike.SelectWaste()
	if err == nil {
		t.Error("Should return error when the waste is empty")
	}

	// If there's a fit in the foundation AND the tableau, the waste card should go to the foundation
	card, _ := cards.ParseCard("A♠")
	klondike.Waste = append(klondike.Waste, *card)
	err = klondike.SelectWaste()
	if err != nil {
		t.Error("The waste card should have fit the foundation with no error")
	}
	if len(klondike.Waste) != 0 {
		t.Error("The waste should be empty")
	}
	if len(klondike.Foundation.Piles[suit.Spades]) != 1 {
		t.Error("The Spades pile should have 1 card")
	}
	if klondike.Score != PointsWasteFoundation {
		t.Error("The score should be %n", PointsWasteFoundation)
	}
	if len(klondike.Foundation.UndoStack) != 1 {
		t.Error("There should be 1 undo action in the foundation")
	}
	if len(klondike.Tableau.UndoStack) != 0 {
		t.Error("There should be 0 undo actions in the tableau")
	}
	// Check that the undo function also impacts the score and the foundation
	klondike.Undo()
	if len(klondike.Waste) != 1 {
		t.Error("The waste should have 1 card again")
	}
	if len(klondike.Foundation.Piles[suit.Spades]) != 0 {
		t.Error("The Spades pile should have 0 cards")
	}
	if klondike.Score != 0 {
		t.Error("The score should be reset")
	}
	if len(klondike.Foundation.UndoStack) != 0 {
		t.Error("There should be 0 undo actions in the foundation")
	}
	if len(klondike.Tableau.UndoStack) != 0 {
		t.Error("There should be 0 undo actions in the tableau")
	}

	// If there's no foundation fit, but there's a tableau fit, it should find the right tableau pile
	card, _ = cards.ParseCard("9♣")
	klondike.Waste = []cards.Card{*card}
	for pileNum, cardString := range []string{"10♦", "9♠", "J♦", "6♣", "3♦", "9♥", "2♦"} {
		card, _ = cards.ParseCard(cardString)
		klondike.Tableau.Piles[pileNum][len(klondike.Tableau.Piles[pileNum])-1] = card
	}
	err = klondike.SelectWaste()
	if err != nil {
		t.Error("The waste card should have found a tableau pile")
	}
	if len(klondike.Waste) != 0 {
		t.Error("The waste should be empty")
	}
	if len(klondike.Tableau.Piles[0]) != 2 {
		t.Error("Tableau pile 0 should now have 2 cards")
	}
	if klondike.Score != PointsWasteTableau {
		t.Error("The score should be %n", PointsWasteTableau)
	}
	if len(klondike.Foundation.UndoStack) != 0 {
		t.Error("There should be 0 undo actions in the foundation")
	}
	if len(klondike.Tableau.UndoStack) != 1 {
		t.Error("There should be 1 undo action in the tableau")
	}

	// check that the undo function affects the score and the tableau
	klondike.Undo()
	if len(klondike.Waste) != 1 {
		t.Error("The waste should have 1 card again")
	}
	if len(klondike.Tableau.Piles[0]) != 1 {
		t.Error("Tableau pile 0 should now have 1 card")
	}
	if klondike.Score != 0 {
		t.Error("The score should be reset")
	}
	if len(klondike.Foundation.UndoStack) != 0 {
		t.Error("There should be 0 undo actions in the foundation")
	}
	if len(klondike.Tableau.UndoStack) != 0 {
		t.Error("There should be 0 undo actions in the tableau")
	}


	// If there's no foundation fit and no tableau fit, return an error
	card, _ = cards.ParseCard("K♥")
	klondike.Waste = []cards.Card{*card}
	err = klondike.SelectWaste()
	if err == nil {
		t.Error("Should have returned an error when there's no foundation or tableau fit")
	}
}

func TestKlondikeGame_SelectTableauInvalidPiles(t *testing.T) {
	k := NewKlondikeGame()

	// pile_num == destination_pile_num - should return error
	if k.SelectTableau(0, 0, 0) == nil {
		t.Error("Should have returned an error if pileNum==destinationPile")
	}
	//valid pile_num, valid single card_num, invalid destination
	if k.SelectTableau(0, -1, 7) == nil {
		t.Error("Should have returned an error if destinationPile is invalid")
	}

	// valid pileNum, invalid cardNum
	if k.SelectTableau(0, 100) == nil {
		t.Error("Should have returned an error if cardNum is invalid")
	}
	if k.SelectTableau(0, -2) == nil {
		t.Error("Should have returned an error if cardNum is negative and can't be transformed")
	}
	// invalid pileNum
	if k.SelectTableau(-1) == nil {
		t.Error("Should have returned an error if pileNum is invalid")
	}
	if k.SelectTableau(7) == nil {
		t.Error("Should have returned an error if pileNum is invalid")
	}

	// valid pileNum, empty pile
	k.Tableau.Piles[0] = []*cards.Card{}
	if k.SelectTableau(0) == nil {
		t.Error("Should have returned an error if pile is empty")
	}
}

func TestKlondikeGame_SelectTableauNoDestination(t *testing.T) {
	k := NewKlondikeGame()
	// valid pile_num, valid single card_num, no destination, no foundation fit, no tableau fit - return error

	for pileNum, cardString := range []string{"10♦", "9♠", "J♦", "6♣", "3♦", "9♥", "2♦"}{
		card, _ := cards.ParseCard(cardString)
		k.Tableau.Piles[pileNum][len(k.Tableau.Piles[pileNum])-1] = card
	}
	if k.SelectTableau(3, -1) == nil {
		t.Error("Should have errored")
	}
	// valid pile_num, valid single card_num, no destination, foundation fit
	card, _ := cards.ParseCard("5♣")
	k.Foundation.Piles[suit.Clubs] =append(k.Foundation.Piles[suit.Clubs], *card)
	if k.SelectTableau(3, -1) != nil {
		t.Error("Should have found a foundation fit, not an error")
	}
	// valid pile_num, valid single card_num, no destination, no foundation fit, tableau fit
	k = NewKlondikeGame()
	for pileNum, cardString := range []string{"10♦", "9♠", "J♦", "6♣", "3♦", "9♥", "2♦"}{
		card, _ := cards.ParseCard(cardString)
		k.Tableau.Piles[pileNum][len(k.Tableau.Piles[pileNum])-1] = card
	}
	if k.SelectTableau(1, 1) != nil {
		t.Error("Should have found a tableau fit, not an error")
	}

	// valid pile_num, valid multi card_num, no destination, tableau fit
	k = NewKlondikeGame()
	for pileNum, cardString := range []string{"10♦", "8♥", "J♦", "6♣", "3♦", "9♥", "2♦"}{
		card, _ := cards.ParseCard(cardString)
		k.Tableau.Piles[pileNum][len(k.Tableau.Piles[pileNum])-1] = card
	}
	card, _ = cards.ParseCard("9♠")
	k.Tableau.Piles[1][0] = card
	if k.SelectTableau(1, 0) != nil {
		t.Error("Should have found a tableau fit")
	}
}

func TestKlondikeGame_SelectTableauValidDestination(t *testing.T) {
	k := NewKlondikeGame()
	for pileNum, cardString := range []string{"10♦", "9♠", "J♦", "6♣", "3♦", "9♥", "2♦"}{
		card, _ := cards.ParseCard(cardString)
		k.Tableau.Piles[pileNum][len(k.Tableau.Piles[pileNum])-1] = card
	}
	// valid pile_num, valid single card_num, valid destination, no tableau fit - error
	if k.SelectTableau(5, 5, 0) == nil {
		t.Error("Should have raised an error with no tableau fit")
	}
	// valid pile_num, valid single card_num, valid destination, tableau fit
	if k.SelectTableau(1, 1, 0) != nil {
		t.Error("Should have found a tableau fit")
	}
}

func TestKlondikeGame_IsSolvable(t *testing.T) {
	k := NewKlondikeGame()
	// remaining stock cards - false
	if k.IsSolvable() {
		t.Error("New game should not be solvable yet")
	}

	// remaining waste cards - false
	k.Deal()
	k.Stock.Cards = []cards.Card{}
	if k.IsSolvable() {
		t.Error("New game should not be solvable with waste cards remaining")
	}

	// concealed cards in the tableau - false
	k.Waste = []cards.Card{}
	if k.IsSolvable() {
		t.Error("New game should not be solvable with concealed cards in the tableau")
	}

	// all tableau cards revealed = true
	for _, pile := range k.Tableau.Piles {
		for _, card := range pile {
			card.Reveal()
		}
	}
	if !k.IsSolvable() {
		t.Error("The game should be solvable with empty deck snd waste, and all tableau cards revealed")
	}


}