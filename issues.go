package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const githubAPIBaseURL = "https://api.github.com/"

type Repository struct {
	Name           string
	TotalIssues    int
	OpenIssueCount int
	Issues         *[]Issue
	Token          string
	IsDirty        bool
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
	url := githubAPIBaseURL + "search/issues?q=repo:" + repo.Name
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
	repo.OpenIssueCount = 0
	for _, issue := range *repo.Issues {
		if issue.State == "open" {
			repo.OpenIssueCount++
		}
	}
	repo.IsDirty = false

	return &results, nil
}

// List all issues in the repository.
func listIssues(repo *Repository) {
	issues := repo.Issues
	fmt.Printf("\n#       Name                                     State      Date\n")
	fmt.Printf("------- ---------------------------------------- ---------- ----------\n")
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
	fmt.Println()
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

// Create a new issue and add it to the repository. Invoke a configurable
// preferred text editor to edit the issue.
func createIssue(repo *Repository) error {

	if repo.Token == "" {
		return fmt.Errorf("You must supply a Personal Access Token to create an issue")
	}

	fmt.Printf("  Enter title: ")
	reader := bufio.NewReader(os.Stdin)
	title, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	title = strings.TrimSpace(title)
	body, err := editWithExternalEditor(title)
	if err != nil {
		return err
	}

	requestBody, err := json.Marshal(map[string]string{
		"title": fmt.Sprintf("%s", title),
		"body":  body,
	})
	if err != nil {
		return err
	}
	url := githubAPIBaseURL + "repos/" + repo.Name + "/issues"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", repo.Token))
	req.Header.Add("Content-type", "application/json")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg := resp.Status
		switch resp.StatusCode {
		case 201:
			fmt.Printf("\n  Created issue: %s\n\n", title)
		case 404:
			msg = "Unauthorized"
		default:
			return fmt.Errorf("Failed to create issue: %s\n", msg)
		}
	}
	resp.Body.Close()

	// Force update of local repo data
	repo.IsDirty = true

	return nil
}

// Update an existing issue. Invoke a configurable preferred text editor
// to edit the issue.
func updateIssue(repo *Repository, id int) error {
	return nil
}

// Delete an issue with the specified index.
func deleteIssue(repo *Repository, index int) error {
	return nil
}
