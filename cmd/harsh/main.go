package main

import (
  "os"
  "fmt"
  "strings"
  "bufio"

  "github.com/cheesyhypocrisy/harsh/internal/executor"
	"github.com/cheesyhypocrisy/harsh/internal/shell"
)

func main() {
  path := os.Getenv("PATH")
  executor.PathDirs = append(executor.PathDirs, strings.Split(path, ":")...)

  histfile, exists := os.LookupEnv("HISTFILE")
  // Only enable history persistence if HISTFILE is set
  if exists {
    file, err := os.Open(histfile)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Unable to read history from file %s with err: %#v\n", histfile, err.Error())
      return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      executor.Hist = append(executor.Hist, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
      fmt.Fprintf(os.Stderr, "Unable to read history from file %s with err: %#v\n", histfile, err.Error())
    }
  }


  if err := shell.Shell(); err != nil {
    fmt.Fprint(os.Stderr, err)
  }
}

