package main

import (
    "os/exec"   
    "os"
    "fmt"
    "io"
    "log"
    "github.com/howeyc/fsnotify"
    "time"
    "path/filepath"
)

func main() {
    if len(os.Args) < 2 {
        log.Println("need more args")
        os.Exit(1)
    }
    go runCommand()
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }

    done := make(chan bool)
    timestamp := int64(0)

    // Process events
    go func() {
        for {
            select {
            case ev := <- watcher.Event:
                if ev.IsModify() {
                    if filepath.Ext(ev.Name) == ".go" {
                        log.Println("event:", ev, ev.Name)

                        newTime := time.Now().Unix()
                        diffTime := newTime - timestamp
                        if diffTime > 3 {
                            go runCommand()
                        }
                        timestamp = newTime
                    }
                }
            case err := <- watcher.Error:
                log.Println("error:", err)
            }
        }
    }()

    err = watcher.Watch("./")
    if err != nil {
        log.Fatal(err)
    }

    <- done

    /* ... do stuff ... */
    watcher.Close()
}

func runCommand() {
    time.Sleep(100 * time.Millisecond)
    fmt.Println()
    fmt.Println()
    fmt.Println("\x1b[33;1m~@~@~@~@~@~@~ Hello, World! ~@~@~@~@~@~@~\x1b[0m")
    fmt.Println()

    cmd := exec.Command(os.Args[1], os.Args[2:]...)
    log.Println("\x1b[35;1m", os.Args[0], os.Args[1], "\x1b[0m")
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
