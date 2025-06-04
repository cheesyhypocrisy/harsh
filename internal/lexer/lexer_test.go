package lexer

import (
  "testing"
)

func TestLexer(t *testing.T) {
  tests := []struct {
    name     string
    input    string
    expected []Token
    hasError bool
  }{
    {
      name: "Simple Command",
      input: "echo hello",
      expected: []Token{
        {Typ: LiteralStr, Literal: "echo"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "hello"},
      },
      hasError: false,
    },
    {
      name: "Command with quoted strings",
      input: "echo 'hello world' \"quoted\" \"test-\"\"concat\"",
      expected: []Token{
        {Typ: LiteralStr, Literal: "echo"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "hello world"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "quoted"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "test-"},
        {Typ: LiteralStr, Literal: "concat"},
      },
      hasError: false,
    },
    {
      name: "Command with redirection",
      input: "echo 'hello world' > output.txt",
      expected: []Token{
        {Typ: LiteralStr, Literal: "echo"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "hello world"},
        {Typ: Space, Literal: " "},
        {Typ: Redirect, Literal: "stdout"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "output.txt"},
      },
      hasError: false,
    },
    {
      name:     "Unmatched quote",
      input:    "echo 'hello",
      expected: []Token{},
      hasError: true,
    },
    {
      name:  "Command with pipe",
      input: "echo hello | grep hello",
      expected: []Token{
        {Typ: LiteralStr, Literal: "echo"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "hello"},
        {Typ: Space, Literal: " "},
        {Typ: Pipe, Literal: "pipe"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "grep"},
        {Typ: Space, Literal: " "},
        {Typ: LiteralStr, Literal: "hello"},
      },
      hasError: false,
    },
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T){
      lexer := NewLexer(test.input)
      tokens, err := lexer.Lex()

      if (err != nil) != test.hasError {
        t.Errorf("Error expectation mismatch - got error: %v, expected error: %v", err, test.hasError)
        return
      }

      if test.hasError {
        return
      }

      if len(tokens) != len(test.expected) {
        t.Errorf("Expected %d tokens, got %d\n", len(test.expected), len(tokens))
      }

      for i, token := range tokens {
        if token.Typ != test.expected[i].Typ || token.Literal != test.expected[i].Literal {
          t.Errorf("Token %d - Expected %+v, got %+v\n", i, test.expected[i], token)
        }
      }
    })
  }
}

