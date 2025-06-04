package parser

import (
  "fmt"
  "github.com/cheesyhypocrisy/harsh/internal/lexer"
)

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

func ParseTokens(tokens []lexer.Token) ([]*Command, error) {
  i := 0
  commands := make([]*Command, 0)
  for i < len(tokens) {
    command := &Command{}
    err := error(nil)
    command, i, err = ParseCommand(tokens, i)
    if err != nil {
      return []*Command{}, err
    }
    if command != nil {
      commands = append(commands, command)
    }
  }

  if len(commands) == 0 {
    return []*Command{}, fmt.Errorf("No command provided!")
  }
  return commands, nil
}

func ParseCommand(tokens []lexer.Token, start int) (*Command, int, error) {
  i := start
  for i < len(tokens) && tokens[i].Typ == lexer.Space {
    i++
  }
  if i >= len(tokens) {
    return nil, len(tokens), nil
  }
  if tokens[i].Typ != lexer.LiteralStr {
    return nil, len(tokens), fmt.Errorf("No command provided!")
  }
  name := tokens[i].Literal
  i++
  tokens = tokens[i:]
  args := make([]string, 0)
  redirs := make([]Redirection, 0)
  for j := 0; j < len(tokens); j++ {
    if tokens[j].Typ == lexer.LiteralStr || (name == "echo" && tokens[j].Typ == lexer.Space && j != 0) {
      args = append(args, tokens[j].Literal)
    } else if tokens[j].Typ == lexer.Redirect {
      redirFd := 1
      redirType := ">"
      switch tokens[j].Literal {
      case "stdout":
        redirType = ">"
        redirFd = 1
      case "stderr":
        redirType = ">"
        redirFd = 2
      }
      j++
      for j < len(tokens) && tokens[j].Typ == lexer.Space {
        j++
      }
      if j >= len(tokens) || tokens[j].Typ != lexer.LiteralStr {
        return nil, 0, fmt.Errorf("Expected file path for redirect!\n")
      }
      redirs = append(redirs, Redirection{Type: redirType, Fd: redirFd, FilePath: tokens[j].Literal})
    } else if tokens[j].Typ == lexer.Append {
      redirFd := 1
      redirType := ">>"
      switch tokens[j].Literal {
      case "stdout":
        redirType = ">>"
        redirFd = 1
      case "stderr":
        redirType = ">>"
        redirFd = 2
      }
      j++
      for j < len(tokens) && tokens[j].Typ == lexer.Space {
        j++
      }
      if j >= len(tokens) || tokens[j].Typ != lexer.LiteralStr {
        return nil, 0, fmt.Errorf("Expected file path for redirect!\n")
      }
      redirs = append(redirs, Redirection{Type: redirType, Fd: redirFd, FilePath: tokens[j].Literal})
    } else if tokens[j].Typ == lexer.Pipe {
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

