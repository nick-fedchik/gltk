package main

import (
	"github.com/gltk/gltk/internal/gl/tag"
	"github.com/spf13/cobra"
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
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return tag.List(cfg, p)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	return cmd
}

func newTagCreateCmd() *cobra.Command {
	var project, tagName, ref, message string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return tag.Create(cfg, p, tagName, ref, message)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&tagName, "tag", "", "Tag name (required)")
	cmd.Flags().StringVar(&ref, "ref", "main", "Ref to tag")
	cmd.Flags().StringVar(&message, "message", "", "Tag message")
	_ = cmd.MarkFlagRequired("tag")
	return cmd
}

func newTagDeleteCmd() *cobra.Command {
	var project, tagName string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return tag.Delete(cfg, p, tagName)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&tagName, "tag", "", "Tag name (required)")
	_ = cmd.MarkFlagRequired("tag")
	return cmd
}
