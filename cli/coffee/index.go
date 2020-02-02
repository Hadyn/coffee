package coffee

import (
    "bufio"
    "encoding/json"
    "fmt"
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

var indexDecodeCmd = &cobra.Command{
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

var indexLookupCmd = &cobra.Command{
    Use:   "lookup <name>",
    Short: "Looks up the identifier of a group for the provided name",
    Long:  ``,

    RunE: func(cmd *cobra.Command, args []string) (err error) {
        var (
            bs  []byte
            fi  *jagex.FileIndex
        )

        bs, err = ioutil.ReadAll(os.Stdin)
        if err != nil {
            return
        }

        fi, err = jagex.DecodeFileIndex(bs)
        if err != nil {
            return
        }

        id, found := fi.FindGroupID(args[0])
        if !found {
            fmt.Println("-1")
            return
        }

        fmt.Printf("%d\n", id)
        return
    },
}
