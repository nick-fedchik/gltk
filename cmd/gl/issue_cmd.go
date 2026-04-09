package main

import (
	"github.com/gltk/gltk/internal/gl/issue"
	"github.com/spf13/cobra"
)

func newIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage GitLab issues",
	}
	cmd.AddCommand(newIssueListCmd(), newIssueCreateCmd(), newIssueCloseCmd(), newIssueBatchCmd())
	return cmd
}

func newIssueListCmd() *cobra.Command {
	var project string
	var page, perPage int
	var state string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return issue.List(cfg, p, state, page, perPage)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&state, "state", "opened", "Issue state: opened, closed, all")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	return cmd
}

func newIssueCreateCmd() *cobra.Command {
	var project string
	var milestone int
	var title, description, labels string
	var assignees []int64
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return issue.Create(cfg, p, title, description, labels, milestone, assignees)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&title, "title", "", "Issue title (required)")
	cmd.Flags().StringVar(&description, "description", "", "Issue description")
	cmd.Flags().StringVar(&labels, "labels", "", "Comma-separated labels")
	cmd.Flags().IntVar(&milestone, "milestone", 0, "Milestone ID")
	cmd.Flags().Int64SliceVar(&assignees, "assignee", nil, "Assignee user ID (repeatable for multiple)")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newIssueCloseCmd() *cobra.Command {
	var project string
	var iid int
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close an issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return issue.Close(cfg, p, iid)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&iid, "iid", 0, "Issue IID (required)")
	_ = cmd.MarkFlagRequired("iid")
	return cmd
}

func newIssueBatchCmd() *cobra.Command {
	var project string
	var filePath string
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Batch create issues from a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return issue.Batch(cfg, p, filePath)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&filePath, "file", "", "JSON file path (required)")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
