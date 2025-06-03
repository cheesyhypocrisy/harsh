package main

import (
  "github.com/chzyer/readline"
	"fmt"
	"os"
  "os/exec"
  "strconv"
  "strings"
)

var pathDirs []string

type builtin int

const (
  unknownBuiltin builtin = iota
  exit
  echo
  _type
  pwd
  cd
)

func lookupBuiltin(command string) builtin {
  switch command {
  case "exit":
    return exit
  case "echo":
    return echo
  case "type":
    return _type
  case "pwd":
    return pwd
  case "cd":
    return cd
  default:
    return unknownBuiltin
  }
}

func findExecutable(command string) (string, error) {
  path := ""
  for _, dir:= range pathDirs {
    path = strings.TrimRight(dir, "/") + "/" + command
    _, err := os.Stat(path)
    if err == nil {
      return path, nil
    } else {
      continue
    }
  }

  return "", fmt.Errorf("Executable not found in PATH: %s", command)
}

func eval(command *Command) (string, error) {
    switch lookupBuiltin(command.Name) {
    case exit:
      if len(command.Args) == 0 {
        os.Exit(0)
      }
      exitCode, err := strconv.Atoi(command.Args[0])
      if err != nil {
        return "", err
      }
      os.Exit(exitCode)
    case echo:
      return strings.Join(command.Args, "")+"\n", nil
    case _type:
      if len(command.Args) == 0 {
        return "", fmt.Errorf("Missing argument for type command\n")
      }
      if lookupBuiltin(command.Args[0]) == unknownBuiltin {
        path, err := findExecutable(command.Args[0])
        if err != nil {
          return "", fmt.Errorf("%s: not found\n", command.Args[0])
        } else {
          return fmt.Sprintf("%s is %s\n", command.Args[0], path), nil
        }
      } else {
        return fmt.Sprintf("%s is a shell builtin\n", command.Args[0]), nil
      }
    case pwd:
      dir, err := os.Getwd()
      return dir + "\n", err
    case cd:
      if len(command.Args) == 0 || command.Args[0] == "~" {
        homeDir, exists := os.LookupEnv("HOME")
        if !exists {
          username := os.Getenv("USER")
          homeDir = fmt.Sprintf("/home/%s", username)
        }

        if err := os.Chdir(homeDir); err != nil {
          return "", fmt.Errorf("cd: %s: No such file or directory\n", homeDir)
        }

        return "", nil
      }
      if err := os.Chdir(command.Args[0]); err != nil {
        return "", fmt.Errorf("cd: %s: No such file or directory\n", command.Args[0]) 
      }
      return "", nil
    default:
      _, err := findExecutable(command.Name)
      if err == nil {
        cmd := exec.Command(command.Name, command.Args...)
        for _, redir := range command.Redirs {
          if redir.Type == ">" {
            if redir.Fd == 1 {
              file, err := os.Create(redir.FilePath)
              if err != nil {
                return "", err
              }
              defer file.Close()
              cmd.Stdout = file
            } else if redir.Fd == 2 {
              file, err := os.Create(redir.FilePath)
              if err != nil {
                return "", err
              }
              defer file.Close()
              cmd.Stderr = file
            }
          } else if redir.Type == ">>" {
            if redir.Fd == 1 {
              file, err := os.OpenFile(redir.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
              if err != nil {
                return "", err
              }
              defer file.Close()
              cmd.Stdout = file
            } else if redir.Fd == 2 {
              file, err := os.OpenFile(redir.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
              if err != nil {
                return "", err
              }
              defer file.Close()
              cmd.Stderr = file
            }
          }
        }
        if cmd.Stdout == nil {
          cmd.Stdout = os.Stdout
        }
        if cmd.Stderr == nil {
          cmd.Stderr = os.Stderr
        }
        _ = cmd.Run()
        command.Redirs = []Redirection{}
        return "", nil
      }
      return "", fmt.Errorf("%s: command not found\n",command.Name)
    }

    return "", nil
}

func shell() error {
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
    tokens, err := NewLexer(line).Lex()
    if err != nil {
      return err
    }

    command, err := ParseTokens(tokens)
    if err != nil {
      return err
    }

    result, err := eval(command)
    stdout := os.Stdout
    stderr := os.Stderr
    for _, redir := range command.Redirs {
      if redir.Type == ">" && redir.Fd == 1 {
        file, err := os.Create(redir.FilePath)
        if err != nil {
          return err
        }
        defer file.Close()
        stdout = file
      } else if redir.Type == ">" && redir.Fd == 2 {
        file, err := os.Create(redir.FilePath)
        if err != nil {
          return err
        }
        defer file.Close()
        stderr = file
      } else if redir.Type == ">>" && redir.Fd == 1 {
        file, err := os.OpenFile(redir.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        if err != nil {
          return err
        }
        defer file.Close()
        stdout = file
      } else if redir.Type == ">>" && redir.Fd == 2 {
        file, err := os.OpenFile(redir.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        if err != nil {
          return err
        }
        defer file.Close()
        stderr = file
      }
    }
    
    if err != nil {
      fmt.Fprint(stderr, err)
    } else {
      fmt.Fprint(stdout, result)
    }
  }
}

func main() {
  path := os.Getenv("PATH")
  pathDirs = append(pathDirs, strings.Split(path, ":")...)

  if err := shell(); err != nil {
    fmt.Fprint(os.Stderr, err)
  }

}

