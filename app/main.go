package main

import (
	"bufio"
	"fmt"
	"os"
  "strconv"
  "strings"
)

func eval(command string, args []string) (string, error) {
    switch command {
    case "exit":
      if len(args) == 0 {
        os.Exit(0)
      }
      exitCode, err := strconv.Atoi(args[0])
      if err != nil {
        return "", err
      }
      os.Exit(exitCode)
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
      fmt.Println(result)
    }
  }
}

func main() {
  if err := shell(); err != nil {
    fmt.Fprint(os.Stderr, err)
  }
}

