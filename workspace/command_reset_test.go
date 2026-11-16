package workspace

import (
    "testing"
    "github.com/spf13/cobra"
)

func TestPersistentFlagLeakage(t *testing.T) {
    var persistentVal string

    rootCmd := &cobra.Command{Use: "root"}
    rootCmd.PersistentFlags().StringVar(&persistentVal, "config", "default-path", "config file path")

    subCmd := &cobra.Command{
        Use: "sub",
        Run: func(cmd *cobra.Command, args []string) {},
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
