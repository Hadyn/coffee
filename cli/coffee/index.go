package coffee

import (
    "bufio"
    "encoding/json"
    "github.com/hadyn/coffee/jagex"
    "github.com/spf13/cobra"
    "io/ioutil"
    "os"
)

var indexCmd = &cobra.Command{
    Use:   "index",
    Short: "Root for file index editing commands",
    Long:  ``,

    Run: func(cmd *cobra.Command, args []string) {
        _ = cmd.Help()
    },
}

var indexDecode = &cobra.Command{
    Use:   "decode",
    Short: "Decodes a file index into an output type",
    Long:  ``,

    RunE: func(cmd *cobra.Command, args []string) (err error) {
        var (
            bs  []byte
            fi  *jagex.FileIndex
            enc []byte
        )

        bs, err = ioutil.ReadAll(os.Stdin)
        if err != nil {
            return
        }

        fi, err = jagex.DecodeFileIndex(bs)
        if err != nil {
            return
        }

        if enc, err = json.Marshal(fi); err != nil {
            return
        }

        fw := bufio.NewWriter(os.Stdout)
        defer fw.Flush()

        _, err = fw.Write(enc)

        return nil
    },
}
