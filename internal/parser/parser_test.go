package parser

import (
  "testing"

  "github.com/cheesyhypocrisy/harsh/internal/lexer"
)

func TestParseTokens(t *testing.T) {
  tests := []struct {
    name string
    tokens []lexer.Token
    expected []*Command
    hasError bool
  }{
    {
      name: "Simple Command",
      tokens: []lexer.Token{
        {Typ: lexer.LiteralStr, Literal: "echo"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "hello"},
      },
      expected: []*Command{
        {
          Name:   "echo",
          Args:   []string{"hello"},
          Redirs: []Redirection{},
        },
      },
      hasError: false,
    },
    {
      name: "Command with redirection",
      tokens: []lexer.Token{
        {Typ: lexer.LiteralStr, Literal: "echo"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "hello"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.Redirect, Literal: "stdout"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "output.txt"},
      },
      expected: []*Command{
        {
          Name: "echo",
          Args: []string{"hello", " "},
          Redirs: []Redirection{
            {Type: ">", Fd: 1, FilePath: "output.txt"},
          },
        },
      },
      hasError: false,
    },
    {
      name: "Command with append redirection",
      tokens: []lexer.Token{
        {Typ: lexer.LiteralStr, Literal: "echo"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "hello"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.Append, Literal: "stdout"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "output.txt"},
      },
      expected: []*Command{
        {
          Name: "echo",
          Args: []string{"hello", " "},
          Redirs: []Redirection{
            {Type: ">>", Fd: 1, FilePath: "output.txt"},
          },
        },
      },
      hasError: false,
    },
    {
      name: "No command",
      tokens: []lexer.Token{
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.Space, Literal: " "},
      },
      expected: []*Command{},
      hasError: true,
    },
    {
      name: "Pipe command",
      tokens: []lexer.Token{
        {Typ: lexer.LiteralStr, Literal: "echo"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "hello"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.Pipe, Literal: "pipe"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "grep"},
        {Typ: lexer.Space, Literal: " "},
        {Typ: lexer.LiteralStr, Literal: "hello"},
      },
      expected: []*Command{
        {
          Name:   "echo",
          Args:   []string{"hello"},
          Redirs: []Redirection{},
        },
        {
          Name:   "grep",
          Args:   []string{"hello"},
          Redirs: []Redirection{},
        },
      },
      hasError: false,
    },
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      commands, err := ParseTokens(test.tokens)

      if (err != nil) != test.hasError {
        t.Errorf("Error expectation mismatch - got error: %v, expected error: %v", err, test.hasError)
        return
      }

      if test.hasError {
        return
      }

      if len(commands) != len(test.expected) {
        t.Errorf("Expected %d commands, got %d", len(test.expected), len(commands))
        return
      }

      for i, cmd := range commands {
        expCmd := test.expected[i]

        if cmd.Name != expCmd.Name {
          t.Errorf("Command %d - Expected name %s, got %s", i, expCmd.Name, cmd.Name)
        }

        if len(cmd.Args) != len(expCmd.Args) {
          t.Errorf("Command %d - Expected %d args, got %d", i, len(expCmd.Args), len(cmd.Args))
          continue
        }

        for j, arg := range cmd.Args {
          if arg != expCmd.Args[j] {
            t.Errorf("Command %d, Arg %d - Expected %s, got %s", i, j, expCmd.Args[j], arg)
          }
        }

        if len(cmd.Redirs) != len(expCmd.Redirs) {
          t.Errorf("Command %d - Expected %d redirections, got %d", i, len(expCmd.Redirs), len(cmd.Redirs))
          continue
        }

        for j, redir := range cmd.Redirs {
          expRedir := expCmd.Redirs[j]
          if redir.Type != expRedir.Type || redir.Fd != expRedir.Fd || redir.FilePath != expRedir.FilePath {
            t.Errorf("Command %d, Redirection %d - Expected %+v, got %+v", i, j, expRedir, redir)
          }
        }
      }
    })
  }
}
