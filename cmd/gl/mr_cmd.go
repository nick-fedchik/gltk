package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/mr"
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
			return mr.List(mustConfig(cmd), project, state, page, perPage)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&state, "state", "opened", "MR state: opened, closed, merged, all")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}

func newMRGetCmd() *cobra.Command {
	var project string
	var mrIID int
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get merge request details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mr.Get(mustConfig(cmd), project, mrIID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("mr")
	return cmd
}

func newMRCreateCmd() *cobra.Command {
	var project, title, source, target, description string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a merge request",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mr.Create(mustConfig(cmd), project, title, source, target, description)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&title, "title", "", "MR title (required)")
	cmd.Flags().StringVar(&source, "source", "", "Source branch (required)")
	cmd.Flags().StringVar(&target, "target", "main", "Target branch")
	cmd.Flags().StringVar(&description, "description", "", "MR description")
	_ = cmd.MarkFlagRequired("project")
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
			return mr.Merge(mustConfig(cmd), project, mrIID, message)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	cmd.Flags().StringVar(&message, "message", "", "Merge commit message")
	_ = cmd.MarkFlagRequired("project")
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
			return mr.Close(mustConfig(cmd), project, mrIID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	_ = cmd.MarkFlagRequired("project")
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
			return mr.Comment(mustConfig(cmd), project, mrIID, body)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().IntVar(&mrIID, "mr", 0, "MR IID (required)")
	cmd.Flags().StringVar(&body, "body", "", "Comment body (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("mr")
	_ = cmd.MarkFlagRequired("body")
	return cmd
}
