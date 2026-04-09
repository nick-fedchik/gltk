package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/tag"
)

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage GitLab repository tags",
	}
	cmd.AddCommand(newTagListCmd(), newTagCreateCmd(), newTagDeleteCmd())
	return cmd
}

func newTagListCmd() *cobra.Command {
	var project string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repository tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tag.List(mustConfig(cmd), project)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}

func newTagCreateCmd() *cobra.Command {
	var project, tagName, ref, message string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tag.Create(mustConfig(cmd), project, tagName, ref, message)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&tagName, "tag", "", "Tag name (required)")
	cmd.Flags().StringVar(&ref, "ref", "main", "Ref to tag")
	cmd.Flags().StringVar(&message, "message", "", "Tag message")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("tag")
	return cmd
}

func newTagDeleteCmd() *cobra.Command {
	var project, tagName string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tag.Delete(mustConfig(cmd), project, tagName)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (required)")
	cmd.Flags().StringVar(&tagName, "tag", "", "Tag name (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("tag")
	return cmd
}
