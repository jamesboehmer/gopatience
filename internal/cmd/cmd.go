package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Cmd struct {
	LastCmd       string
	CommandPrompt string
	FunctionMap   map[string]func(string) (bool, error)
	Error         error
	PostCmd       func(bool, string) bool
	PreCmd        func(string) string
	PreLoop       func()
	PostLoop      func()
}

func (cmd *Cmd) preLoop() {
	return
}

func (cmd *Cmd) postLoop() {
	return
}

func (cmd *Cmd) preCmd(line string) string {
	return line
}

func (cmd *Cmd) postCmd(stop bool, line string) bool {
	return stop
}

func (cmd *Cmd) ParseLine(line string) (string, string, string) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", "", line
	} else if line == "?" {
		line = "help " + line[1:]
	}

	splits := strings.Fields(line)
	_cmd, _arg := splits[0], strings.Join(splits[1:], " ")
	return _cmd, _arg, line
}

func (cmd *Cmd) EmptyLine() (bool, error) {
	if cmd.LastCmd != "" {
		return cmd.OneCmd(cmd.LastCmd)
	} else {
		return false, nil
	}
}

func (cmd *Cmd) DefaultCmd(line string) (bool, error) {
	return false, errors.New(fmt.Sprintf("*** Unknown syntax: %s ***", line))
}

func (cmd *Cmd) OneCmd(line string) (bool, error) {
	_cmd, _arg, _line := cmd.ParseLine(line)
	if len(_line) == 0 {
		return cmd.EmptyLine()
	}
	if len(_cmd) == 0 {
		return cmd.DefaultCmd(_line)
	}

	function, found := cmd.FunctionMap[_cmd]
	if found {
		cmd.LastCmd = line
		return function(_arg)
	} else {
		return cmd.DefaultCmd(_line)
	}
}

func (cmd *Cmd) init() {

	if cmd.PreLoop == nil {
		cmd.PreLoop = cmd.preLoop
	}

	if cmd.PreCmd == nil {
		cmd.PreCmd = cmd.preCmd
	}

	if cmd.PostCmd == nil {
		cmd.PostCmd = cmd.postCmd
	}
	if cmd.PostLoop == nil {
		cmd.PostLoop = cmd.postLoop
	}

	if cmd.CommandPrompt == "" {
		cmd.CommandPrompt = "cmd> "
	}
}

func (cmd *Cmd) CommandLoop() {
	cmd.init()
	reader := bufio.NewReader(os.Stdin)
	stopLooping := false

	cmd.PreLoop()
	for !stopLooping {
		fmt.Print(cmd.CommandPrompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}
		line = cmd.PreCmd(strings.TrimSpace(line))
		stop, err := cmd.OneCmd(line)
		if err != nil {
			cmd.Error = err
		}
		stopLooping = cmd.PostCmd(stop, line)
	}
	cmd.PostLoop()
}
