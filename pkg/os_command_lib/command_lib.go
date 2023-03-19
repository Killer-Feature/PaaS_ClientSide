package os_command_lib

import (
	"strings"
)

type Code string

const (
	ResourceAlreadyBusy Code = "16"
)

type Command string

type Parser func(output []byte, extraData interface{}) error

type CommandAndParser struct {
	Command   Command
	Parser    Parser
	Condition Condition
}

type Condition uint8

const (
	Sufficient Condition = iota
	Required
	Anyway
)

func (c CommandAndParser) WithArgs(args ...string) CommandAndParser {
	c.Command = c.Command + " " + Command(strings.Join(args, " "))
	return c
}

func (c CommandAndParser) String() string {
	return string(c.Command)
}

func (c CommandAndParser) WithEnv(env, val, sep string) CommandAndParser {
	c.Command = Command(env + " = " + val + sep + c.String())
	return c
}

func (c CommandAndParser) Pipe(cmds ...CommandAndParser) CommandAndParser {

	strCmds := make([]string, 0, 1+len(cmds))
	strCmds = append(strCmds, c.String())
	for i := range cmds {
		strCmds = append(strCmds, cmds[i].String())
	}

	c.Command = Command(strings.Join(strCmds, " | "))
	c.Parser = cmds[len(cmds)-1].Parser
	c.Condition = cmds[len(cmds)-1].Condition
	return c
}
