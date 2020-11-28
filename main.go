// A tool to create, read, update, and delete GitHub issues from the command
// line. Built for exercise 4.11 from The Go Programming Language.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	CommandRefresh
)

var (
	appBaseName string
	settings    Settings
)

var (
	token  = flag.String("token", "", "GitHub Personal Access Token")
	editor = flag.String("editor", "", "Preferred external text editor")
)

func main() {

	var err error

	flag.Parse()

	appBaseName = filepath.Base(os.Args[0])
	checkParms()

	repo := Repository{
		Name:    flag.CommandLine.Args()[0],
		Token:   *token,
		IsDirty: true,
	}

	loadSettings(&settings, defaultSettingsFilename)
	if len(*editor) > 0 {
		settings.EditorCommand = *editor
	}
	if settings.EditorCommand == "" {
		switch runtime.GOOS {
		case "windows":
			settings.EditorCommand = "notepad.exe"
		case "linux":
			settings.EditorCommand = "vi"
		}
	}

	isQuitting := false
	for !isQuitting {
		if repo.IsDirty {
			_, err = fetchIssues(&repo)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		cmd := showMainMenu(&repo)
		switch cmd {
		case CommandList:
			listIssues(&repo)
		case CommandRead:
			readIssue(&repo)
		case CommandCreate:
			err = createIssue(&repo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] Error creating issue: %s\n\n", err.Error())
			}
		case CommandRefresh:
			repo.IsDirty = true
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
