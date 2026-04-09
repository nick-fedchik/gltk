package project

import (
	"encoding/json"
	"fmt"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func List(cfg *config.Config, groupID int) error {
	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	projects, _, err := client.Groups.ListGroupProjects(groupID, &glclient.ListGroupProjectsOptions{
		ListOptions: glclient.ListOptions{PerPage: 100},
	})
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	fmt.Println("Projects:")
	for _, p := range projects {
		fmt.Printf("  #%d: %s (%s) → %s\n", p.ID, p.Name, p.PathWithNamespace, p.WebURL)
	}
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
