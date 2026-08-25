package main

import (
	"testing"
)

func TestPersistentFlagLeakage(t *testing.T) {
	var persistentVal string

	rootCmd := &Command{Use: "root"}
	rootCmd.PersistentFlags().StringVar(&persistentVal, "config", "default-path", "config file path")

	subCmd := &Command{
		Use: "sub",
		Run: func(cmd *Command, args []string) {},
	}
	rootCmd.AddCommand(subCmd)

	// First execution: Explicitly set the persistent flag
	rootCmd.SetArgs([]string{"sub", "--config", "custom-path"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("First execution failed: %v", err)
	}
	if persistentVal != "custom-path" {
		t.Errorf("Expected 'custom-path', got '%s'", persistentVal)
	}

	// Second execution: Do not set the persistent flag
	rootCmd.SetArgs([]string{"sub"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Second execution failed: %v", err)
	}
	if persistentVal != "default-path" {
		t.Errorf("Expected flag to reset to 'default-path', but it leaked: '%s'", persistentVal)
	}
}

func TestLocalFlagLeakage(t *testing.T) {
	var localVal string

	rootCmd := &Command{Use: "root"}
	rootCmd.Flags().StringVar(&localVal, "local", "default-local", "local flag")

	rootCmd.Run = func(cmd *Command, args []string) {}

	// First execution: Explicitly set the local flag
	rootCmd.SetArgs([]string{"--local", "custom-local"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("First execution failed: %v", err)
	}
	if localVal != "custom-local" {
		t.Errorf("Expected 'custom-local', got '%s'", localVal)
	}

	// Second execution: Do not set the local flag
	rootCmd.SetArgs([]string{})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Second execution failed: %v", err)
	}
	if localVal != "default-local" {
		t.Errorf("Expected flag to reset to 'default-local', but it leaked: '%s'", localVal)
	}
}

func TestChildCommandInheritanceIsolation(t *testing.T) {
	var persistentVal string

	rootCmd := &Command{Use: "root"}
	rootCmd.PersistentFlags().StringVar(&persistentVal, "config", "default-path", "config file path")

	subCmd1 := &Command{
		Use: "sub1",
		Run: func(cmd *Command, args []string) {},
	}
	subCmd2 := &Command{
		Use: "sub2",
		Run: func(cmd *Command, args []string) {},
	}
	rootCmd.AddCommand(subCmd1, subCmd2)

	// First execution: sub1 with custom config
	rootCmd.SetArgs([]string{"sub1", "--config", "custom-path-1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("First execution failed: %v", err)
	}
	if persistentVal != "custom-path-1" {
		t.Errorf("Expected 'custom-path-1', got '%s'", persistentVal)
	}

	// Second execution: sub2 without config (should use default)
	rootCmd.SetArgs([]string{"sub2"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Second execution failed: %v", err)
	}
	if persistentVal != "default-path" {
		t.Errorf("Expected flag to reset to 'default-path' for sub2, but it leaked: '%s'", persistentVal)
	}

	// Third execution: sub1 again without config (should use default)
	rootCmd.SetArgs([]string{"sub1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Third execution failed: %v", err)
	}
	if persistentVal != "default-path" {
		t.Errorf("Expected flag to reset to 'default-path' for sub1, but it leaked: '%s'", persistentVal)
	}
}

func TestBoundVariableReset(t *testing.T) {
	var intVal int
	var boolVal bool
	var stringVal string

	rootCmd := &Command{Use: "root"}
	rootCmd.PersistentFlags().IntVar(&intVal, "int-flag", 42, "int flag")
	rootCmd.PersistentFlags().BoolVar(&boolVal, "bool-flag", true, "bool flag")
	rootCmd.PersistentFlags().StringVar(&stringVal, "string-flag", "hello", "string flag")

	rootCmd.Run = func(cmd *Command, args []string) {}

	// First execution: change all values
	rootCmd.SetArgs([]string{"--int-flag", "100", "--bool-flag", "false", "--string-flag", "world"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("First execution failed: %v", err)
	}
	if intVal != 100 || boolVal != false || stringVal != "world" {
		t.Errorf("Expected 100, false, world; got %d, %v, %s", intVal, boolVal, stringVal)
	}

	// Second execution: no flags provided, should reset to defaults
	rootCmd.SetArgs([]string{})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Second execution failed: %v", err)
	}
	if intVal != 42 || boolVal != true || stringVal != "hello" {
		t.Errorf("Expected flags to reset to defaults (42, true, hello); got %d, %v, %s", intVal, boolVal, stringVal)
	}
}
