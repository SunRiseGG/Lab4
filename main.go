package main

import (
	"fmt"
	"./engine"
	"bufio"
	"strings"
	"os"
)


type printCommand struct {
	arg string
}

func (print *printCommand) Execute(loop engine.Handler) {
	fmt.Println(print.arg)
}

type deleteCommand struct {
	str, symbol string
}

func (delete *deleteCommand) Execute(loop engine.Handler) {
	result := strings.ReplaceAll(delete.str, delete.symbol, "");
	loop.Post(&printCommand{arg: result})
}

func parse(commandLine string) engine.Command {
	pos1 := strings.Index(commandLine, string('"'))
	pos2 := strings.LastIndex(commandLine, string('"'))
	substr := ""
	if pos1 != -1 && pos2 != -1 {
	  substr = commandLine[pos1:pos2 + 1]
		commandLine = strings.ReplaceAll(commandLine, substr, "")
		substr = strings.TrimLeft(substr, string('"'))
    substr = strings.TrimRight(substr, string('"'))
	}
  elements := strings.Fields(commandLine)
  if elements[0] == "print" {
		if substr != "" {
			if len(elements) == 1 {
        return &printCommand{arg: substr}
			} else {
				return &printCommand{arg: "SYNTAX ERROR: Print has too many arguments"}
			}
		} else if len(elements) == 2 {
      return &printCommand{arg: elements[1]}
    } else if len(elements) == 1 {
      return &printCommand{arg: "SYNTAX ERROR: Print doesn't have enough arguments"}
    } else {
      return &printCommand{arg: "SYNTAX ERROR: Print has too many arguments"}
    }
  } else if elements[0] == "delete" {
		if substr != "" {
			if len(elements) == 2 {
        return &deleteCommand{str: substr, symbol: elements[1]}
			} else if len(elements) < 2 {
	      return &printCommand{arg: "SYNTAX ERROR: Delete doesn't have enough arguments"}
	    }	else {
				return &printCommand{arg: "SYNTAX ERROR: Delete has too many arguments"}
			}
		} else if len(elements) == 3 {
      return &deleteCommand{str: elements[1], symbol: elements[2]}
    } else if len(elements) < 3 {
      return &printCommand{arg: "SYNTAX ERROR: Delete doesn't have enough arguments"}
    } else {
      return &printCommand{arg: "SYNTAX ERROR: Delete has too many arguments"}
    }
  } else {
    return &printCommand{arg: "SYNTAX ERROR: Unknown command"}
  }
}

func main() {
	eventLoop := new(engine.Loop)
	eventLoop.Start()
	if input, err := os.Open("./input.txt"); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd := parse(commandLine)
			eventLoop.Post(cmd)
		}
	}
	eventLoop.AwaitFinish()
}
