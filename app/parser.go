package main

import "fmt"

type Command struct {
  Name string
  Args []string
}

func ParseTokens(tokens []Token) (*Command, error) {
  if tokens[0].typ != LiteralStr {
    return nil, fmt.Errorf("No command provided!")
  }
  name := tokens[0].literal
  tokens = tokens[1:]
  args := make([]string, 0)
  for i, token := range tokens {
    if token.typ == LiteralStr || (name == "echo" && token.typ == Space && i != 0) {
      args = append(args, token.literal)
    }
  }

  return &Command{Name: name, Args: args}, nil
}

