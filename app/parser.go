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

func ParseTokens(tokens []Token) (*Command, error) {
  if tokens[0].typ != LiteralStr {
    return nil, fmt.Errorf("No command provided!")
  }
  name := tokens[0].literal
  tokens = tokens[1:]
  args := make([]string, 0)
  redirs := make([]Redirection, 0)
  for i := 0; i < len(tokens); i++ {
    if tokens[i].typ == LiteralStr || (name == "echo" && tokens[i].typ == Space && i != 0) {
      args = append(args, tokens[i].literal)
    } else if tokens[i].typ == Redirect {
      redirFd := 1
      redirType := ">"
      switch tokens[i].literal {
      case "stdout":
        redirType = ">"
        redirFd = 1
      case "stderr":
        redirType = ">"
        redirFd = 2
      }
      i++
      for i < len(tokens) && tokens[i].typ == Space {
        i++
      }
      if i >= len(tokens) || tokens[i].typ != LiteralStr {
        return nil, fmt.Errorf("Expected file path for redirect!\n")
      }
      redirs = append(redirs, Redirection{Type: redirType, Fd: redirFd, FilePath: tokens[i].literal})
    } else if tokens[i].typ == Append {
      redirFd := 1
      redirType := ">>"
      switch tokens[i].literal {
      case "stdout":
        redirType = ">>"
        redirFd = 1
      case "stderr":
        redirType = ">>"
        redirFd = 2
      }
      i++
      for i < len(tokens) && tokens[i].typ == Space {
        i++
      }
      if i >= len(tokens) || tokens[i].typ != LiteralStr {
        return nil, fmt.Errorf("Expected file path for redirect!\n")
      }
      redirs = append(redirs, Redirection{Type: redirType, Fd: redirFd, FilePath: tokens[i].literal})
    }
  }

  return &Command{Name: name, Args: args, Redirs: redirs}, nil
}

