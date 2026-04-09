package main

import (
	"github.com/gltk/gltk/internal/gl/comment"
	"github.com/spf13/cobra"
)

func newCommentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Manage GitLab issue and MR comments",
	}
	cmd.AddCommand(newCommentListCmd(), newCommentAddCmd(), newCommentDeleteCmd())
	return cmd
}

func newCommentListCmd() *cobra.Command {
	var project, resourceType string
	var resourceID int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List comments on an issue or MR",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return comment.List(cfg, p, resourceType, resourceID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&resourceType, "type", "issue", "Resource type: issue, mr")
	cmd.Flags().IntVar(&resourceID, "id", 0, "Issue or MR IID (required)")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newCommentAddCmd() *cobra.Command {
	var project, resourceType, body string
	var resourceID int
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a comment to an issue or MR",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return comment.Add(cfg, p, resourceType, resourceID, body)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&resourceType, "type", "issue", "Resource type: issue, mr")
	cmd.Flags().IntVar(&resourceID, "id", 0, "Issue or MR IID (required)")
	cmd.Flags().StringVar(&body, "body", "", "Comment body (required)")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("body")
	return cmd
}

func newCommentDeleteCmd() *cobra.Command {
	var project, resourceType string
	var resourceID, noteID int
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return comment.Delete(cfg, p, resourceType, resourceID, noteID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&resourceType, "type", "issue", "Resource type: issue, mr")
	cmd.Flags().IntVar(&resourceID, "id", 0, "Issue or MR IID (required)")
	cmd.Flags().IntVar(&noteID, "note", 0, "Note ID (required)")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("note")
	return cmd
}
