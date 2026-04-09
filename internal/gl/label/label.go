package label

import (
	"fmt"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func List(cfg *config.Config, projectID, groupID, page, perPage int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	if groupID != 0 {
		opts := &glclient.ListGroupLabelsOptions{
			ListOptions: glclient.ListOptions{
				Page:    int64(page),
				PerPage: int64(perPage),
			},
		}
		labels, _, err := client.GroupLabels.ListGroupLabels(groupID, opts)
		if err != nil {
			return fmt.Errorf("failed to list group labels: %w", err)
		}
		if len(labels) == 0 {
			fmt.Println("No labels found")
			return nil
		}
		fmt.Printf("Group Labels (%d):\n\n", len(labels))
		for _, l := range labels {
			desc := ""
			if l.Description != "" {
				desc = "  — " + l.Description
			}
			fmt.Printf("  %-30s  %s%s\n", l.Name, l.Color, desc)
		}
		return nil
	}

	opts := &glclient.ListLabelsOptions{
		ListOptions: glclient.ListOptions{
			Page:    int64(page),
			PerPage: int64(perPage),
		},
	}
	labels, _, err := client.Labels.ListLabels(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list project labels: %w", err)
	}
	if len(labels) == 0 {
		fmt.Println("No labels found")
		return nil
	}
	fmt.Printf("Project Labels (%d):\n\n", len(labels))
	for _, l := range labels {
		desc := ""
		if l.Description != "" {
			desc = "  — " + l.Description
		}
		fmt.Printf("  %-30s  %s  (open issues: %d)%s\n", l.Name, l.Color, l.OpenIssuesCount, desc)
	}
	return nil
}

func Create(cfg *config.Config, projectID, groupID int, name, color, desc string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	if groupID != 0 {
		opts := &glclient.CreateGroupLabelOptions{
			Name:  &name,
			Color: &color,
		}
		if desc != "" {
			opts.Description = &desc
		}
		label, _, err := client.GroupLabels.CreateGroupLabel(groupID, opts)
		if err != nil {
			return fmt.Errorf("failed to create group label: %w", err)
		}
		fmt.Printf("✅ Group label created\n")
		fmt.Printf("  Name:  %s\n", label.Name)
		fmt.Printf("  Color: %s\n", label.Color)
		fmt.Printf("  ID:    %d\n", label.ID)
		return nil
	}

	opts := &glclient.CreateLabelOptions{
		Name:  &name,
		Color: &color,
	}
	if desc != "" {
		opts.Description = &desc
	}
	label, _, err := client.Labels.CreateLabel(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create label: %w", err)
	}
	fmt.Printf("✅ Label created\n")
	fmt.Printf("  Name:  %s\n", label.Name)
	fmt.Printf("  Color: %s\n", label.Color)
	fmt.Printf("  ID:    %d\n", label.ID)
	return nil
}
