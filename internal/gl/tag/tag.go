package tag

import (
	"fmt"
	"strconv"

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
	id, err := strconv.Atoi(project)
	if err == nil {
		return id
	}
	return project
}

func List(cfg *config.Config, project string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.ListTagsOptions{
		ListOptions: glclient.ListOptions{Page: 1, PerPage: 100},
	}

	tags, _, err := client.Tags.ListTags(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found")
		return nil
	}

	fmt.Printf("Tags (%d total):\n\n", len(tags))
	for _, t := range tags {
		shortID := t.Commit.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		fmt.Printf("  %-20s  %s  %s\n", t.Name, shortID, t.Message)
	}
	return nil
}

func Create(cfg *config.Config, project, tagName, ref, message string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	opts := &glclient.CreateTagOptions{
		TagName: glclient.Ptr(tagName),
		Ref:     glclient.Ptr(ref),
		Message: glclient.Ptr(message),
	}

	created, _, err := client.Tags.CreateTag(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	fmt.Printf("✓ Tag created\n")
	fmt.Printf("  Name: %s\n", created.Name)
	fmt.Printf("  Commit: %s\n", created.Commit.ID[:8])
	return nil
}

func Delete(cfg *config.Config, project, tagName string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}
	projectID := getProjectID(project)

	_, err = client.Tags.DeleteTag(projectID, tagName)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	fmt.Printf("✓ Tag deleted\n")
	fmt.Printf("  Name: %s\n", tagName)
	return nil
}
