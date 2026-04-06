package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestInitRequiresRepoFlag(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"init", "demo"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected init without --repo to fail")
	}
	if !strings.Contains(err.Error(), "required flag(s) \"repo\" not set") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddRequiresRepoFlag(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"add"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected add without --repo to fail")
	}
	if !strings.Contains(err.Error(), "required flag(s) \"repo\" not set") {
		t.Fatalf("unexpected error: %v", err)
	}
}
