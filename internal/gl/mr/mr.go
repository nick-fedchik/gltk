package mr

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

func List(cfg *config.Config, project, state string, page, perPage int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.ListProjectMergeRequestsOptions{
		State: glclient.Ptr(state),
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}

	mrs, resp, err := client.MergeRequests.ListProjectMergeRequests(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list merge requests: %w", err)
	}

	if len(mrs) == 0 {
		fmt.Printf("No merge requests found (state: %s)\n", state)
		return nil
	}

	fmt.Printf("Merge Requests (%d total):\n\n", resp.TotalItems)
	for _, mr := range mrs {
		status := "✓"
		if mr.State == "opened" {
			status = "→"
		} else if mr.State == "closed" {
			status = "✗"
		}

		fmt.Printf("%s !%d  %-40s  [%s]  %s → %s\n",
			status,
			mr.IID,
			truncate(mr.Title, 40),
			mr.State,
			mr.SourceBranch,
			mr.TargetBranch,
		)
	}
	fmt.Printf("\nPage %d/%d\n", page, resp.TotalPages)
	return nil
}

func Get(cfg *config.Config, project string, mrIID int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	mrData, _, err := client.MergeRequests.GetMergeRequest(projectID, int64(mrIID), &glclient.GetMergeRequestsOptions{})
	if err != nil {
		return fmt.Errorf("failed to get merge request: %w", err)
	}

	fmt.Printf("!%d — %s\n", mrData.IID, mrData.Title)
	fmt.Printf("State: %s\n", mrData.State)
	fmt.Printf("Author: %s\n", mrData.Author.Name)
	fmt.Printf("Branch: %s → %s\n", mrData.SourceBranch, mrData.TargetBranch)
	fmt.Printf("Upvotes: %d | Downvotes: %d\n", mrData.Upvotes, mrData.Downvotes)
	fmt.Printf("Created: %s\n", mrData.CreatedAt.Format("2006-01-02 15:04:05"))
	if mrData.Description != "" {
		fmt.Printf("\nDescription:\n%s\n", mrData.Description)
	}
	fmt.Printf("\nWeb URL: %s\n", mrData.WebURL)
	return nil
}

func Create(cfg *config.Config, project, title, source, target, description string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.CreateMergeRequestOptions{
		Title:        glclient.Ptr(title),
		SourceBranch: glclient.Ptr(source),
		TargetBranch: glclient.Ptr(target),
		Description:  glclient.Ptr(description),
	}

	mr, _, err := client.MergeRequests.CreateMergeRequest(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create merge request: %w", err)
	}

	fmt.Printf("✓ Merge request created\n")
	fmt.Printf("  IID: !%d\n", mr.IID)
	fmt.Printf("  Title: %s\n", mr.Title)
	fmt.Printf("  URL: %s\n", mr.WebURL)
	return nil
}

func Merge(cfg *config.Config, project string, mrIID int, message string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.AcceptMergeRequestOptions{}
	if message != "" {
		opts.MergeCommitMessage = glclient.Ptr(message)
	}

	mrData, _, err := client.MergeRequests.AcceptMergeRequest(projectID, int64(mrIID), opts)
	if err != nil {
		return fmt.Errorf("failed to merge request: %w", err)
	}

	fmt.Printf("✓ Merge request merged\n")
	fmt.Printf("  IID: !%d\n", mrData.IID)
	fmt.Printf("  State: %s\n", mrData.State)
	return nil
}

func Close(cfg *config.Config, project string, mrIID int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.UpdateMergeRequestOptions{
		StateEvent: glclient.Ptr("close"),
	}

	mrData, _, err := client.MergeRequests.UpdateMergeRequest(projectID, int64(mrIID), opts)
	if err != nil {
		return fmt.Errorf("failed to close merge request: %w", err)
	}

	fmt.Printf("✓ Merge request closed\n")
	fmt.Printf("  IID: !%d\n", mrData.IID)
	return nil
}

func Comment(cfg *config.Config, project string, mrIID int, body string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.CreateMergeRequestNoteOptions{
		Body: glclient.Ptr(body),
	}

	note, _, err := client.Notes.CreateMergeRequestNote(projectID, int64(mrIID), opts)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	fmt.Printf("✓ Comment added\n")
	fmt.Printf("  Note ID: %d\n", note.ID)
	fmt.Printf("  Author: %s\n", note.Author.Name)
	return nil
}
