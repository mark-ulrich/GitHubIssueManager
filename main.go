// A tool to create, read, update, and delete GitHub issues from the command
// line. Built for exercise 4.11 from The Go Programming Language.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Command int

const (
	CommandInvalid Command = iota
	CommandQuit
	CommandCreate
	CommandRead
	CommandUpdate
	CommandDelete
	CommandList
)

var (
	appBaseName string
)

func main() {

	appBaseName = filepath.Base(os.Args[0])
	checkParms()

	repo := Repository{Name: os.Args[1]}
	_, err := fetchIssues(&repo)
	if err != nil {
		log.Fatal(err.Error())
	}

	isQuitting := false
	for !isQuitting {
		cmd := showMainMenu(&repo)
		switch cmd {
		case CommandList:
			listIssues(&repo)
		case CommandRead:
			readIssue(&repo)
		case CommandQuit:
			fmt.Printf("Quitting...\n")
			isQuitting = true
		case CommandInvalid:
			fmt.Printf("[!] Invalid Command\n ")
		}
	}

}

// Verify the user has supplied the necessary command-line arguments.
func checkParms() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <repo>\n", appBaseName)
		os.Exit(1)
	}
}
