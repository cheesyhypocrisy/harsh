package main

import "fmt"

type TokenType int

const (
  LiteralStr TokenType = iota
  Space
)

type Token struct {
  typ TokenType
  literal string
}

type Lexer struct {
  input string
  position int
}

func NewLexer(input string) *Lexer {
  return &Lexer{
    input: input,
    position: 0,
  }
}

func (l *Lexer) Lex() ([]Token, error) {
  tokens := []Token{}
  for l.position < len(l.input) {
    switch l.input[l.position] {
      case '\'':
        l.position++
        start := l.position
        end := l.position
        for end < len(l.input) && l.input[end] != '\'' {
          end++
        }
        if end == len(l.input) {
          return []Token{}, fmt.Errorf("Unmatched ', expected ' at the end of the input")
        }
        tokens = append(tokens, Token{typ: LiteralStr, literal: l.input[start:end]})
        l.position = end+1
      case '"':
        l.position++
        end := l.position
        curr := ""
        for end < len(l.input) && l.input[end] != '"' {
          if l.input[end] == '\\' {
            if end + 1 < len(l.input) {
              next := l.input[end+1]
              if next == '$' || next == '\\' || next == '"' || next == '\n' {
                curr += string(next)
                end = end+2
                continue
              }
            }
          }
          curr += string(l.input[end])
          end++
        }
        if end == len(l.input) {
          return []Token{}, fmt.Errorf("Unmatched \", expected \" at the end of the input")
        }
        tokens = append(tokens, Token{typ: LiteralStr, literal: curr})
        l.position = end+1
      case ' ':
        tokens = append(tokens, Token{typ: Space, literal: string(l.input[l.position])})
        for l.position < len(l.input) && l.input[l.position] == ' ' {
          l.position++
        }
      case '\\':
        if l.position + 1 < len(l.input) {
          tokens = append(tokens, Token{typ: LiteralStr, literal: string(l.input[l.position+1])})
          l.position++
        }
        l.position++
      default:
        curr := ""
        end := l.position
        for end < len(l.input) && (l.input[end] != ' ' && l.input[end] != '\'') {
          if l.input[end] == '\\' {
            if end + 1 < len(l.input) {
              next := l.input[end+1]
              curr += string(next)
              end = end+2
              continue
            }
          }
          curr += string(l.input[end])
          end++
        }
        tokens = append(tokens, Token{typ: LiteralStr, literal: string(curr)})
        l.position = end
    }
  }

  return tokens, nil
}

