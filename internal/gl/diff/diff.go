package diff

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

func Summary(cfg *config.Config, project string, mr int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	mrData, _, err := client.MergeRequests.GetMergeRequest(projectID, int64(mr), &glclient.GetMergeRequestsOptions{})
	if err != nil {
		return fmt.Errorf("failed to get merge request: %w", err)
	}

	fmt.Printf("=== Merge Request !%d ===\n\n", mrData.IID)
	fmt.Printf("Title: %s\n", mrData.Title)
	fmt.Printf("State: %s\n", mrData.State)
	fmt.Printf("Author: %s\n", mrData.Author.Name)

	fmt.Printf("\nBranches:\n")
	fmt.Printf("  Source: %s\n", mrData.SourceBranch)
	fmt.Printf("  Target: %s\n", mrData.TargetBranch)

	fmt.Printf("\nMetrics:\n")
	fmt.Printf("  Upvotes: %d ⬆\n", mrData.Upvotes)
	fmt.Printf("  Downvotes: %d ⬇\n", mrData.Downvotes)

	if mrData.MergedAt != nil {
		fmt.Printf("  Merged: %s\n", mrData.MergedAt.Format("2006-01-02 15:04:05"))
	}

	if mrData.ClosedAt != nil {
		fmt.Printf("  Closed: %s\n", mrData.ClosedAt.Format("2006-01-02 15:04:05"))
	}

	if mrData.Description != "" {
		fmt.Printf("\nDescription:\n")
		desc := mrData.Description
		if len(desc) > 200 {
			desc = desc[:200] + "..."
		}
		fmt.Printf("  %s\n", desc)
	}

	fmt.Printf("\n═══════════════════════════════════════\n")
	fmt.Printf("View full diff and compare at:\n  %s\n", mrData.WebURL)
	fmt.Printf("═══════════════════════════════════════\n")
	return nil
}
