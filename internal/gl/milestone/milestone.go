package milestone

import (
	"fmt"
	"time"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func List(cfg *config.Config, groupID int, state string, page, perPage int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	opts := &glclient.ListGroupMilestonesOptions{
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}

	if state != "" {
		opts.State = &state
	}

	milestones, resp, err := client.GroupMilestones.ListGroupMilestones(groupID, opts)
	if err != nil {
		return fmt.Errorf("failed to list milestones: %w", err)
	}

	if len(milestones) == 0 {
		fmt.Println("No milestones found.")
		return nil
	}

	fmt.Println("Milestones:")
	for _, m := range milestones {
		dueStr := ""
		if m.DueDate != nil {
			dueStr = fmt.Sprintf(" → %s", time.Time(*m.DueDate).Format("2006-01-02"))
		}
		fmt.Printf("  #%d: %s (%s)%s\n", m.IID, m.Title, m.State, dueStr)
	}

	if resp.NextPage > 0 {
		fmt.Printf("\nPage %d/%d | Next: --page=%d\n", page, resp.TotalPages, resp.NextPage)
	}
	return nil
}

func Create(cfg *config.Config, groupID int, title, description, startDate, dueDate string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	opts := &glclient.CreateGroupMilestoneOptions{
		Title: &title,
	}

	if description != "" {
		opts.Description = &description
	}

	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return fmt.Errorf("invalid start date format (use YYYY-MM-DD): %w", err)
		}
		iso := glclient.ISOTime(t)
		opts.StartDate = &iso
	}

	if dueDate != "" {
		t, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format (use YYYY-MM-DD): %w", err)
		}
		iso := glclient.ISOTime(t)
		opts.DueDate = &iso
	}

	m, _, err := client.GroupMilestones.CreateGroupMilestone(groupID, opts)
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}

	fmt.Printf("✅ Created milestone [%d] #%d: %s\n", m.ID, m.IID, m.Title)
	return nil
}

func Update(cfg *config.Config, groupID, milestoneID int, title, description, startDate, dueDate, state string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	opts := &glclient.UpdateGroupMilestoneOptions{}

	if title != "" {
		opts.Title = &title
	}
	if description != "" {
		opts.Description = &description
	}
	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return fmt.Errorf("invalid start date format (use YYYY-MM-DD): %w", err)
		}
		iso := glclient.ISOTime(t)
		opts.StartDate = &iso
	}
	if dueDate != "" {
		t, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format (use YYYY-MM-DD): %w", err)
		}
		iso := glclient.ISOTime(t)
		opts.DueDate = &iso
	}
	if state != "" {
		opts.StateEvent = &state
	}

	m, _, err := client.GroupMilestones.UpdateGroupMilestone(groupID, int64(milestoneID), opts)
	if err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	fmt.Printf("✓ Milestone updated\n")
	fmt.Printf("  #%d: %s\n", m.IID, m.Title)
	if m.DueDate != nil {
		fmt.Printf("  Due: %s\n", time.Time(*m.DueDate).Format("2006-01-02"))
	}
	return nil
}

func Delete(cfg *config.Config, groupID, milestoneID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	_, err = client.GroupMilestones.DeleteGroupMilestone(groupID, int64(milestoneID))
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	fmt.Printf("✓ Milestone #%d deleted\n", milestoneID)
	return nil
}
