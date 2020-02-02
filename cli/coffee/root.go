package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "coffee",
    Short: "Coffee is a Runescape cache editing toolset",
    Long: ``,
    Run: func(cmd *cobra.Command, args []string) {
    },
}

func Execute() {
    
}