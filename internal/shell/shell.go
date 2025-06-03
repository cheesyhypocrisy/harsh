package shell

import (
  "strings"

  "github.com/cheesyhypocrisy/harsh/internal/executor"
	"github.com/cheesyhypocrisy/harsh/internal/lexer"
	"github.com/cheesyhypocrisy/harsh/internal/parser"
  "github.com/chzyer/readline"
)

func Shell() error {
  autocomplete := &Autocomplete{
    tabCount: 0,
  }
  rl, err := readline.NewEx(&readline.Config{
    Prompt: "$ ",
    AutoComplete: autocomplete,
    InterruptPrompt: "^C",
    EOFPrompt:       "exit",
  })
  if err != nil {
    return err
  }
  defer rl.Close()

	for {
    line, err := rl.Readline()
    if err != nil {
      return err
    }
    
    line = strings.TrimSpace(line)
    executor.Hist = append(executor.Hist, line)
    tokens, err := lexer.NewLexer(line).Lex()
    if err != nil {
      return err
    }

    commands, err := parser.ParseTokens(tokens)
    if err != nil {
      return err
    }
    executor.Eval(commands)
  }
}
