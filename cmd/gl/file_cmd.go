package main

import (
	"github.com/gltk/gltk/internal/gl/file"
	"github.com/spf13/cobra"
)

func newFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Read, write, and list files in GitLab repositories",
	}
	cmd.AddCommand(newFileReadCmd(), newFileWriteCmd(), newFileListCmd())
	return cmd
}

func newFileReadCmd() *cobra.Command {
	var project, filePath, ref, output string
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read a file from a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return file.Read(cfg, p, filePath, ref, output)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&filePath, "path", "", "File path (required)")
	cmd.Flags().StringVar(&ref, "ref", "main", "Git reference")
	cmd.Flags().StringVar(&output, "output", "", "Save to file (default: stdout)")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func newFileWriteCmd() *cobra.Command {
	var project, filePath, input, message, ref string
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Write a file to a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return file.Write(cfg, p, filePath, input, message, ref)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&filePath, "path", "", "File path (required)")
	cmd.Flags().StringVar(&input, "input", "", "Local file to upload (required)")
	cmd.Flags().StringVar(&message, "message", "", "Commit message (required)")
	cmd.Flags().StringVar(&ref, "ref", "main", "Git reference")
	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("message")
	return cmd
}

func newFileListCmd() *cobra.Command {
	var project, path, ref string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files in a repository directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return file.List(cfg, p, path, ref)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().StringVar(&path, "path", "", "Directory path")
	cmd.Flags().StringVar(&ref, "ref", "main", "Git reference")
	return cmd
}
