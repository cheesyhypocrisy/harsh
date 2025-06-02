package main

import (
	"bufio"
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

func eval(command string, args []string) (string, error) {
    switch lookupBuiltin(command) {
    case exit:
      if len(args) == 0 {
        os.Exit(0)
      }
      exitCode, err := strconv.Atoi(args[0])
      if err != nil {
        return "", err
      }
      os.Exit(exitCode)
    case echo:
      return strings.Join(args, " ")+"\n", nil
    case _type:
      if len(args) == 0 {
        return "", fmt.Errorf("Missing argument for type command\n")
      }
      if lookupBuiltin(args[0]) == unknownBuiltin {
        path, err := findExecutable(args[0])
        if err != nil {
          return "", fmt.Errorf("%s: not found\n", args[0])
        } else {
          return fmt.Sprintf("%s is %s\n", args[0], path), nil
        }
      } else {
        return fmt.Sprintf("%s is a shell builtin\n", args[0]), nil
      }
    case pwd:
      dir, err := os.Getwd()
      return dir + "\n", err
    case cd:
      if len(args) == 0 || args[0] == "~" {
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
      if err := os.Chdir(args[0]); err != nil {
        return "", fmt.Errorf("cd: %s: No such file or directory\n", args[0]) 
      }
      return "", nil
    default:
      _, err := findExecutable(command)
      if err == nil {
        result, err := exec.Command(command, args...).Output()
        return string(result), err
      }
      return "", fmt.Errorf("%s: command not found\n",strings.TrimSpace(command))
    }

    return "", nil
}

func shell() error {
	for {
    fmt.Fprint(os.Stdout, "$ ")
  
    fullCommand, err := bufio.NewReader(os.Stdin).ReadString('\n')
    if err != nil {
      return err
    }
    
    fullCommand = strings.TrimSpace(fullCommand)
    cliArgs := strings.Split(fullCommand, " ")
    command := cliArgs[0]
    args := cliArgs[1:]

    result, err := eval(command, args)
    
    if err != nil {
      fmt.Fprint(os.Stderr, err)
    } else {
      fmt.Print(result)
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

