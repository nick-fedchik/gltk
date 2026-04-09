package search

import (
	"fmt"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func newClient(cfg *config.Config) (*glclient.Client, error) {
	token := cfg.Token
	client, err := glclient.NewClient(token, glclient.WithBaseURL(cfg.GitLabURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}
	return client, nil
}

func getProjectID(project string) interface{} {
	var id int
	if _, err := fmt.Sscanf(project, "%d", &id); err == nil {
		return id
	}
	return project
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func Issues(cfg *config.Config, project, search, state string, page, limit int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.ListProjectIssuesOptions{
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: int64(limit),
		},
	}

	if state != "all" {
		opts.State = glclient.Ptr(state)
	}

	if search != "" {
		opts.Search = glclient.Ptr(search)
	}

	issues, _, err := client.Issues.ListProjectIssues(projectID, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(issues) == 0 {
		fmt.Println("No issues found")
		return nil
	}

	fmt.Printf("Issue Search Results (%d found):\n\n", len(issues))

	for _, issue := range issues {
		status := "✓"
		if issue.State == "closed" {
			status = "✗"
		}

		fmt.Printf("%s #%d  %s\n",
			status,
			issue.IID,
			truncate(issue.Title, 70),
		)
		fmt.Printf("   State: %s | Author: %s\n", issue.State, issue.Author.Name)

		if issue.Description != "" {
			desc := issue.Description
			if len(desc) > 80 {
				desc = desc[:80] + "..."
			}
			fmt.Printf("   %s\n", desc)
		}
		fmt.Printf("\n")
	}
	return nil
}

func MRs(cfg *config.Config, project, search, state string, page, limit int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.ListProjectMergeRequestsOptions{
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: int64(limit),
		},
	}

	if state != "all" {
		opts.State = glclient.Ptr(state)
	}

	if search != "" {
		opts.Search = glclient.Ptr(search)
	}

	mrs, _, err := client.MergeRequests.ListProjectMergeRequests(projectID, opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(mrs) == 0 {
		fmt.Println("No merge requests found")
		return nil
	}

	fmt.Printf("Merge Request Search Results (%d found):\n\n", len(mrs))

	for _, mr := range mrs {
		status := "→"
		if mr.State == "merged" {
			status = "✓"
		} else if mr.State == "closed" {
			status = "✗"
		}

		fmt.Printf("%s !%d  %s\n",
			status,
			mr.IID,
			truncate(mr.Title, 70),
		)
		fmt.Printf("   State: %s | Author: %s\n", mr.State, mr.Author.Name)
		fmt.Printf("   Branch: %s → %s\n", mr.SourceBranch, mr.TargetBranch)

		if mr.Description != "" {
			desc := mr.Description
			if len(desc) > 80 {
				desc = desc[:80] + "..."
			}
			fmt.Printf("   %s\n", desc)
		}
		fmt.Printf("\n")
	}
	return nil
}
