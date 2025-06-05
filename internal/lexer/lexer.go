package lexer

import "fmt"

type TokenType int

const (
	LiteralStr TokenType = iota
	Space
	Redirect
	Append
	Pipe
)

type Token struct {
	Typ     TokenType
	Literal string
}

type Lexer struct {
	input    string
	position int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:    input,
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
			tokens = append(tokens, Token{Typ: LiteralStr, Literal: l.input[start:end]})
			l.position = end + 1
		case '"':
			l.position++
			end := l.position
			curr := ""
			for end < len(l.input) && l.input[end] != '"' {
				if l.input[end] == '\\' {
					if end+1 < len(l.input) {
						next := l.input[end+1]
						if next == '$' || next == '\\' || next == '"' || next == '\n' {
							curr += string(next)
							end = end + 2
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
			tokens = append(tokens, Token{Typ: LiteralStr, Literal: curr})
			l.position = end + 1
		case ' ':
			tokens = append(tokens, Token{Typ: Space, Literal: string(l.input[l.position])})
			for l.position < len(l.input) && l.input[l.position] == ' ' {
				l.position++
			}
		case '\\':
			if l.position+1 < len(l.input) {
				tokens = append(tokens, Token{Typ: LiteralStr, Literal: string(l.input[l.position+1])})
				l.position++
			}
			l.position++
		case '1':
			curr := l.position + 1
			if curr < len(l.input) && l.input[curr] == ' ' {
				curr++
			}
			if curr+1 < len(l.input) && l.input[curr:curr+2] == ">>" {
				tokens = append(tokens, Token{Typ: Append, Literal: "stdout"})
				l.position = curr + 2
			} else if curr < len(l.input) && l.input[curr] == '>' {
				tokens = append(tokens, Token{Typ: Redirect, Literal: "stdout"})
				l.position = curr + 1
			} else {
				tokens = append(tokens, Token{Typ: LiteralStr, Literal: string('1')})
				l.position++
			}
		case '2':
			curr := l.position + 1
			if curr < len(l.input) && l.input[curr] == ' ' {
				curr++
			}
			if curr+1 < len(l.input) && l.input[curr:curr+2] == ">>" {
				tokens = append(tokens, Token{Typ: Append, Literal: "stderr"})
				l.position = curr + 2
			} else if curr < len(l.input) && l.input[curr] == '>' {
				tokens = append(tokens, Token{Typ: Redirect, Literal: "stderr"})
				l.position = curr + 1
			} else {
				tokens = append(tokens, Token{Typ: LiteralStr, Literal: string('2')})
				l.position++
			}
		case '>':
			if l.position+1 < len(l.input) && l.input[l.position+1] == '>' {
				tokens = append(tokens, Token{Typ: Append, Literal: "stdout"})
				l.position += 2
			} else {
				tokens = append(tokens, Token{Typ: Redirect, Literal: "stdout"})
				l.position++
			}
		case '|':
			tokens = append(tokens, Token{Typ: Pipe, Literal: "pipe"})
			l.position++
		default:
			curr := ""
			end := l.position
			for end < len(l.input) && (l.input[end] != ' ' && l.input[end] != '\'' && l.input[end] != '"') {
				if l.input[end] == '\\' {
					if end+1 < len(l.input) {
						next := l.input[end+1]
						curr += string(next)
						end = end + 2
						continue
					}
				}
				curr += string(l.input[end])
				end++
			}
			tokens = append(tokens, Token{Typ: LiteralStr, Literal: string(curr)})
			l.position = end
		}
	}

	return tokens, nil
}
