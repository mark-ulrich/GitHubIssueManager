package main

import (
	"fmt"
	"strings"
)

// showMainMenu displays the programs main menu and reads an input string. It
// returns a Command based on the input.
func showMainMenu(repo *Repository) Command {
	menuOptions := []string{
		"List Issues",
		"Read Issue",
		"Create Issue",
		"Update Issue",
		"Delete Issue",
		"Refresh",
		"Quit",
	}

	fmt.Printf("[ Repo: %s ] Issues: %d (%d open / %d closed)\n", repo.Name, repo.TotalIssues, repo.OpenIssueCount, repo.TotalIssues-repo.OpenIssueCount)
	for i, str := range menuOptions {
		fmt.Printf("  (%d)  %s\n", i+1, str)
	}
	selected, err, in := promptInt("  > ")
	if err != nil {
		in = strings.ToLower(strings.TrimSpace(in))
		if in == "q" || in == "quit" {
			return CommandQuit
		}
		return CommandInvalid
	}

	if selected < 1 || selected > len(menuOptions) {
		return CommandInvalid
	}

	selected--
	switch strings.ToLower(strings.Split(menuOptions[selected], " ")[0]) {
	case "read":
		return CommandRead
	case "update":
		return CommandUpdate
	case "delete":
		return CommandDelete
	case "create":
		return CommandCreate
	case "quit":
		return CommandQuit
	case "list":
		return CommandList
	case "refresh":
		return CommandRefresh
	}

	return CommandInvalid
}
