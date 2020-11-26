// A tool to create, read, update, and delete GitHub issues from the command
// line. Built for exercise 4.11 from The Go Programming Language.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
	repoName := os.Args[1]
	_, err := fetchIssues(repoName)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkParms() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <repo>\n", appBaseName)
		os.Exit(1)
	}
}

// fetchIssues retrieves a list of issues for a named GitHub repository
func fetchIssues(repoName string) (*IssuesSearchResult, error) {

	// Perform search
	url := githubIssueSearchBaseUrl + "repo:" + repoName
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 422 {
		resp.Body.Close()
		return nil, fmt.Errorf("Unable to find repository: %s (Do you have permission to access this repository?)", repoName)
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

	return &results, nil
}
