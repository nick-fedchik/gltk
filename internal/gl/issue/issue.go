package issue

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func List(cfg *config.Config, projectID int, state string, page, perPage int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	opts := &glclient.ListProjectIssuesOptions{
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}
	if state != "" && state != "all" {
		opts.State = &state
	}

	issues, resp, err := client.Issues.ListProjectIssues(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	for _, issue := range issues {
		fmt.Printf("#%d: %s (%s)\n", issue.IID, issue.Title, issue.State)
	}

	if resp.NextPage > 0 {
		fmt.Printf("\nPage %d/%d | Next: %d\n", page, resp.TotalPages, resp.NextPage)
	}
	return nil
}

func Create(cfg *config.Config, projectID int, title, description, labels string, milestone int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	opts := &glclient.CreateIssueOptions{
		Title: &title,
	}

	if description != "" {
		opts.Description = &description
	}

	if labels != "" {
		labelList := glclient.LabelOptions(strings.Split(labels, ","))
		opts.Labels = &labelList
	}

	if milestone > 0 {
		m := int64(milestone)
		opts.MilestoneID = &m
	}

	issue, _, err := client.Issues.CreateIssue(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	fmt.Printf("✅ Created issue #%d: %s\n", issue.IID, issue.Title)
	return nil
}

func Close(cfg *config.Config, projectID, iid int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	state := "closed"
	opts := &glclient.UpdateIssueOptions{
		StateEvent: &state,
	}

	issue, _, err := client.Issues.UpdateIssue(projectID, int64(iid), opts)
	if err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}

	fmt.Printf("✅ Closed issue #%d: %s\n", issue.IID, issue.Title)
	return nil
}

func Batch(cfg *config.Config, projectID int, filePath string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var issues []struct {
		Title       string   `json:"title"`
		Description string   `json:"description,omitempty"`
		Labels      []string `json:"labels,omitempty"`
		MilestoneID int      `json:"milestone_id,omitempty"`
	}

	if err := json.Unmarshal(data, &issues); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Creating %d issues...\n", len(issues))

	for i, issue := range issues {
		title := issue.Title
		opts := &glclient.CreateIssueOptions{
			Title: &title,
		}

		if issue.Description != "" {
			desc := issue.Description
			opts.Description = &desc
		}

		if len(issue.Labels) > 0 {
			labels := glclient.LabelOptions(issue.Labels)
			opts.Labels = &labels
		}

		if issue.MilestoneID > 0 {
			mid := int64(issue.MilestoneID)
			opts.MilestoneID = &mid
		}

		created, _, err := client.Issues.CreateIssue(projectID, opts)
		if err != nil {
			fmt.Printf("❌ [%d/%d] Failed: %s - %v\n", i+1, len(issues), issue.Title, err)
			continue
		}

		fmt.Printf("✅ [%d/%d] Created #%d: %s\n", i+1, len(issues), created.IID, created.Title)
	}

	fmt.Println("✨ Batch creation complete!")
	return nil
}
