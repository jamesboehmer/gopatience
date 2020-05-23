package main

import (
	"fmt"
	"github.com/jamesboehmer/gopatience/internal/cmd"
	"github.com/jamesboehmer/gopatience/pkg/games/solitaire"
)

type KlondikeCmd struct {
	cmd.Cmd
	klondike *solitaire.KlondikeGame
}

func (cmd *KlondikeCmd) printGame() {
	fmt.Println("print game...")
}

func (cmd *KlondikeCmd) doQuit(_ string) (bool, error) {
	return true, nil
}

func (cmd *KlondikeCmd) doDeal(_ string) (bool, error) {
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
	fmt.Println("KlondikeCmd.postCmd")
	return stop
}

func main() {
	new(KlondikeCmd).Init().CommandLoop()
}
