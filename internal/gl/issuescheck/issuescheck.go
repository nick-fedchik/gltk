package issuescheck

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func Run(cfg *config.Config, projectID int, interactive, listOnly bool, closeIDs string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return fmt.Errorf("failed to create GitLab client: %w", err)
	}

	fmt.Printf("🔗 Connected to: %s\n", cfg.GitLabURL)
	fmt.Printf("📦 Project ID: %d\n\n", projectID)

	switch {
	case listOnly:
		return listIssues(client, projectID)
	case interactive:
		return interactiveMode(client, projectID)
	case closeIDs != "":
		return closeIssues(client, projectID, closeIDs)
	default:
		fmt.Fprintf(os.Stderr, `gl-issues-check - Review and manage GitLab issues

Usage:
  gl-issues-check [options]

Options:
  -project=ID       Project ID (default: 1)
  -list            List opened issues only
  -interactive     Interactive mode: review issues and choose to close/skip
  -close=IDS       Close specific issues (comma-separated, e.g., "42,99,123")
`)
	}
	return nil
}

func listIssues(client *glclient.Client, projectID int) error {
	fmt.Println("📋 Opened Issues:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	opts := &glclient.ListProjectIssuesOptions{
		State: stringPtr("opened"),
		ListOptions: glclient.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	page := int64(1)
	totalCount := 0

	for {
		opts.ListOptions.Page = page

		issues, resp, err := client.Issues.ListProjectIssues(projectID, opts)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {
			totalCount++
			fmt.Printf("  #%-4d | %-60s | %s\n",
				issue.IID,
				truncate(issue.Title, 60),
				issue.State,
			)
		}

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Total: %d opened issues\n\n", totalCount)
	return nil
}

func interactiveMode(client *glclient.Client, projectID int) error {
	fmt.Println("🎯 Interactive Mode")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	opts := &glclient.ListProjectIssuesOptions{
		State: stringPtr("opened"),
		ListOptions: glclient.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	reader := bufio.NewReader(os.Stdin)
	closedCount := 0
	skippedCount := 0
	page := int64(1)

	for {
		opts.ListOptions.Page = page

		issues, resp, err := client.Issues.ListProjectIssues(projectID, opts)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {
			fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("Issue #%d\n", issue.IID)
			fmt.Printf("Title: %s\n", issue.Title)
			fmt.Printf("State: %s\n", issue.State)
			if issue.Description != "" {
				fmt.Printf("Description: %s\n", truncate(issue.Description, 100))
			}
			fmt.Printf("Created: %s\n\n", issue.CreatedAt)

			fmt.Print("Action? [c]lose, [s]kip, [q]uit: ")
			input, _ := reader.ReadString('\n')
			action := strings.TrimSpace(input)

			switch action {
			case "c", "C":
				if closeIssue(client, projectID, int(issue.IID)) {
					closedCount++
					fmt.Println("✅ Issue closed")
				}
			case "s", "S":
				skippedCount++
				fmt.Println("⏭️  Skipped")
			case "q", "Q":
				fmt.Println("\n👋 Exiting interactive mode")
				goto summary
			default:
				fmt.Println("❓ Invalid option")
			}
		}

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

summary:
	fmt.Printf("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Summary:\n")
	fmt.Printf("  ✅ Closed:  %d\n", closedCount)
	fmt.Printf("  ⏭️  Skipped: %d\n", skippedCount)
	fmt.Println()
	return nil
}

func closeIssues(client *glclient.Client, projectID int, closeIDsStr string) error {
	ids := strings.Split(strings.TrimSpace(closeIDsStr), ",")

	fmt.Printf("🔄 Closing %d issue(s)...\n", len(ids))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	closedCount := 0
	failedCount := 0

	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		iid, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Printf("❌ Invalid issue ID: %s\n", idStr)
			failedCount++
			continue
		}

		if closeIssue(client, projectID, iid) {
			closedCount++
		} else {
			failedCount++
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Summary: %d closed, %d failed\n", closedCount, failedCount)
	fmt.Println()
	return nil
}

func closeIssue(client *glclient.Client, projectID, iid int) bool {
	state := "closed"
	opts := &glclient.UpdateIssueOptions{
		StateEvent: &state,
	}

	issue, _, err := client.Issues.UpdateIssue(projectID, int64(iid), opts)
	if err != nil {
		fmt.Printf("❌ Failed to close #%d: %v\n", iid, err)
		return false
	}

	fmt.Printf("✅ Closed #%d: %s\n", issue.IID, truncate(issue.Title, 50))
	return true
}

func stringPtr(s string) *string {
	return &s
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
