package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/config"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "gl",
		Short:   "GitLab API client — issues, MRs, pipelines, jobs, and more",
		Version: "2.0.0",
	}

	root.PersistentFlags().String("config", "", "Config file path (default: config.yaml)")
	root.PersistentFlags().String("gitlab-url", "", "GitLab instance URL (overrides config/env)")
	root.PersistentFlags().String("token", "", "GitLab personal access token (overrides config/env)")
	root.PersistentFlags().String("format", "text", "Output format: text, json, table")
	root.PersistentFlags().Bool("verbose", false, "Verbose output")
	root.PersistentFlags().Bool("no-color", false, "Disable color output")

	root.AddCommand(
		newArtifactCmd(),
		newAuthCmd(),
		newBranchCmd(),
		newCommentCmd(),
		newCommitCmd(),
		newDiffCmd(),
		newFileCmd(),
		newIssueCmd(),
		newIssuesCheckCmd(),
		newJobCmd(),
		newLabelCmd(),
		newMilestoneCmd(),
		newMRCmd(),
		newPipelineCmd(),
		newProjectCmd(),
		newReportCmd(),
		newRunnerCmd(),
		newSearchCmd(),
		newTagCmd(),
		newUserCmd(),
	)

	return root
}

func mustConfig(cmd *cobra.Command) *config.Config {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	if url, _ := cmd.Flags().GetString("gitlab-url"); url != "" {
		cfg.GitLabURL = url
	}
	if token, _ := cmd.Flags().GetString("token"); token != "" {
		cfg.Token = token
	}
	return cfg
}
