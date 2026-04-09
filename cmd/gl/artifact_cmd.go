package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/artifact"
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
	var projectID, jobID, page, perPage int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List job artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return artifact.List(mustConfig(cmd), projectID, jobID, page, perPage)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (optional)")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 20, "Results per page")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}

func newArtifactDownloadCmd() *cobra.Command {
	var projectID, jobID int
	var output string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download job artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return artifact.Download(mustConfig(cmd), projectID, jobID, output)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	cmd.Flags().StringVar(&output, "output", "", "Output file path")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}

func newArtifactDeleteCmd() *cobra.Command {
	var projectID, jobID int
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete job artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return artifact.Delete(mustConfig(cmd), projectID, jobID)
		},
	}
	cmd.Flags().IntVar(&projectID, "project", 0, "Project ID (required)")
	cmd.Flags().IntVar(&jobID, "job", 0, "Job ID (required)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("job")
	return cmd
}
