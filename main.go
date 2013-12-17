package main

import (
    "os/exec"   
    "os"
    "fmt"
    "io"
    "log"
    "github.com/howeyc/fsnotify"
    "time"
)

func main() {
    runTest()
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
                log.Println("event:", ev)

                // runTest()
                newTime := time.Now().Unix()
                diffTime := newTime - timestamp
                if diffTime > 4 {
                    go runTest()
                }
                timestamp = newTime
            case err := <-watcher.Error:
                log.Println("error:", err)
            }
        }
    }()

    err = watcher.Watch("./")
    if err != nil {
        log.Fatal(err)
    }

    <-done

    /* ... do stuff ... */
    watcher.Close()
}

func runTest() {
    time.Sleep(100 * time.Millisecond)
    fmt.Println()
    fmt.Println()
    fmt.Println()
    fmt.Println("====== bugging ======")
    fmt.Println()
    fmt.Println()
    fmt.Println()
    cmd := exec.Command("go", "test")
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
