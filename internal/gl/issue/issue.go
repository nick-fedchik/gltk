package issue

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gltk/gltk/internal/config"
	glclient "gitlab.com/gitlab-org/api/client-go"
)

func List(cfg *config.Config, projectID interface{}, state string, page, perPage int) error {
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

func resolveUsernamesToIDs(cfg *config.Config, usernames []string) ([]int64, error) {
	if len(usernames) == 0 {
		return nil, nil
	}
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, u := range usernames {
		users, _, err := client.Users.ListUsers(&glclient.ListUsersOptions{Username: glclient.Ptr(u)})
		if err != nil {
			return nil, fmt.Errorf("failed to look up user %q: %w", u, err)
		}
		if len(users) == 0 {
			return nil, fmt.Errorf("user not found: %q", u)
		}
		ids = append(ids, users[0].ID)
	}
	return ids, nil
}

func Create(cfg *config.Config, projectID interface{}, title, description, labels string, milestone int, assigneeIDs []int64, assigneeUsernames []string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	resolvedIDs, err := resolveUsernamesToIDs(cfg, assigneeUsernames)
	if err != nil {
		return err
	}
	allAssigneeIDs := append(assigneeIDs, resolvedIDs...)

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

	if len(allAssigneeIDs) > 0 {
		opts.AssigneeIDs = &allAssigneeIDs
	}

	issue, _, err := client.Issues.CreateIssue(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	fmt.Printf("✅ Created issue #%d: %s\n", issue.IID, issue.Title)
	if len(issue.Assignees) > 0 {
		names := make([]string, len(issue.Assignees))
		for i, a := range issue.Assignees {
			names[i] = a.Username
		}
		fmt.Printf("   Assigned to: %s\n", strings.Join(names, ", "))
	}
	return nil
}

func Close(cfg *config.Config, projectID interface{}, iid int) error {
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

func Batch(cfg *config.Config, projectID interface{}, filePath string) error {
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
