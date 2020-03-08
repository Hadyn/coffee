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

	archiveCmd.AddCommand(archiveDecompressCmd)

	indexCmd.AddCommand(indexDecodeCmd)
	indexCmd.AddCommand(indexLookupCmd)

	spriteDecodeCmd.PersistentFlags().StringVarP(
		&spriteDecodeFormat,
		"format",
		"f",
		"png",
		"",
	)

	spriteCmd.AddCommand(spriteDecodeCmd)

	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(spriteCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
