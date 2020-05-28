package main

import (
	"fmt"
	"github.com/jamesboehmer/gopatience/internal/cmd"
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"github.com/jamesboehmer/gopatience/pkg/games/solitaire"
	"strings"
)

type KlondikeCmd struct {
	cmd.Cmd
	klondike  *solitaire.KlondikeGame
	lastError error
}

const (
	StyleBright    = "\033[1m"
	StyleReset     = "\033[0m"
	ColorForeBlack = "\033[30m"
	ColorForeRed   = "\033[31m"
	ColorBackRed   = "\033[41m"
	ClearScreen    = "\033[H\033[J"
)

func (cmd *KlondikeCmd) printGame() {
	paintSuit := func(s suit.Suit) string {
		var color = ""
		if s.Color() == suit.Red {
			color = ColorForeRed
		}
		return fmt.Sprintf("%s%s%s%s", StyleBright, color, s, StyleReset)
	}

	max := func (x, y int) int {
		if x > y {
		return x
	}
		return y
	}
	paintCard := func(card cards.Card, leftPad int, rightPad int) string {
		cardString := card.String()
		if !card.Revealed {
			cardString = "#"
		}
		length := len([]rune(cardString))
		left := strings.Repeat(" ", max(leftPad-length, 0))
		right := strings.Repeat(" ", max(rightPad-length-len(left), 0))

		if !card.Revealed || card.Pip == "" {
			return fmt.Sprintf("%s%s%s%s%s", left, StyleBright, cardString, StyleReset, right)
		}
		return fmt.Sprintf("%s%s%s%s%s%s", left, StyleBright, card.Pip, paintSuit(card.Suit), StyleReset, right)

	}
	buffer := strings.Builder{}
	buffer.WriteString(ClearScreen)
	var status string
	if cmd.klondike.IsSolved() {
		status = "Solved!"
	} else if cmd.lastError != nil {
		status = cmd.lastError.Error()
		cmd.lastError = nil
	} else {
		status = ""
	}
	buffer.WriteString(fmt.Sprintf("%s%s%s%s%s\n", StyleBright, ColorBackRed, ColorForeBlack, status, StyleReset))

	buffer.WriteString(fmt.Sprintf("Score: %d\n", cmd.klondike.Score))

	buffer.WriteString(fmt.Sprintf("Stock: %d\n", cmd.klondike.Stock.Remaining()))

	var paintedCards []string
	for _, c := range cmd.klondike.Waste {
		paintedCards = append(paintedCards, paintCard(c, 0, 0))
	}
	waste := fmt.Sprintf("[%s]", strings.Join(paintedCards, ", "))
	buffer.WriteString(fmt.Sprintf("Waste: %s\n", waste))

	//TODO: foundation
	//TODO: tableau

	buffer.WriteString(strings.Repeat("\n", 19))
	fmt.Println(buffer.String())

}

func (cmd *KlondikeCmd) doQuit(_ string) (bool, error) {
	return true, nil
}

func (cmd *KlondikeCmd) doDeal(_ string) (bool, error) {
	cmd.klondike.Deal()
	cmd.CommandPrompt = fmt.Sprintf("klondike[deal]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doNew(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[new]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doWaste(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[waste]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doFoundation(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[foundation]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doSave(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[save]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doLoad(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[load]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doSolve(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[solve]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doUndo(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[undo]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doTableau(_ string) (bool, error) {
	cmd.CommandPrompt = fmt.Sprintf("klondike[tableau]> ")
	return false, nil
}

func (cmd *KlondikeCmd) Init() *KlondikeCmd {
	cmd.PostCmd = cmd.postCmd
	cmd.LastCmd = ""
	cmd.PreLoop = func () {cmd.printGame()}
	cmd.CommandPrompt = "klondike> "
	cmd.FunctionMap = map[string]func(string) (bool, error){
		"d":          cmd.doDeal,
		"deal":       cmd.doDeal,
		"w":          cmd.doWaste,
		"waste":      cmd.doWaste,
		"n":          cmd.doNew,
		"new":        cmd.doNew,
		"t":          cmd.doTableau,
		"tableau":    cmd.doTableau,
		"f":          cmd.doFoundation,
		"foundation": cmd.doFoundation,
		"s":          cmd.doSave,
		"save":       cmd.doSave,
		"l":          cmd.doLoad,
		"load":       cmd.doLoad,
		"solve":      cmd.doSolve,
		"u":          cmd.doUndo,
		"undo":       cmd.doUndo,
		"q":          cmd.doQuit,
		"quit":       cmd.doQuit,
	}
	cmd.klondike = solitaire.NewKlondikeGame()
	return cmd
}

func (cmd *KlondikeCmd) postCmd(stop bool, line string) bool {
	cmd.printGame()
	return stop
}

func main() {
	new(KlondikeCmd).Init().CommandLoop()
}
