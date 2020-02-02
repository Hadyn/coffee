package coffee

import (
    "bufio"
    "errors"
    "github.com/hadyn/coffee/jagex"
    "github.com/spf13/cobra"
    "os"
    "strconv"
)

var (
    indexFilePath  string
    blocksFilePath string
)

var cacheCmd = &cobra.Command{
    Use:   "cache",
    Short: "Root for cache editing commands",
    Long:  ``,

    Run: func(cmd *cobra.Command, args []string) {
        _ = cmd.Help()
    },
}


var cacheReadCmd = &cobra.Command{
    Use:   "read <file-type> <file-id>",
    Short: "Reads a file from a cache",
    Long: ``,
    Args: cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) (err error) {
        var (
            fileType int
            fileID   int
        )

        fileType, err = strconv.Atoi(args[0])
        if err != nil {
            return
        }

        fileID, err = strconv.Atoi(args[1])
        if err != nil {
            return
        }

        var (
            index  *os.File
            blocks *os.File
        )

        index, err = openIndexFile(os.O_RDONLY, 0777)
        if err != nil {
            return
        }

        defer index.Close()

        blocks, err = openBlocksFile(os.O_RDONLY, 0777)
        if err != nil {
            return
        }

        defer blocks.Close()

        r := jagex.NewCacheReader(index, blocks, fileType)

        var file []byte
        file, err = r.Read(fileID)
        if err != nil {
            return
        }

        fw := bufio.NewWriter(os.Stdout)
        defer fw.Flush()

        _, err = fw.Write(file)
        return
    },
}

func openIndexFile(flag int, perm os.FileMode) (file *os.File, err error) {
    if indexFilePath == "" {
        return nil, errors.New("cache index file path was not set")
    }
    file, err = os.OpenFile(indexFilePath, flag, perm)
    return
}

func openBlocksFile(flag int, perm os.FileMode) (file *os.File, err error) {
    if blocksFilePath == "" {
        return nil, errors.New("cache blocks file path was not set")
    }
    file, err = os.OpenFile(blocksFilePath, flag, perm)
    return
}