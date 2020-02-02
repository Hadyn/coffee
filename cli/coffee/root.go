package coffee

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "coffee",
    Short: "Coffee is a Runescape cache and file editing toolset",
    Long:  ``,
    Run: func(cmd *cobra.Command, args []string) {
        _ = cmd.Help()
    },
}

func Execute() {
    cacheCmd.AddCommand(cacheReadCmd)

    cacheCmd.PersistentFlags().StringVar(
        &indexFilePath,
        "index",
        "",
        indexFlagUsage,
    )
    _ = cacheCmd.MarkPersistentFlagRequired("index")

    cacheCmd.PersistentFlags().StringVar(
        &blocksFilePath,
        "blocks",
        "main_file_cache.dat2",
        blocksFlagUsage,
    )

    rootCmd.AddCommand(cacheCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
