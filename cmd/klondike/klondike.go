package main

import (
	"errors"
	"fmt"
	"github.com/jamesboehmer/gopatience/internal/cmd"
	"github.com/jamesboehmer/gopatience/pkg/cards"
	"github.com/jamesboehmer/gopatience/pkg/cards/suit"
	"github.com/jamesboehmer/gopatience/pkg/games/solitaire"
	"strconv"
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

	max := func(x, y int) int {
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

	var foundation []string
	for _, suit := range []suit.Suit{suit.Spades, suit.Diamonds, suit.Clubs, suit.Hearts} {
		pile := cmd.klondike.Foundation.Piles[suit]
		if len(pile) == 0 {
			foundation = append(foundation, fmt.Sprintf("[%s]", paintSuit(suit)))
		} else {
			paintedCard := paintCard(pile[len(pile)-1], 3, 0)
			foundation = append(foundation, fmt.Sprintf("%-3s", paintedCard))
		}
	}
	buffer.WriteString(fmt.Sprintf("Foundation: %s\n", strings.Join(foundation, " ")))

	buffer.WriteString("Tableau:\n")
	columnWidth := 3
	spacer := strings.Repeat(" ", 2)
	var tableauHeaders []string
	for pileNum, _ := range cmd.klondike.Tableau.Piles {
		tableauHeaders = append(tableauHeaders, fmt.Sprintf("%-*d", columnWidth, pileNum))
	}
	buffer.WriteString(
		fmt.Sprintf("%s\n", strings.Join(tableauHeaders, spacer)))
	var tableauDividers []string
	for range cmd.klondike.Tableau.Piles {
		tableauDividers = append(tableauDividers, strings.Repeat("-", columnWidth))
	}
	buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(tableauDividers, spacer)))

	//transpose the tableau piles
	rows := [13][]string{}
	for rowNum, _ := range rows {
		rows[rowNum] = make([]string, len(cmd.klondike.Tableau.Piles), len(cmd.klondike.Tableau.Piles))
		for pileNum, pile := range cmd.klondike.Tableau.Piles {
			if len(pile) > rowNum {
				rows[rowNum][pileNum] = paintCard(*pile[rowNum], 0, columnWidth)
			} else if rowNum == 0 {
				rows[rowNum][pileNum] = "[ ]"
			} else {
				rows[rowNum][pileNum] = "   "
			}
		}
	}
	for _, row := range rows {
		buffer.WriteString(fmt.Sprintf("%s\n", strings.Join(row, spacer)))
	}

	buffer.WriteString(strings.Repeat("\n", 6))
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
	cmd.klondike = solitaire.NewKlondikeGame()
	cmd.CommandPrompt = fmt.Sprintf("klondike[new]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doWaste(line string) (bool, error) {
	args := strings.Fields(line)
	if len(args) > 0 {
		pileNum, err := strconv.ParseInt(args[0], 0, 0)
		if err == nil {
			cmd.lastError = cmd.klondike.SelectWaste(int(pileNum))
			cmd.CommandPrompt = fmt.Sprintf("klondike[waste %d]> ", pileNum)
			return false, nil
		}
	}

	cmd.lastError = cmd.klondike.SelectWaste()
	cmd.CommandPrompt = fmt.Sprintf("klondike[waste]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doFoundation(line string) (bool, error) {
	args := strings.Fields(line)
	if len(args) > 0 {
		suit, found := map[string]suit.Suit{"c": suit.Clubs, "s": suit.Spades, "d": suit.Diamonds, "h": suit.Hearts}[args[0]]
		if !found {
			cmd.lastError = errors.New("no such suit")
			return false, nil
		} else {
			if len(args) < 2 {
				cmd.lastError = cmd.klondike.SelectFoundation(suit)
				cmd.CommandPrompt = fmt.Sprintf("klondike[foundation %s]> ", args[0])
				return false, nil
			} else {
				pileNum, err := strconv.ParseInt(args[1], 0, 0)
				if err == nil {
					cmd.lastError = cmd.klondike.SelectFoundation(suit, int(pileNum))
					cmd.CommandPrompt = fmt.Sprintf("klondike[foundation %s %d]> ", args[0], pileNum)
					return false, nil
				}
			}
		}
	}

	cmd.lastError = errors.New("usage: foundation c|s|d|h [pile]")
	return false, nil
}

func (cmd *KlondikeCmd) doSave(_ string) (bool, error) {
	cmd.lastError = errors.New("save not implemented yet")
	cmd.CommandPrompt = fmt.Sprintf("klondike[save]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doLoad(_ string) (bool, error) {
	cmd.lastError = errors.New("load not implemented yet")
	cmd.CommandPrompt = fmt.Sprintf("klondike[load]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doSolve(_ string) (bool, error) {
	cmd.lastError = cmd.klondike.Solve()
	cmd.CommandPrompt = fmt.Sprintf("klondike[solve]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doUndo(_ string) (bool, error) {
	cmd.lastError = cmd.klondike.Undo()
	cmd.CommandPrompt = fmt.Sprintf("klondike[undo]> ")
	return false, nil
}

func (cmd *KlondikeCmd) doTableau(line string) (bool, error) {
	args := strings.Fields(line)
	if len(args) == 0 {
		cmd.lastError = errors.New("usage: tableau <pileNum> [<cardNum> [toPile]]")
		return false, nil
	}
	var pileNum int64
	var otherArgs []int
	pileNum, err := strconv.ParseInt(args[0], 0, 0)
	if err != nil {
		cmd.lastError = err
		return false, nil
	}
	if len(args) > 1 {
		cardNum, err := strconv.ParseInt(args[1], 0, 0)
		if err != nil {
			cmd.lastError = err
			return false, nil
		}
		otherArgs = append(otherArgs, int(cardNum))
	}
	if len(args) > 2 {
		toPile, err := strconv.ParseInt(args[2], 0, 0)
		if err != nil {
			cmd.lastError = err
			return false, nil
		}
		otherArgs = append(otherArgs, int(toPile))
	}
	cmd.lastError = cmd.klondike.SelectTableau(int(pileNum), otherArgs...)
	cmd.CommandPrompt = fmt.Sprintf("klondike[tableau %s]> ", line)
	return false, nil
}

func (cmd *KlondikeCmd) Init() *KlondikeCmd {
	cmd.PostCmd = cmd.postCmd
	cmd.LastCmd = ""
	cmd.PreLoop = func() { cmd.printGame() }
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
