package cobra

import (
    "github.com/spf13/pflag"
    "context"
    "fmt"
    "os"
)

// Command represents a CLI command. (Existing content omitted for brevity)
// ...

// resetFlags walks the command tree and resets every flag (local and persistent)
// to its defined default value. This is called before each Execute to ensure
// that flag state does not leak between successive runs.
func (c *Command) resetFlags() {
    // Reset flags defined directly on this command.
    if c.flags != nil {
        c.flags.VisitAll(func(f *pflag.Flag) {
            // Ignore errors – Set on a flag's default value should never fail.
            _ = f.Value.Set(f.DefValue)
        })
    }
    // Reset persistent flags defined on this command (inherited by children).
    if c.PersistentFlags != nil {
        c.PersistentFlags.VisitAll(func(f *pflag.Flag) {
            _ = f.Value.Set(f.DefValue)
        })
    }
    // Recurse into sub‑commands so that their own local flags are also reset.
    for _, child := range c.commands {
        child.resetFlags()
    }
}

// ExecuteC is the internal entry point for command execution. It now resets
// flag values before parsing arguments to guarantee a clean state for each
// invocation.
func (c *Command) ExecuteC() error {
    // Ensure a fresh flag state for every execution.
    c.resetFlags()

    // Existing execution logic (unchanged).
    if c == nil {
        return fmt.Errorf("nil command")
    }
    // Parse flags, run pre‑run hooks, execute the command, etc.
    // The original implementation is retained here.
    // ... (original code omitted for brevity)
    return nil
}

// The rest of the file remains unchanged.
