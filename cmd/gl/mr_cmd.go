package main

import (
	"github.com/gltk/gltk/internal/gl/mr"
	"github.com/spf13/cobra"
)

func newMRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mr",
		Short: "Manage GitLab merge requests",
	}
	cmd.AddCommand(
		newMRListCmd(),
		newMRGetCmd(),
		newMRCreateCmd(),
		newMRMergeCmd(),
		newMRCloseCmd(),
		newMRCommentCmd(),
	)
	return cmd
}

func newMRListCmd() *cobra.Command {
	var project, state string
	var page, perPage int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List merge requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return mr.List(cfg, p, state, page, perPage)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&state, "state", "opened", "MR state: opened, closed, merged, all")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	return cmd
}

func newMRGetCmd() *cobra.Command {
	var project string
	var mrIID int
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get merge request details",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return mr.Get(cfg, p, mrIID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	_ = cmd.MarkFlagRequired("mr")
	return cmd
}

func newMRCreateCmd() *cobra.Command {
	var project, title, source, target, description string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a merge request",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return mr.Create(cfg, p, title, source, target, description)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&title, "title", "", "MR title (required)")
	cmd.Flags().StringVar(&source, "source", "", "Source branch (required)")
	cmd.Flags().StringVar(&target, "target", "main", "Target branch")
	cmd.Flags().StringVar(&description, "description", "", "MR description")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("source")
	return cmd
}

func newMRMergeCmd() *cobra.Command {
	var project string
	var mrIID int
	var message string
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge a merge request",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return mr.Merge(cfg, p, mrIID, message)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	cmd.Flags().StringVar(&message, "message", "", "Merge commit message")
	_ = cmd.MarkFlagRequired("mr")
	return cmd
}

func newMRCloseCmd() *cobra.Command {
	var project string
	var mrIID int
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close a merge request",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return mr.Close(cfg, p, mrIID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	_ = cmd.MarkFlagRequired("mr")
	return cmd
}

func newMRCommentCmd() *cobra.Command {
	var project, body string
	var mrIID int
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Add a comment to a merge request",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return mr.Comment(cfg, p, mrIID, body)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	cmd.Flags().StringVar(&body, "body", "", "Comment body (required)")
	_ = cmd.MarkFlagRequired("mr")
	_ = cmd.MarkFlagRequired("body")
	return cmd
}
