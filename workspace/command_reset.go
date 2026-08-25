package main

import (
	"github.com/spf13/pflag"
)

// ResetFlags resets all flags (local and persistent) on the command and its children
// to their default values. This ensures isolation between repeated Execute() calls.
func (c *Command) ResetFlags() {
	c.resetFlags()
	for _, child := range c.commands {
		child.ResetFlags()
	}
}

// resetFlags resets the flags on this command only (not children)
func (c *Command) resetFlags() {
	// Reset local flags
	c.localFlags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		if f.Changed {
			f.Changed = false
		}
	})
	// Reset persistent flags
	c.persistentFlags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		if f.Changed {
			f.Changed = false
		}
	})
	// Reset inherited flags (flags from parent that are inherited)
	c.inheritedFlags().VisitAll(func(f *pflag.Flag) {
		f.Value.Set(f.DefValue)
		if f.Changed {
			f.Changed = false
		}
	})
}

// ExecuteC is the core execution function that adds flag reset before parsing
func (c *Command) ExecuteC() error {
	// Reset all flags to default values before each execution
	c.ResetFlags()
	
	// ... existing ExecuteC logic
	// This would include: ParseFlags, Run, etc.
	return nil
}

// Execute adds flag reset before the existing execution logic
func (c *Command) Execute() error {
	// Reset all flags to default values before each execution
	c.ResetFlags()
	
	// ... existing Execute logic
	return c.ExecuteC()
}

// ExecuteContext adds flag reset before the existing execution logic
func (c *Command) ExecuteContext(ctx context.Context) error {
	// Reset all flags to default values before each execution
	c.ResetFlags()
	
	// ... existing ExecuteContext logic
	return c.ExecuteC()
}
