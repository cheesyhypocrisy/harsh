package shell

import (
  "fmt"
  "strings"
  "unicode"
  "os"
  "sort"

  "github.com/cheesyhypocrisy/harsh/internal/executor"
)

type Autocomplete struct {
  lastPos int
  tabCount int
  lastLine string
}

func (a *Autocomplete) Do(line []rune, pos int) (newLine [][]rune, length int) {
  if a.lastLine == string(line) && a.lastPos == pos {
    fmt.Println(a.tabCount)
    a.tabCount = 1
  } else {
    a.tabCount = 0
  }

  a.lastLine = string(line)
  a.lastPos = pos
  start := pos
  for start > 0 && !unicode.IsSpace(line[start-1]) {
    start--
  }

  input := string(line[start:pos])

  suggestions := [][]rune{}
  builtins := []string{"exit", "echo", "type", "pwd", "cd"}

  commandsSet := make(map[string]bool)
  for _, builtin := range builtins {
    commandsSet[builtin] = true
  }
  for _, dir := range executor.PathDirs { 
    files, err := os.ReadDir(dir)
    if err != nil {
      continue
    }
    for _, file := range files {
      if !file.IsDir() {
        commandsSet[file.Name()] = true
      }
    } 
  }

  for command := range commandsSet {
    if strings.HasPrefix(command, input) {
      suggestions = append(suggestions, []rune(command[len(input):]+" "))
    }
  }

  if len(suggestions) == 0 {
    fmt.Fprintf(os.Stdout, "\x07")
  }

  if len(suggestions) > 1 {
    sort.Slice(suggestions, func(i, j int) bool {
      return string(suggestions[i]) < string(suggestions[j])
    })

    smallest := suggestions[0][:len(suggestions[0])-1]
    isCommonPrefix := true
    for i := 1; i < len(suggestions); i++ {
      if !strings.HasPrefix(string(suggestions[i][:len(suggestions[i])-1]),string(smallest)) {
        isCommonPrefix = false
      }
    }
    if isCommonPrefix {
      return [][]rune{smallest}, pos-start
    }

    if a.tabCount == 0 {
      a.tabCount++
      fmt.Fprintf(os.Stdout, "\a")
      return [][]rune{}, 0
    } else if a.tabCount == 1 {
      a.tabCount = 0
    }

    fmt.Fprintf(os.Stdout,"\r\n")
    for i, suggestion := range suggestions {
      fmt.Fprintf(os.Stdout, "%s", input + string(suggestion))
      if i != len(suggestions)-1 {
        fmt.Fprintf(os.Stdout, " ")
      }
    }
    fmt.Printf("\n")
    fmt.Printf("$ %s", string(line))
    return [][]rune{}, 0 // Don't want suggestions to be tabbable hence handled above
  }

  return suggestions, pos-start
}

