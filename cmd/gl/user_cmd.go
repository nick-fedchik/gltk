package main

import (
	"github.com/spf13/cobra"
	"github.com/gltk/gltk/internal/gl/user"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage GitLab users",
	}
	cmd.AddCommand(
		newUserListCmd(),
		newUserGetCmd(),
		newUserCreateCmd(),
		newUserSetPasswordCmd(),
		newUserBlockCmd(),
		newUserUnblockCmd(),
		newUserDeleteCmd(),
	)
	return cmd
}

func newUserListCmd() *cobra.Command {
	var search string
	var active bool
	var page, perPage int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.List(mustConfig(cmd), search, active, page, perPage)
		},
	}
	cmd.Flags().StringVar(&search, "search", "", "Search by username, name, or email")
	cmd.Flags().BoolVar(&active, "active", false, "Show only active users")
	cmd.Flags().IntVar(&page, "page", 1, "Page number")
	cmd.Flags().IntVar(&perPage, "per-page", 50, "Results per page")
	return cmd
}

func newUserGetCmd() *cobra.Command {
	var username string
	var id int
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get user details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Get(mustConfig(cmd), username, id)
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Username to look up")
	cmd.Flags().IntVar(&id, "id", 0, "User ID to look up")
	return cmd
}

func newUserCreateCmd() *cobra.Command {
	var username, email, name, password string
	var admin, skipConfirm, resetPw bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Create(mustConfig(cmd), username, email, name, password, admin, skipConfirm, resetPw)
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Username (required)")
	cmd.Flags().StringVar(&email, "email", "", "Email (required)")
	cmd.Flags().StringVar(&name, "name", "", "Display name (required)")
	cmd.Flags().StringVar(&password, "password", "", "Password")
	cmd.Flags().BoolVar(&admin, "admin", false, "Make admin")
	cmd.Flags().BoolVar(&skipConfirm, "skip-confirmation", false, "Skip email confirmation")
	cmd.Flags().BoolVar(&resetPw, "reset-password", false, "Send password reset email")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newUserSetPasswordCmd() *cobra.Command {
	var username, password string
	var id int
	cmd := &cobra.Command{
		Use:   "set-password",
		Short: "Set user password",
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.SetPassword(mustConfig(cmd), username, id, password)
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Username")
	cmd.Flags().IntVar(&id, "id", 0, "User ID")
	cmd.Flags().StringVar(&password, "password", "", "New password (required)")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func newUserBlockCmd() *cobra.Command {
	var username string
	var id int
	cmd := &cobra.Command{
		Use:   "block",
		Short: "Block a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Block(mustConfig(cmd), username, id)
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Username")
	cmd.Flags().IntVar(&id, "id", 0, "User ID")
	return cmd
}

func newUserUnblockCmd() *cobra.Command {
	var username string
	var id int
	cmd := &cobra.Command{
		Use:   "unblock",
		Short: "Unblock a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Unblock(mustConfig(cmd), username, id)
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Username")
	cmd.Flags().IntVar(&id, "id", 0, "User ID")
	return cmd
}

func newUserDeleteCmd() *cobra.Command {
	var username string
	var id int
	var hardDelete bool
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Delete(mustConfig(cmd), username, id, hardDelete)
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Username")
	cmd.Flags().IntVar(&id, "id", 0, "User ID")
	cmd.Flags().BoolVar(&hardDelete, "hard-delete", false, "Hard delete (removes contributions)")
	return cmd
}
