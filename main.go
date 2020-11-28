// A tool to create, read, update, and delete GitHub issues from the command
// line. Built for exercise 4.11 from The Go Programming Language.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

type Repository struct {
	Name           string
	TotalIssues    int
	OpenIssueCount int
	Issues         *[]Issue
}

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []Issue
}

type Issue struct {
	Title     string
	User      User
	Labels    []Label
	State     string
	Assignees []User
	Comments  int
	CreatedAt string `json:"created_at"`
	Body      string
}

type User struct {
	Login string
}

type Label struct {
	Name string
}

const (
	githubIssueSearchBaseUrl = "https://api.github.com/search/issues?q="
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

// Reads an integer from standard input. On error, the returned int will be 0,
// and the error will be returned the caller along with the actual read string.
func readInt() (int, error, string) {
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		return 0, err, in
	}
	in = strings.ToLower(strings.TrimSpace(in))
	if err != nil {
		return 0, err, in
	}
	num, err := strconv.Atoi(in)
	if err != nil {
		return num, err, in
	}
	return num, nil, in
}

// Displays the given promp and calls readInt() to read an integer from standard
// input. On error, the returned int is 0 and the error will be returned the
// caller along with the actual read string.
func promptInt(prompt string) (int, error, string) {
	fmt.Print(prompt)
	return readInt()
}

// fetchIssues retrieves a list of issues for a named GitHub repository
func fetchIssues(repo *Repository) (*IssuesSearchResult, error) {

	// Perform search
	url := githubIssueSearchBaseUrl + "repo:" + repo.Name
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 422 {
		resp.Body.Close()
		return nil, fmt.Errorf("Unable to find repository: %s (Do you have permission to access this repository?)", repo.Name)
	} else if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetch failed: %s", resp.Status)
	}

	// Search was successful; unmarshal data
	var results IssuesSearchResult
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}

	// Fill out repo struct
	repo.Issues = &results.Items
	repo.TotalIssues = results.TotalCount
	for _, issue := range *repo.Issues {
		if issue.State == "open" {
			repo.OpenIssueCount++
		}
	}

	return &results, nil
}

// showMainMenu displays the programs main menu and reads an input string. It
// returns a Command based on the input.
func showMainMenu(repo *Repository) Command {
	menuOptions := []string{
		"List Issues",
		"Read Issue",
		"Create Issue",
		"Update Issue",
		"Delete Issue",
		"Quit",
	}

	fmt.Printf("\n[ Repo: %s ] Issues: %d (%d open / %d closed)\n", repo.Name, repo.TotalIssues, repo.OpenIssueCount, repo.TotalIssues-repo.OpenIssueCount)
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
	}

	return CommandInvalid
}

// List all issues in the repository.
func listIssues(repo *Repository) {
	issues := repo.Issues
	fmt.Println()
	for i, issue := range *issues {
		const MaxTitleLength = 40
		dateString := strings.Split(issue.CreatedAt, "T")[0]
		title := issue.Title
		if len(title) > MaxTitleLength {
			title = title[:MaxTitleLength-3] + "..."
		}
		stateString := fmt.Sprintf("[%s]", strings.Title(issue.State))
		numberString := fmt.Sprintf("[%d]", i+1)
		fmt.Printf("%-8s%-60s %-10s %s\n", numberString, title, stateString, dateString)
	}
}

// Prompt user for issue number and display the issue.
func readIssue(repo *Repository) {
	issues := repo.Issues
	var issueNumber int
	var err error
	for {
		issueNumber, err, _ = promptInt("\n  Enter an issue number: ")
		if err != nil {

		}
		if issueNumber < 1 || issueNumber > repo.TotalIssues {
			fmt.Printf("[!] Invalid isssue number: %d. Valid issues are 1-%d\n", issueNumber, repo.TotalIssues)
			continue
		}
		break // We have a valid issue number
	}
	issueIndex := issueNumber - 1
	issue := &(*issues)[issueIndex]

	dateString := strings.Split(issue.CreatedAt, "T")[0]
	fmt.Printf("\nTitle:  %s\nAuthor: %s\nDate:   %s\nState:  %s\n\n%s\n", issue.Title, issue.User.Login, dateString, strings.Title(issue.State), issue.Body)
}
