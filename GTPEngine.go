package go_GTPEngine

import (
	"fmt"
	"strconv"
	"strings"
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
	name         string
	version      string

	cmd_idx  int
	response string

	Clear   func()
	GenMove func(StoneType) (int, int)
	Play    func(StoneType, int, int) error
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

		e.response = "=" + strconv.Itoa(e.cmd_idx) + " "

		cmd := e.ParseCmd(line)
		e.Process(cmd)

		fmt.Print(e.response, "\n\n")
		e.cmd_idx++
	}
}

func (e *Engine) ParseCmd(line string) Cmd {
	if len(line) == 0 {
		return Cmd{Type: CMD_QUIT}
	}

	tokens := strings.Fields(line)

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

func (e *Engine) Process(cmd Cmd) {
	switch cmd.Type {
	case CMD_QUIT:
		e.running = false

	case CMD_PROTOCOL_VERSION:
		e.response += e.protocol_ver
	case CMD_NAME:
		e.response += e.name
	case CMD_VERSION:
		e.response += e.version

	case CMD_CLEAR_BOARD:
		if e.Clear != nil {
			e.Clear()
		}

	case CMD_PLAY:
		if len(cmd.Args) <= 1 {
			e.response = "?" + e.response[1:] + "syntax error"
			break
		}

		if e.Play != nil {
			stone := parseColor(cmd.Args[0])
			x, y := parseCoordinate(cmd.Args[1])

			if !(stone > 0 && x >= -1 && y >= -1) {
				e.response = "?" + e.response[1:] + "invalid color or coordinate"
			}

			err := e.Play(stone, x, y)

			if err != nil {
				e.response = "?" + e.response[1:] + "illegal move"
				break
			}
		}

	case CMD_GENMOVE:
		if e.GenMove != nil {
			stone := parseColor(cmd.Args[0])

			x, y := e.GenMove(stone)

			if x == -1 && y == -1 {
				e.response += "pass"
			} else if x == -2 && y == -2 {
				e.response += "resign"
			} else {
				pos_list := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

				e.response += pos_list[x:x]
				e.response += strconv.Itoa(19 - y)
			}
		}
	}
}

func parseColor(arg string) StoneType {
	switch strings.ToLower(arg) {
	case "black":
		fallthrough
	case "b":
		return BLACK

	case "white":
		fallthrough
	case "w":
		return WHITE

	default:
		return NONE
	}
}

func parseCoordinate(arg string) (int, int) {
	if strings.ToLower(arg) == "pass" {
		return -1, -1
	}

	pos_list := "abcdefghijklmnopqrstuvwxyz"

	x := strings.Index(pos_list, strings.ToLower(arg)[:0])
	y, _ := strconv.Atoi(arg[1:])

	return x, 19 - y
}
