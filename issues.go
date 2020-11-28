package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// fetchIssues retrieves a list of issues for a named GitHub repository
func fetchIssues(repo *Repository) (*IssuesSearchResult, error) {

	// Perform search
	const githubIssueSearchBaseUrl = "https://api.github.com/search/issues?q="
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
		fmt.Printf("%-8s%-40s %-10s %s\n", numberString, title, stateString, dateString)
	}
}

// Prompt user for issue number and display the issue.
func readIssue(repo *Repository) {
	issues := repo.Issues
	var (
		issueNumber int
		err         error
		in          string
	)

	for {
		issueNumber, err, in = promptInt("\n  Enter an issue number: ")
		if err != nil || issueNumber < 1 || issueNumber > repo.TotalIssues {
			fmt.Printf("[!] Invalid isssue number: %s. Valid issues are 1-%d\n", in, repo.TotalIssues)
			continue
		}
		break // We have a valid issue number
	}

	issueIndex := issueNumber - 1
	issue := &(*issues)[issueIndex]
	dateString := strings.Split(issue.CreatedAt, "T")[0]
	fmt.Printf("\nTitle:  %s\nAuthor: %s\nDate:   %s\nState:  %s\n\n%s\n", issue.Title, issue.User.Login, dateString, strings.Title(issue.State), issue.Body)
}
