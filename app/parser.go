package main

import "fmt"

type Redirection struct {
  Type string
  Fd int
  FilePath string
}

type Command struct {
  Name string
  Args []string
  Redirs []Redirection
}
func ParseTokens(tokens []Token) ([]*Command, error) {
  i := 0
  commands := make([]*Command, 0)
  for i < len(tokens) {
    command := &Command{}
    err := error(nil)
    command, i, err = ParseCommand(tokens, i)
    if err != nil {
      return []*Command{}, err
    }
    commands = append(commands, command)
  }
  return commands, nil
}
func ParseCommand(tokens []Token, start int) (*Command, int, error) {
  i := start
  for i < len(tokens) && tokens[i].typ == Space {
    i++
  }
  if tokens[i].typ != LiteralStr {
    return nil, len(tokens), fmt.Errorf("No command provided!")
  }
  name := tokens[i].literal
  i++
  tokens = tokens[i:]
  args := make([]string, 0)
  redirs := make([]Redirection, 0)
  for j := 0; j < len(tokens); j++ {
    if tokens[j].typ == LiteralStr || (name == "echo" && tokens[j].typ == Space && j != 0) {
      args = append(args, tokens[j].literal)
    } else if tokens[j].typ == Redirect {
      redirFd := 1
      redirType := ">"
      switch tokens[j].literal {
      case "stdout":
        redirType = ">"
        redirFd = 1
      case "stderr":
        redirType = ">"
        redirFd = 2
      }
      j++
      for j < len(tokens) && tokens[j].typ == Space {
        j++
      }
      if j >= len(tokens) || tokens[j].typ != LiteralStr {
        return nil, 0, fmt.Errorf("Expected file path for redirect!\n")
      }
      redirs = append(redirs, Redirection{Type: redirType, Fd: redirFd, FilePath: tokens[j].literal})
    } else if tokens[j].typ == Append {
      redirFd := 1
      redirType := ">>"
      switch tokens[j].literal {
      case "stdout":
        redirType = ">>"
        redirFd = 1
      case "stderr":
        redirType = ">>"
        redirFd = 2
      }
      j++
      for j < len(tokens) && tokens[j].typ == Space {
        j++
      }
      if j >= len(tokens) || tokens[j].typ != LiteralStr {
        return nil, 0, fmt.Errorf("Expected file path for redirect!\n")
      }
      redirs = append(redirs, Redirection{Type: redirType, Fd: redirFd, FilePath: tokens[j].literal})
    } else if tokens[j].typ == Pipe {
      if name == "echo" {
        end := len(args)-1
        for ; end >= 0 && args[end] == " "; end-- {}
        args = args[:end+1]
      }
      return &Command{Name: name, Args: args, Redirs: redirs}, j+i+2, nil
    }
  }

  return &Command{Name: name, Args: args, Redirs: redirs}, len(tokens)+start+1, nil
}

