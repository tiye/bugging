
Bugging: go debugging
------

This command just runs `go test` at working directory on file changes.

Code is mainly glued togather from these two snippets:

* https://github.com/howeyc/fsnotify
* http://stackoverflow.com/q/8875038/883571

### Usage

Installation:

```
go get github.com/jiyinyiyong/bugging
```

Watch files in current directory and repeat command on events:

```
bugging go test
```

This command will try running `go test` on every file change.

### Bugs

* files in deeper level are not watched
* quite many file changes when saving with Sublime Text

### License

MIT