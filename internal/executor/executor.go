package executor

import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "io"

  "github.com/cheesyhypocrisy/harsh/internal/parser"
)

var PathDirs []string
var Hist []string

type builtin int

const (
  unknownBuiltin builtin = iota
  exit
  echo
  _type
  pwd
  cd
  history
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
  case "history":
    return history
  default:
    return unknownBuiltin
  }
}

func findExecutable(command string) (string, error) {
  path := ""
  for _, dir:= range PathDirs {
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

type Runnable struct {
  isBuiltin bool
  Start func(stdin io.Reader, stdout, stderr io.Writer)
  Wait func()
}

func WrapBuiltin(command *parser.Command) Runnable {
  return Runnable {
    isBuiltin: true,
    Start: func(stdin io.Reader, stdout, stderr io.Writer) {
      switch lookupBuiltin(command.Name) {
      case exit:
        code := 0
        err := error(nil)
        if len(command.Args) >= 0 {
          code, err = strconv.Atoi(command.Args[0])
          if err != nil {
            fmt.Fprintln(stderr, err)
            return
          }
        }
        os.Exit(code)
      case echo:
        output := strings.Join(command.Args, "")
        fmt.Fprintln(stdout, output)
      case _type:
        if len(command.Args) == 0 {
          fmt.Fprintln(stderr, "Missing argument for type command")
          return
        }
        if lookupBuiltin(command.Args[0]) == unknownBuiltin {
          path, err := findExecutable(command.Args[0])
          if err != nil {
            fmt.Fprintf(stderr, "%s: not found\n", command.Args[0])
          } else {
            fmt.Fprintf(stdout, "%s is %s\n", command.Args[0], path)
          }
        } else {
          fmt.Fprintf(stdout, "%s is a shell builtin\n", command.Args[0])
        }
      case pwd:
        dir, err := os.Getwd()
        if err != nil {
          fmt.Fprintln(stderr, err)
          return
        }
        fmt.Fprintln(stdout, dir)
      case cd:
        if len(command.Args) == 0 || command.Args[0] == "~" {
          homeDir, exists := os.LookupEnv("HOME")
          if !exists {
            username := os.Getenv("USER")
            homeDir = fmt.Sprintf("/home/%s", username)
          }

          if err := os.Chdir(homeDir); err != nil {
            fmt.Fprintf(stderr, "cd: %s: No such file or directory\n", homeDir)
          }

          return
        }
        if err := os.Chdir(command.Args[0]); err != nil {
          fmt.Fprintf(stderr, "cd: %s: No such file or directory\n", command.Args[0]) 
        }
      case history:
        limit := len(Hist)
        err := error(nil)
        if len(command.Args) != 0 {
          limit, err = strconv.Atoi(command.Args[0])
          if err != nil {
            fmt.Fprintf(stderr, "%s", err.Error())
            return
          }
        }
        for i := max(0, len(Hist)-limit) ; i < len(Hist); i++ {
          fmt.Fprintf(stdout, "%d %s\n", i+1, Hist[i])
        }
      }
    },
    Wait: func() {
      return
    },
  }
}

func WrapExternal(command *parser.Command) Runnable {
  var cmd *exec.Cmd
  return Runnable {
    Start: func(stdin io.Reader, stdout, stderr io.Writer) {
      cmd = exec.Command(command.Name, command.Args...)
      cmd.Stdin = stdin
      cmd.Stdout = stdout
      cmd.Stderr = stderr

      cmd.Start()
    },
    Wait: func() {
      if cmd == nil {
        return
      }
      cmd.Wait()
    },
  }
}

func Eval(commands []*parser.Command) {
  runnables := make([]Runnable, 0)
  for _, command := range commands {
    if lookupBuiltin(command.Name) != unknownBuiltin {
      runnables = append(runnables, WrapBuiltin(command))
    } else{
      _, err := findExecutable(command.Name)
      if err == nil {
        runnables = append(runnables, WrapExternal(command))
      } else {
        fmt.Fprintf(os.Stderr, "%s: command not found\n",command.Name)
      }
    }
  }

  pipes := make([]*os.File, 0, 2*len(runnables))

  for i := 0; i < len(runnables); i++ {
    var stdin io.Reader = os.Stdin
    var stdout io.Writer = os.Stdout
    var stderr io.Writer = os.Stderr

    if i > 0 {
      stdin = pipes[2*(i-1)]
    }

    if i < len(runnables)-1 {
      r, w, _ := os.Pipe()
      pipes = append(pipes, r, w)
      stdout = w
    }

    if i == len(runnables)-1 {
      for _, redir := range commands[i].Redirs {
        if redir.Type == ">" {
          if redir.Fd == 1 {
            file, err := os.Create(redir.FilePath)
            if err != nil {
              fmt.Fprint(os.Stderr, err)
            }
            defer file.Close()
            stdout = file
          } else if redir.Fd == 2 {
            file, err := os.Create(redir.FilePath)
            if err != nil {
              fmt.Fprint(os.Stderr, err)
            }
            defer file.Close()
            stderr = file
          }
        } else if redir.Type == ">>" {
          if redir.Fd == 1 {
            file, err := os.OpenFile(redir.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
            if err != nil {
              fmt.Fprint(os.Stderr, err)
            }
            defer file.Close()
            stdout = file
          } else if redir.Fd == 2 {
            file, err := os.OpenFile(redir.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
            if err != nil {
              fmt.Fprint(os.Stderr, err)
            }
            defer file.Close()
            stderr = file
          }
        }
      }
      commands[i].Redirs = []parser.Redirection{}
    }

    runnables[i].Start(stdin, stdout, stderr)
  }

  for _, pipe := range pipes {
    pipe.Close()
  }

  for _, r := range runnables {
    r.Wait()
  }
}

