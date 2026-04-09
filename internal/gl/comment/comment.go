package comment

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

func List(cfg *config.Config, project, resourceType string, resourceID int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}

	// Convert to int for API call
	var projInt int
	if id, err := strconv.Atoi(project); err == nil {
		projInt = id
	} else {
		return fmt.Errorf("project must be numeric ID or convertible path")
	}

	var notes []*glclient.Note

	if resourceType == "mr" {
		notes, _, err = client.Notes.ListMergeRequestNotes(projInt, int64(resourceID), nil)
	} else {
		notes, _, err = client.Notes.ListIssueNotes(projInt, int64(resourceID), nil)
	}

	if err != nil {
		return fmt.Errorf("failed to list comments: %w", err)
	}

	if len(notes) == 0 {
		fmt.Printf("No comments on %s #%d\n", resourceType, resourceID)
		return nil
	}

	fmt.Printf("Comments on %s #%d:\n\n", resourceType, resourceID)
	for _, note := range notes {
		fmt.Printf("#%d | %s | %s\n", note.ID, note.Author.Username, note.CreatedAt)
		fmt.Printf("   %s\n\n", note.Body)
	}
	return nil
}

func Add(cfg *config.Config, project, resourceType string, resourceID int, body string) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}

	// Convert to int for API call
	var projInt int
	if i, err := strconv.Atoi(project); err == nil {
		projInt = i
	} else {
		return fmt.Errorf("project must be numeric ID")
	}

	var note *glclient.Note

	if resourceType == "mr" {
		opts := &glclient.CreateMergeRequestNoteOptions{
			Body: glclient.Ptr(body),
		}
		note, _, err = client.Notes.CreateMergeRequestNote(projInt, int64(resourceID), opts)
	} else {
		opts := &glclient.CreateIssueNoteOptions{
			Body: glclient.Ptr(body),
		}
		note, _, err = client.Notes.CreateIssueNote(projInt, int64(resourceID), opts)
	}

	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	fmt.Printf("✓ Comment added to %s #%d\n", resourceType, resourceID)
	fmt.Printf("  ID: %d\n", note.ID)
	fmt.Printf("  Author: %s\n", note.Author.Username)
	return nil
}

func Delete(cfg *config.Config, project, resourceType string, resourceID, noteID int) error {
	client, err := newClient(cfg)
	if err != nil {
		return err
	}

	// Convert to int for API call
	var projInt int
	if i, err := strconv.Atoi(project); err == nil {
		projInt = i
	} else {
		return fmt.Errorf("project must be numeric ID")
	}

	if resourceType == "mr" {
		_, err = client.Notes.DeleteMergeRequestNote(projInt, int64(resourceID), int64(noteID))
	} else {
		_, err = client.Notes.DeleteIssueNote(projInt, int64(resourceID), int64(noteID))
	}

	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	fmt.Printf("✓ Comment #%d deleted from %s #%d\n", noteID, resourceType, resourceID)
	return nil
}
