package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/issue"
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
	var projectID, page, perPage int
	var state string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issue.List(mustConfig(cmd), projectID, state, page, perPage)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().StringVar(&state, "state", "opened", "Issue state: opened, closed, all")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}

func newIssueCreateCmd() *cobra.Command {
	var projectID, milestone int
	var title, description, labels string
	var assignees []int64
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issue.Create(mustConfig(cmd), projectID, title, description, labels, milestone, assignees)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().StringVar(&title, "title", "", "Issue title (required)")
	cmd.Flags().StringVar(&description, "description", "", "Issue description")
	cmd.Flags().StringVar(&labels, "labels", "", "Comma-separated labels")
	cmd.Flags().IntVar(&milestone, "milestone", 0, "Milestone ID")
	cmd.Flags().Int64SliceVar(&assignees, "assignee", nil, "Assignee user ID (repeatable for multiple)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newIssueCloseCmd() *cobra.Command {
	var projectID, iid int
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close an issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issue.Close(mustConfig(cmd), projectID, iid)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().IntVar(&iid, "iid", 0, "Issue IID (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("iid")
	return cmd
}

func newIssueBatchCmd() *cobra.Command {
	var projectID int
	var filePath string
	cmd := &cobra.Command{
		Use:   "batch",
		Short: "Batch create issues from a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return issue.Batch(mustConfig(cmd), projectID, filePath)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().StringVar(&filePath, "file", "", "JSON file path (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}
