package cards

import (
	"github.com/jamesboehmer/gopatience/pkg/cards/pip"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"strings"
	"testing"
)

func TestCardFaces(t *testing.T) {
	deck := NewDeck(1, 2)
	facePips := map[pip.Pip]interface{}{pip.Jack: nil, pip.Queen: nil, pip.King: nil}
	for _, card := range deck.Cards {
		_, found := facePips[card.Pip]
		if found {
			if !card.IsFace() {
				t.Errorf("%s should be considered a face card", card.String())
			}
		} else {
			if card.IsFace() {
				t.Errorf("%s should not be considered a face card", card.String())
			}
		}
	}
}

func TestNewDeck(t *testing.T) {
	deck := NewDeck(1, 0)
	if len(deck.Cards) != 52 {
		t.Error("Standard deck should have 52 cards.")
	}
	deck = NewDeck(2, 0)
	if len(deck.Cards) != 104 {
		t.Error("Double deck should have 104 cards.")
	}
	deck = NewDeck(1, 2)
	if len(deck.Cards) != 54 {
		t.Error("Double deck should have 54 cards.")
	}
}

func TestDeck_Shuffle(t *testing.T) {
	deck := NewDeck(1, 0)
	if deck.IsShuffled {
		t.Error("New deck should not be shuffled")
	}
	originalCards := strings.Builder{}
	for _, card := range deck.Cards {
		originalCards.WriteString(card.String())
	}

	deck.Shuffle()
	if !deck.IsShuffled {
		t.Error("Shuffled deck isShuffled should be true")
	}
	shuffledCards := strings.Builder{}
	for _, card := range deck.Cards {
		shuffledCards.WriteString(card.String())
	}
	if originalCards.String() == shuffledCards.String() {
		t.Error("Shuffled cards should be in a different order.")
	}

	deck = NewDeck(1, 0).Shuffle()
	if !deck.IsShuffled {
		t.Error("Continuation-shuffled deck isShuffled should be true")
	}
}

func TestDeck_Deal_Remaining(t *testing.T) {
	deck := NewDeck(1, 0)
	if deck.Remaining() != 52 {
		t.Error("New Deck should have 52 cards.")
	}
	card1 := deck.Cards[0]
	card2 := deck.Cards[1]
	card3 := deck.Cards[2]
	dealtCard1 := *deck.Deal()
	dealtCard2 := *deck.Deal()
	dealtCard3 := *deck.Deal()
	if card1 != dealtCard1 || card2 != dealtCard2 || card3 != dealtCard3 {
		t.Error("Dealt cards should be unique.")
	}
	if deck.Remaining() != 49 {
		t.Error("Deck should have 49 cards remaining.")
	}
	for i := 0; i < 49; i++ {
		deck.Deal()
	}
	if deck.Remaining() != 0 {
		t.Error("Deck should have 0 cards remaining.")
	}
	card := deck.Deal()
	if card != nil {
		t.Error("Deck should deal nil when empty")
	}

}

func TestCard_Reveal_Conceal(t *testing.T) {
	card := new(Card).Reveal()
	if !card.Revealed {
		t.Error("Card should be revealed")
	}
	card.Conceal()
	if card.Revealed {
		t.Error("Card should be concealed")
	}
}

func TestCard_String(t *testing.T) {
	deck := NewDeck(1, 1)
	for deck.Remaining() > 0 {
		card := deck.Deal()
		cardString := []rune(card.String())
		if string(cardString[0]) != "|" {
			t.Error("Concealed card string should start with |")
		}
		cardString = []rune(card.Reveal().String())
		if string(cardString[0]) == "|" {
			t.Error("Revealed card string should not start with |")
		}
		if card.Suit == "" {
			if string(cardString) != "*" {
				t.Error("Joker string should be an asterisk")
			}
		} else {
			suitString := cardString[len(cardString)-1:]
			if string(suitString) != string(card.Suit) {
				t.Error("Card string suit should match the card's Suit")
			}
			pipString := cardString[:len(cardString)-1]
			if string(pipString) != string(card.Pip) {
				t.Error("Card string pip should match the card's Pip")
			}
		}

	}
}

func TestParseCard(t *testing.T) {
	_, err := ParseCard("|10x")
	if err == nil {
		t.Error("Bad suit string should have returned an error")
	}
	_, err = ParseCard("Z♠")
	if err == nil {
		t.Error("Bad pip string should have returned an error")
	}
	card, err := ParseCard("*")
	if err != nil {
		t.Error("Joker card should have been parsed")
	}
	if card.Pip != "" || card.Suit != "" {
		t.Error("Joker suit and pip should be empty")
	}
	for _, suit := range suit.Suits {
		for _, pip := range pip.Pips {
			for _, revealed := range []bool{true, false} {
				card := Card{pip, suit, revealed}
				cardString := card.String()
				parsedCard, err := ParseCard(cardString)
				if err != nil {
					t.Error("Card parse should not have failed.")
				} else {
					if parsedCard.Pip != pip {
						t.Error("Parsed card pip doesn't match.")
					}
					if parsedCard.Suit != suit {
						t.Error("Parsed card suit doesn't match.")
					}
					if parsedCard.Revealed != revealed {
						t.Error("Parsed card state doesn't match.")
					}
				}

			}
		}
	}
}

func TestSuitColor(t *testing.T) {
	if suit.Spades.Color() != suit.Black {
		t.Errorf("%s should be black", suit.Spades)
	}
	if suit.Clubs.Color() != suit.Black {
		t.Errorf("%s should be black", suit.Clubs)
	}
	if suit.Diamonds.Color() != suit.Red {
		t.Errorf("%s should be red", suit.Diamonds)
	}
	if suit.Hearts.Color() != suit.Red {
		t.Errorf("%s should be red", suit.Hearts)
	}
}
