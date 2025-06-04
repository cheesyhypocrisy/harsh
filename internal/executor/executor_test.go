package executor

import (
  "testing"
  "bytes"
  "strings"

  "github.com/cheesyhypocrisy/harsh/internal/parser"
)

func TestLookupBuiltin(t *testing.T) {
  tests := []struct {
    name     string
    command  string
    expected builtin
  }{
    {"exit command", "exit", exit},
    {"echo command", "echo", echo},
    {"type command", "type", _type},
    {"pwd command", "pwd", pwd},
    {"cd command", "cd", cd},
    {"history command", "history", history},
    {"unknown command", "unknown", unknownBuiltin},
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      result := lookupBuiltin(test.command)
      if result != test.expected {
        t.Errorf("lookupBuiltin(%s) = %v, expected %v", test.command, result, test.expected)
      }
    })
  }
}

func TestFindExecutable(t *testing.T) {
  originalPathDirs := PathDirs
  defer func() { PathDirs = originalPathDirs }()

  tests := []struct {
    name      string
    pathDirs  []string
    command   string
    shouldErr bool
  }{
    {
      name:      "Command exists",
      pathDirs:  []string{"/bin", "/usr/bin"},
      command:   "ls",
      shouldErr: false,
    },
    {
      name:      "Command doesn't exist",
      pathDirs:  []string{"/bin", "/usr/bin"},
      command:   "unknown",
      shouldErr: true,
    },
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      PathDirs = test.pathDirs

      path, err := findExecutable(test.command)

      if (err != nil) != test.shouldErr {
        t.Errorf("Error expectation mismatch - got error: %v, expected error: %v", err, test.shouldErr)
        return
      }

      if !test.shouldErr && !strings.HasSuffix(path, test.command) {
        t.Errorf("Expected path to end with %s, got %s", test.command, path)
      }
    })
  }
}

func TestWrapBuiltin(t *testing.T) {
  tests := []struct {
    name           string
    command        *parser.Command
    expectedOutput string
    expectedError  string
  }{
    {
      name: "Echo command",
      command: &parser.Command{
        Name: "echo",
        Args: []string{"hello", "world"},
      },
      expectedOutput: "helloworld\n",
      expectedError:  "",
    },
    {
      name: "Type command with builtin",
      command: &parser.Command{
        Name: "type",
        Args: []string{"echo"},
      },
      expectedOutput: "echo is a shell builtin\n",
      expectedError:  "",
    },
    {
      name: "Type command with no args",
      command: &parser.Command{
        Name: "type",
        Args: []string{},
      },
      expectedOutput: "",
      expectedError:  "Missing argument for type command\n",
    },
  }

  for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
      runnable := WrapBuiltin(test.command)

      if !runnable.isBuiltin {
        t.Errorf("Expected isBuiltin to be true")
      }

      stdout := &bytes.Buffer{}
      stderr := &bytes.Buffer{}

      runnable.Start(nil, stdout, stderr)

      if stdout.String() != test.expectedOutput {
        t.Errorf("Expected output %q, got %q", test.expectedOutput, stdout.String())
      }

      if stderr.String() != test.expectedError {
        t.Errorf("Expected stderr %q, got %q", test.expectedError, stderr.String())
      }
    })
  }
}
