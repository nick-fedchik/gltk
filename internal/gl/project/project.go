package project

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gltk/gltk/internal/config"
	glclient "gitlab.com/gitlab-org/api/client-go"
)

func List(cfg *config.Config, groupID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	var projects []*glclient.Project
	if groupID != 0 {
		projects, _, err = client.Groups.ListGroupProjects(groupID, &glclient.ListGroupProjectsOptions{
			ListOptions: glclient.ListOptions{PerPage: 100},
		})
	} else {
		projects, _, err = client.Projects.ListProjects(&glclient.ListProjectsOptions{
			ListOptions: glclient.ListOptions{PerPage: 100},
			Membership:  glclient.Ptr(true),
		})
	}
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	fmt.Println("Projects:")
	for _, p := range projects {
		fmt.Printf("  #%d: %s (%s) → %s\n", p.ID, p.Name, p.PathWithNamespace, p.WebURL)
	}
	return nil
}

func Members(cfg *config.Config, projectID string) error {
	if projectID == "" {
		projectID = cfg.ProjectID
	}
	projectID = strings.TrimRight(projectID, "/")
	if projectID == "" {
		return fmt.Errorf("project ID or path required: use --project or set default-project in config")
	}

	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	members, _, err := client.ProjectMembers.ListAllProjectMembers(projectID, &glclient.ListProjectMembersOptions{
		ListOptions: glclient.ListOptions{PerPage: 100},
	})
	if err != nil {
		return fmt.Errorf("failed to list members: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUSERNAME\tNAME\tACCESS LEVEL\tSTATE")
	fmt.Fprintln(w, "--\t--------\t----\t------------\t-----")
	for _, m := range members {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			m.ID, m.Username, m.Name,
			accessLevelName(int(m.AccessLevel)),
			m.State,
		)
	}
	w.Flush()
	fmt.Printf("\nTotal: %d members\n", len(members))
	return nil
}

func GetByID(cfg *config.Config, projectID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	p, _, err := client.Projects.GetProject(projectID, &glclient.GetProjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func accessLevelName(level int) string {
	switch level {
	case 10:
		return "Guest"
	case 20:
		return "Reporter"
	case 30:
		return "Developer"
	case 40:
		return "Maintainer"
	case 50:
		return "Owner"
	default:
		return fmt.Sprintf("%d", level)
	}
}

func ByPath(cfg *config.Config, path string) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	p, _, err := client.Projects.GetProject(path, &glclient.GetProjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get project by path: %w", err)
	}

	fmt.Printf("Project: %s (ID: %d)\n", p.Name, p.ID)
	fmt.Printf("Path: %s\n", p.PathWithNamespace)
	fmt.Printf("URL: %s\n", p.WebURL)
	fmt.Printf("Visibility: %s\n", p.Visibility)
	return nil
}
