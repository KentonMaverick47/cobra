package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
)

// ResetPersistentFlags resets all persistent flags to their default values.
// Call this before each Execute() invocation to prevent flag value leakage.
func ResetPersistentFlags(cmd *cobra.Command) {
    cmd.Flags().VisitAll(func(f *pflag.Flag) {
        if f.Changed {
            f.Value.Set(f.DefValue)
        }
    })
    // Also reset persistent flags on parent commands
    if cmd.HasParent() {
        ResetPersistentFlags(cmd.Parent())
    }
}

func main() {
    var rootCmd = &cobra.Command{
        Use:   "cobra-bounty",
        Short: "Fix persistent flag leakage demo",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello, Bounty Hunter!")
        },
    }
    
    // Reset flags before execution to prevent leakage
    ResetPersistentFlags(rootCmd)
    
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
