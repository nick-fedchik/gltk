package file

import (
	"encoding/base64"
	"fmt"
	"os"

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

func parseProjectID(project string) (interface{}, error) {
	var id int
	if _, err := fmt.Sscanf(project, "%d", &id); err == nil {
		return id, nil
	}
	return project, nil
}

func Read(cfg *config.Config, project, filePath, ref, output string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}

	projectID, err := parseProjectID(project)
	if err != nil {
		return fmt.Errorf("error parsing project ID: %w", err)
	}

	fileInfo, resp, err := client.RepositoryFiles.GetFile(projectID, filePath, &glclient.GetFileOptions{Ref: glclient.Ptr(ref)})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	content, err := base64.StdEncoding.DecodeString(fileInfo.Content)
	if err != nil {
		return fmt.Errorf("failed to decode file content: %w", err)
	}

	if output != "" {
		if err := os.WriteFile(output, content, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("✓ File written to %s (%d bytes)\n", output, len(content))
	} else {
		fmt.Print(string(content))
	}
	return nil
}

func Write(cfg *config.Config, project, filePath, input, message, ref string) error {
	content, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	encodedContent := base64.StdEncoding.EncodeToString(content)

	client, err := newClient(cfg)
	if err != nil {
		return err
	}

	projectID, err := parseProjectID(project)
	if err != nil {
		return fmt.Errorf("error parsing project ID: %w", err)
	}

	// Check if file exists
	_, resp, _ := client.RepositoryFiles.GetFile(projectID, filePath, &glclient.GetFileOptions{Ref: glclient.Ptr(ref)})
	fileExists := resp != nil && resp.StatusCode == 200

	var fileInfo *glclient.FileInfo
	if fileExists {
		fileInfo, _, err = client.RepositoryFiles.UpdateFile(projectID, filePath, &glclient.UpdateFileOptions{
			Content:       glclient.Ptr(encodedContent),
			CommitMessage: glclient.Ptr(message),
			Branch:        glclient.Ptr(ref),
			Encoding:      glclient.Ptr("base64"),
		})
		if err != nil {
			return fmt.Errorf("failed to update file: %w", err)
		}
		fmt.Printf("✓ File updated\n")
	} else {
		fileInfo, _, err = client.RepositoryFiles.CreateFile(projectID, filePath, &glclient.CreateFileOptions{
			Content:       glclient.Ptr(encodedContent),
			CommitMessage: glclient.Ptr(message),
			Branch:        glclient.Ptr(ref),
			Encoding:      glclient.Ptr("base64"),
		})
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		fmt.Printf("✓ File created\n")
	}

	fmt.Printf("  Path: %s\n", fileInfo.FilePath)
	fmt.Printf("  Size: %d bytes\n", len(content))
	return nil
}

func List(cfg *config.Config, project, path, ref string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}

	projectID, err := parseProjectID(project)
	if err != nil {
		return fmt.Errorf("error parsing project ID: %w", err)
	}

	entries, _, err := client.Repositories.ListTree(projectID, &glclient.ListTreeOptions{
		Path:      glclient.Ptr(path),
		Ref:       glclient.Ptr(ref),
		Recursive: glclient.Ptr(false),
	})
	if err != nil {
		return fmt.Errorf("failed to list directory: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("Directory is empty: %s\n", path)
		return nil
	}

	fmt.Printf("Contents of %s (ref: %s):\n\n", path, ref)
	for _, entry := range entries {
		icon := "📄"
		if entry.Type == "tree" {
			icon = "📁"
		} else if entry.Type == "commit" {
			icon = "🔗"
		}

		fmt.Printf("%s  %-40s  %s\n", icon, entry.Name, entry.ID[:8])
	}
	return nil
}
