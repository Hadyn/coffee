package coffee

import (
    "bufio"
    "github.com/hadyn/coffee/jagex"
    "github.com/spf13/cobra"
    "io/ioutil"
    "os"
)

var archiveCmd = &cobra.Command{
    Use:   "archive",
    Short: "Root for archive editing commands",
    Long:  ``,

    Run: func(cmd *cobra.Command, args []string) {
        _ = cmd.Help()
    },
}

var archiveDecompressCmd = &cobra.Command{
    Use:   "decompress",
    Short: "Decompresses an archive read in from stdin",
    Long:  ``,
    RunE: func(cmd *cobra.Command, args []string) (err error) {
        var (
            bs []byte
            d  []byte
        )

        bs, err = ioutil.ReadAll(os.Stdin)
        if err != nil {
            return
        }

        d, err = jagex.DecompressFileArchive(bs)
        if err != nil {
            return
        }

        fw := bufio.NewWriter(os.Stdout)
        defer fw.Flush()

        _, err = fw.Write(d)
        return
    },
}
