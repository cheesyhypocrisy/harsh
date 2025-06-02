package main

import (
	"bufio"
	"fmt"
	"os"
  "strconv"
  "strings"
)

type builtin int

const (
  unknownBuiltin builtin = iota
  exit
  echo
  _type
)

func lookupBuiltin(command string) builtin {
  switch command {
  case "exit":
    return exit
  case "echo":
    return echo
  case "type":
    return _type
  default:
    return unknownBuiltin
  }
}

func eval(command string, args []string) (string, error) {
    switch lookupBuiltin(command) {
    case exit:
      if len(args) == 0 {
        os.Exit(0)
      }
      exitCode, err := strconv.Atoi(args[0])
      if err != nil {
        return "", err
      }
      os.Exit(exitCode)
    case echo:
      return strings.Join(args, " ")+"\n", nil
    case _type:
      if len(args) == 0 {
        return "", fmt.Errorf("Missing argument for type command\n")
      }
      if lookupBuiltin(args[0]) == unknownBuiltin {
        return "", fmt.Errorf("%s: not found\n", args[0])
      } else {
        return fmt.Sprintf("%s is a shell builtin\n", args[0]), nil
      }
    default:
      return "", fmt.Errorf("%s: command not found\n",strings.TrimSpace(command))
    }

    return "", nil
}

func shell() error {
	for {
    fmt.Fprint(os.Stdout, "$ ")
  
    fullCommand, err := bufio.NewReader(os.Stdin).ReadString('\n')
    if err != nil {
      return err
    }
    
    fullCommand = strings.TrimSpace(fullCommand)
    cliArgs := strings.Split(fullCommand, " ")
    command := cliArgs[0]
    args := cliArgs[1:]

    result, err := eval(command, args)
    
    if err != nil {
      fmt.Fprint(os.Stderr, err)
    } else {
      fmt.Print(result)
    }
  }
}

func main() {
  if err := shell(); err != nil {
    fmt.Fprint(os.Stderr, err)
  }
}

