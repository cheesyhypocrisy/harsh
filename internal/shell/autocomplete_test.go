package shell

import (
  "testing"

  "github.com/cheesyhypocrisy/harsh/internal/executor"
)

func TestAutocomplete(t *testing.T) {
  originalPathDirs := executor.PathDirs
  defer func() { executor.PathDirs = originalPathDirs }()

  tests := []struct {
    name            string
    line            string
    pos             int
    pathDirs        []string
    expectSuggestion bool
    expectedLength   int
  }{
    {
      name:            "Partial command 'ec'",
      line:            "ec",
      pos:             2,
      pathDirs:        []string{"/bin", "/usr/bin"},
      expectSuggestion: true,
      expectedLength:   2,
    },
    {
      name:            "Empty input",
      line:            "",
      pos:             0,
      pathDirs:        []string{"/bin", "/usr/bin"},
      expectSuggestion: false,
      expectedLength:   0,
    },
    {
      name:            "Complete command 'echo'",
      line:            "echo",
      pos:             4,
      pathDirs:        []string{"/bin", "/usr/bin"},
      expectSuggestion: false,
      expectedLength:   4,
    },
    {
      name:            "Command with space 'echo '",
      line:            "echo ",
      pos:             5,
      pathDirs:        []string{"/bin", "/usr/bin"},
      expectSuggestion: false,
      expectedLength:   0,
    },
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      executor.PathDirs = test.pathDirs

      autocomplete := &Autocomplete{}

      suggestions, length := autocomplete.Do([]rune(test.line), test.pos)

      if test.expectSuggestion && len(suggestions) == 0 {
        t.Errorf("Expected suggestions, got none")
      }

      if length != test.expectedLength {
        t.Errorf("Expected length %d, got %d", test.expectedLength, length)
      }
    })
  }
}
