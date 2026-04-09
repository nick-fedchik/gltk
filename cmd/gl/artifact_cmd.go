package main

import (
	"github.com/gltk/gltk/internal/gl/artifact"
	"github.com/spf13/cobra"
)

func newArtifactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "Manage GitLab job artifacts",
	}
	cmd.AddCommand(newArtifactListCmd(), newArtifactDownloadCmd(), newArtifactDeleteCmd())
	return cmd
}

func newArtifactListCmd() *cobra.Command {
	var project string
	var jobID, page, perPage int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List job artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return artifact.List(cfg, p, jobID, page, perPage)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (optional)")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	return cmd
}

func newArtifactDownloadCmd() *cobra.Command {
	var project string
	var jobID int
	var output string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download job artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return artifact.Download(cfg, p, jobID, output)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newArtifactDeleteCmd() *cobra.Command {
	var project string
	var jobID int
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete job artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustConfig(cmd)
			p, err := resolveProject(cfg, project)
			if err != nil {
				return err
			}
			return artifact.Delete(cfg, p, jobID)
		},
	}
	cmd.Flags().StringVar(&project, "project", "", "Project ID or path (default: from config)")
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}
