package go_GTPEngine

import (
	"strings"
	"fmt"
)

type CommandType int

const (
	CMD_QUIT CommandType = iota
	CMD_PROTOCOL_VERSION
	CMD_NAME
	CMD_VERSION
	CMD_CLEAR_BOARD
	CMD_PLAY
	CMD_GENMOVE
	CMD_SHOWBOARD
	CMD_UNKNOWN
)

type StoneType int

const (
	NONE StoneType = iota
	BLACK
	WHITE
)

type Cmd struct {
	Type CommandType
	Args []string
}

type Engine struct {
	running bool

	protocol_ver string
	name string
	version string

	Clear func()
	GenMove func(StoneType) (int, int)
	Play func(StoneType, int, int)
}

func InitEngine(protocol_ver, name, version string) *Engine {
	engine := new(Engine)

	engine.running = false

	engine.protocol_ver = protocol_ver
	engine.name = name
	engine.version = version

	return engine
}

func (e *Engine) Run() {
	e.running = true

	for e.running {
		var line string
		fmt.Scanln(line)

		cmd := e.ParseCmd(line)
		response := e.Process(cmd)

		fmt.Println(response, "\n")
	}
}

func (e* Engine) ParseCmd(line string) Cmd {
	if len(line) == 0 {
		return Cmd{Type: CMD_QUIT}
	}

	tokens := strings.Split(line, ' ')

	ret := Cmd{Type: CMD_UNKNOWN, Args: tokens[1:]}

	switch strings.ToLower(tokens[0]) {
	case "quit":
		ret.Type = CMD_QUIT
	case "protocol_version":
		ret.Type = CMD_PROTOCOL_VERSION
	case "name":
		ret.Type = CMD_NAME
	case "version":
		ret.Type = CMD_VERSION
	case "clear_board":
		ret.Type = CMD_CLEAR_BOARD
	case "play":
		ret.Type = CMD_PLAY
	case "genmove":
		ret.Type = CMD_GENMOVE
	case "showboard":
		ret.Type = CMD_SHOWBOARD
	}

	return ret
}

func (e *Engine) Process(cmd Cmd) string {
	switch cmd.Type {
	case CMD_QUIT:
		e.running = false
		return ""

	case CMD_PROTOCOL_VERSION:
		return e.protocol_ver
	case CMD_NAME:
		return e.name
	case CMD_VERSION:
		return e.version

	case CMD_CLEAR_BOARD:
		if e.Clear != nil {
			e.Clear()
		}

		return ""

	case CMD_PLAY:
		if e.Play != nil {
		}

		return ""

	case CMD_GENMOVE:
		if e.GenMove != nil {

		}
		return ""

	default:
		return ""
	}
}
