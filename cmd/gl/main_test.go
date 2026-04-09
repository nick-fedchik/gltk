package main

import (
	"bytes"
	"testing"
)

func TestGLHelp(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--help"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gl --help failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected help output, got empty")
	}
}

func TestGLSubcommandsHelp(t *testing.T) {
	subcommands := [][]string{
		{"artifact", "--help"},
		{"auth", "--help"},
		{"branch", "--help"},
		{"comment", "--help"},
		{"commit", "--help"},
		{"diff", "--help"},
		{"file", "--help"},
		{"issue", "--help"},
		{"issues-check", "--help"},
		{"job", "--help"},
		{"label", "--help"},
		{"milestone", "--help"},
		{"mr", "--help"},
		{"pipeline", "--help"},
		{"project", "--help"},
		{"report", "--help"},
		{"runner", "--help"},
		{"search", "--help"},
		{"sync", "--help"},
		{"tag", "--help"},
		{"user", "--help"},
	}
	for _, args := range subcommands {
		cmd := NewRootCmd()
		cmd.SetArgs(args)
		var buf bytes.Buffer
		cmd.SetOut(&buf)
		if err := cmd.Execute(); err != nil {
			t.Errorf("gl %v failed: %v", args, err)
		}
	}
}
