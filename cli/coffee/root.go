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
        `The index file which is used to quickly lookup a file's length and starting block. The file
is typically named "main_file_cache.idx{type}" where type corresponds to the type of file
that the index was built for.`,
    )
    _ = cacheCmd.MarkPersistentFlagRequired("index")

    cacheCmd.PersistentFlags().StringVar(
        &blocksFilePath,
        "blocks",
        "main_file_cache.dat2",
        `The blocks file which contains the data for all of the files contained within the cache.`,
    )

    rootCmd.AddCommand(cacheCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
