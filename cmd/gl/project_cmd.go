package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/project"
)

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "View and query GitLab projects",
	}
	cmd.AddCommand(newProjectListCmd(), newProjectGetCmd(), newProjectByPathCmd())
	return cmd
}

func newProjectListCmd() *cobra.Command {
	var groupID int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects in a group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return project.List(mustConfig(cmd), groupID)
		},
	}
	cmd.Flags().IntVar(&groupID, "group", 0, "Group ID")
	return cmd
}

func newProjectGetCmd() *cobra.Command {
	var projectID int
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get project by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			return project.GetByID(mustConfig(cmd), projectID)
		},
	}
	cmd.Flags().IntVar(&projectID, "id", 0, "Project ID (required)")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newProjectByPathCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "by-path",
		Short: "Get project by path (e.g. group/project)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return project.ByPath(mustConfig(cmd), path)
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "Project path (required)")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}
