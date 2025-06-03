package main

import (
  "os"
  "fmt"
  "strings"

  "github.com/cheesyhypocrisy/harsh/internal/executor"
	"github.com/cheesyhypocrisy/harsh/internal/shell"
)

func main() {
  path := os.Getenv("PATH")
  executor.PathDirs = append(executor.PathDirs, strings.Split(path, ":")...)

  if err := shell.Shell(); err != nil {
    fmt.Fprint(os.Stderr, err)
  }
}

