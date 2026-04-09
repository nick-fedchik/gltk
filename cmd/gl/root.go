package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gltk/gltk/internal/config"
	"github.com/spf13/cobra"
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

// resolveProject returns the project from the flag, falling back to cfg.ProjectID.
// Returns an error if neither is set.
func resolveProject(cfg *config.Config, flagValue string) (string, error) {
	p := strings.TrimRight(flagValue, "/")
	if p == "" {
		p = strings.TrimRight(cfg.ProjectID, "/")
	}
	if p == "" {
		return "", fmt.Errorf("project is required: use --project flag or set default-project in config.yaml")
	}
	return p, nil
}
