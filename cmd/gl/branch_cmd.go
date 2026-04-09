package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/branch"
)

func newBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Manage GitLab repository branches",
	}
	cmd.AddCommand(
		newBranchListCmd(),
		newBranchGetCmd(),
		newBranchCreateCmd(),
		newBranchDeleteCmd(),
		newBranchProtectCmd(),
	)
	return cmd
}

func newBranchListCmd() *cobra.Command {
	var project string
	var page, perPage int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			return branch.List(mustConfig(cmd), project, page, perPage)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}

func newBranchGetCmd() *cobra.Command {
	var project, branchName string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get branch details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return branch.Get(mustConfig(cmd), project, branchName)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&branchName, "branch", "", "Branch name (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("branch")
	return cmd
}

func newBranchCreateCmd() *cobra.Command {
	var project, branchName, ref string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			return branch.Create(mustConfig(cmd), project, branchName, ref)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&branchName, "branch", "", "New branch name (required)")
	cmd.Flags().StringVar(&ref, "ref", "main", "Reference to branch from")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("branch")
	return cmd
}

func newBranchDeleteCmd() *cobra.Command {
	var project, branchName string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			return branch.Delete(mustConfig(cmd), project, branchName)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&branchName, "branch", "", "Branch name (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("branch")
	return cmd
}

func newBranchProtectCmd() *cobra.Command {
	var project, branchName string
	cmd := &cobra.Command{
		Use:   "protect",
		Short: "Protect a branch",
		RunE: func(cmd *cobra.Command, args []string) error {
			return branch.Protect(mustConfig(cmd), project, branchName)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&branchName, "branch", "", "Branch name (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("branch")
	return cmd
}
