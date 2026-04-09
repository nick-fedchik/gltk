package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/auth"
)

func newAuthCmd() *cobra.Command {
	var testCmd string
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Test GitLab authentication and connectivity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth.Test(mustConfig(cmd), testCmd)
		},
	}
	cmd.Flags().StringVar(&testCmd, "test", "user", "Test command: user, groups, projects, health")
	return cmd
}
