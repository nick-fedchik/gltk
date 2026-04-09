package commit

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gltk/gltk/internal/config"
	glclient "gitlab.com/gitlab-org/api/client-go"
)

// ActionSpec describes a single file action in a commit.
type ActionSpec struct {
	Action    string `json:"action"`     // create | update | delete | move | chmod
	FilePath  string `json:"file_path"`  // path in repository (required)
	LocalFile string `json:"local_file"` // path on local disk (required for create/update)
	Content   string `json:"content"`    // inline content (alternative to local_file)
}

// CommitSpec is the JSON input format for gl-commit.
type CommitSpec struct {
	Branch  string       `json:"branch"`
	Message string       `json:"message"`
	Actions []ActionSpec `json:"actions"`
}

func Create(cfg *config.Config, projectID interface{}, specFile string) error {
	// Read commit spec
	data, err := os.ReadFile(specFile)
	if err != nil {
		return fmt.Errorf("failed to read commit spec: %w", err)
	}

	var spec CommitSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("failed to parse commit spec JSON: %w", err)
	}

	if spec.Branch == "" {
		return fmt.Errorf("commit spec must include 'branch'")
	}
	if spec.Message == "" {
		return fmt.Errorf("commit spec must include 'message'")
	}
	if len(spec.Actions) == 0 {
		return fmt.Errorf("commit spec must include at least one action")
	}

	client, err := cfg.NewGitLabClient()
	if err != nil {
		return err
	}

	// Build API actions
	actions := make([]*glclient.CommitActionOptions, 0, len(spec.Actions))
	for _, a := range spec.Actions {
		if a.FilePath == "" {
			return fmt.Errorf("each action must have 'file_path'")
		}

		action := glclient.FileActionValue(a.Action)
		opt := &glclient.CommitActionOptions{
			Action:   &action,
			FilePath: glclient.Ptr(a.FilePath),
		}

		// Load content from local file if specified
		if a.LocalFile != "" {
			content, err := os.ReadFile(a.LocalFile)
			if err != nil {
				return fmt.Errorf("failed to read local file %s: %w", a.LocalFile, err)
			}
			s := string(content)
			opt.Content = &s
		} else if a.Content != "" {
			opt.Content = &a.Content
		}

		actions = append(actions, opt)
	}

	opts := &glclient.CreateCommitOptions{
		Branch:        glclient.Ptr(spec.Branch),
		CommitMessage: glclient.Ptr(spec.Message),
		Actions:       actions,
	}

	c, _, err := client.Commits.CreateCommit(projectID, opts)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	fmt.Printf("✅ Commit created\n")
	fmt.Printf("  SHA:     %s\n", c.ID)
	fmt.Printf("  Branch:  %s\n", spec.Branch)
	fmt.Printf("  Message: %s\n", c.Title)
	fmt.Printf("  Files:   %d actions\n", len(spec.Actions))
	fmt.Printf("  URL:     %s\n", c.WebURL)
	return nil
}
