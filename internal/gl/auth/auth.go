package auth

import (
	"fmt"
	"log"

	glclient "gitlab.com/gitlab-org/api/client-go"
	"github.com/gltk/gltk/internal/config"
)

func Test(cfg *config.Config, testCmd string) error {
	log.Printf("GitLab Instance: %s", cfg.GitLabURL)
	log.Printf("Test Command: %s", testCmd)

	switch testCmd {
	case "user":
		return testUser(cfg.GitLabURL, cfg.Token)
	case "groups":
		return testGroups(cfg.GitLabURL, cfg.Token)
	case "projects":
		return testProjects(cfg.GitLabURL, cfg.Token)
	case "health":
		return testHealth(cfg.GitLabURL)
	default:
		return fmt.Errorf("unknown test command: %s", testCmd)
	}
}

func testHealth(baseURL string) error {
	log.Println("\n=== Testing GitLab Health ===")
	client, err := glclient.NewClient("", glclient.WithBaseURL(baseURL))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	version, _, err := client.Version.GetVersion()
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	log.Printf("✓ GitLab is healthy")
	log.Printf("  Version: %s", version.Version)
	log.Printf("  Revision: %s", version.Revision)
	return nil
}

func testUser(baseURL, token string) error {
	log.Println("\n=== Testing User Authentication ===")

	var client *glclient.Client
	var err error

	if token != "" {
		log.Printf("Using token auth (token: %s...)", token[:20])
		client, err = glclient.NewClient(token, glclient.WithBaseURL(baseURL))
	} else {
		return fmt.Errorf("no authentication provided. Use token")
	}

	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	log.Printf("✓ Authentication successful")
	log.Printf("  User ID: %d", user.ID)
	log.Printf("  Username: %s", user.Username)
	log.Printf("  Email: %s", user.Email)
	log.Printf("  Admin: %v", user.IsAdmin)
	log.Printf("  State: %s", user.State)
	return nil
}

func testGroups(baseURL, token string) error {
	log.Println("\n=== Testing Groups Access ===")

	if token == "" {
		return fmt.Errorf("token required for groups test")
	}

	client, err := glclient.NewClient(token, glclient.WithBaseURL(baseURL))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	groups, _, err := client.Groups.ListGroups(&glclient.ListGroupsOptions{
		ListOptions: glclient.ListOptions{PerPage: 100},
	})
	if err != nil {
		return fmt.Errorf("failed to list groups: %w", err)
	}

	log.Printf("✓ Found %d groups:", len(groups))
	for _, g := range groups {
		log.Printf("  - %s (ID: %d, Path: %s)", g.Name, g.ID, g.Path)
	}
	return nil
}

func testProjects(baseURL, token string) error {
	log.Println("\n=== Testing Projects Access ===")

	if token == "" {
		return fmt.Errorf("token required for projects test")
	}

	client, err := glclient.NewClient(token, glclient.WithBaseURL(baseURL))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	projects, _, err := client.Projects.ListProjects(&glclient.ListProjectsOptions{
		ListOptions: glclient.ListOptions{PerPage: 100},
	})
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}

	log.Printf("✓ Found %d projects:", len(projects))
	for _, p := range projects {
		log.Printf("  - %s (ID: %d, Path: %s)", p.Name, p.ID, p.PathWithNamespace)
	}
	return nil
}
