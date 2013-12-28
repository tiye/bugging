package main

import (
    "os/exec"   
    "os"
    "fmt"
    "io"
    "github.com/howeyc/fsnotify"
    "time"
    "strings"
    "path/filepath"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("need more args")
        os.Exit(1)
    }
    go runCommand()
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        panic(err)
    }

    done := make(chan bool)
    timestamp := int64(0)

    // Process events
    go func() {
        for {
            select {
            case ev := <- watcher.Event:
                if ev.IsModify() {
                    extname := filepath.Ext(ev.Name)
                    if extname != ".tmp" {
                        fmt.Println("event:", ev, ev.Name)

                        newTime := time.Now().Unix()
                        diffTime := newTime - timestamp
                        if diffTime > 2 {
                            go runCommand()
                            timestamp = newTime
                        }
                    }
                }
            case err := <- watcher.Error:
                fmt.Println("error:", err)
            }
        }
    }()

    err = watcher.Watch("./")
    if err != nil {
        panic(err)
    }

    <- done

    /* ... do stuff ... */
    watcher.Close()
}

func runCommand() {
    time.Sleep(100 * time.Millisecond)
    fmt.Println()
    fmt.Print("\x1b[33;1m@@@@@@@@@@@@ ")
    fmt.Print(strings.Join(os.Args[1:], " "))
    fmt.Println("\x1b[0m")
    fmt.Println()

    cmd := exec.Command(os.Args[1], os.Args[2:]...)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println(err)
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        fmt.Println(err)
    }
    err = cmd.Start()
    if err != nil {
        fmt.Println(err)
    }
    go io.Copy(os.Stdout, stdout) 
    go io.Copy(os.Stderr, stderr) 
    cmd.Wait()
}
