package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: folder-watcher <folder-to-watch> <command> <...args> $filepath")
		os.Exit(0)
	}
	folderToWatch := os.Args[1]
	commands := os.Args[2:]
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				rest := append(commands[1:], event.Name)

				log.Printf("executing %s %s", commands[0], strings.Join(rest, " "))
				out, _ := exec.Command(commands[0], rest...).CombinedOutput()
				log.Printf("\n%s", out)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(folderToWatch)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
