package coffee

import (
    "bufio"
    "errors"
    "github.com/hadyn/coffee/jagex/sprite"
    "github.com/spf13/cobra"
    "image/gif"
    "image/jpeg"
    "image/png"
    "io/ioutil"
    "os"
    "strconv"
)

var spriteCmd = &cobra.Command{
    Use:   "sprite",
    Short: "Sprite editing commands",
    Long:  ``,

    Run: func(cmd *cobra.Command, args []string) {
        _ = cmd.Help()
    },
}

var (
    spriteDecodeFormat string
)

var spriteDecodeCmd = &cobra.Command{
    Use:   "decode [n]",
    Short: "Decode a sprite group and output a single sprite",
    Long:  ``,
    RunE: func(cmd *cobra.Command, args []string) (err error) {
        n := 0
        if len(args) > 0 {
            n, err = strconv.Atoi(args[0])
            if err != nil {
                return
            }

            if n < 0 {
                return errors.New(
                    "expected the sprite identifier to be greater than or equal to zero",
                )
            }
        }

        in, err := ioutil.ReadAll(os.Stdin)
        if err != nil {
            return
        }

        var (
            group = sprite.DecodeGroup(in)
            child = group.Get(n)
        )

        w := bufio.NewWriter(os.Stdout)
        defer w.Flush()

        switch spriteDecodeFormat {
        case "png":
            err = png.Encode(w, child.ToImage())
        case "jpeg":
            err = jpeg.Encode(w, child.ToImage(), nil)
        case "gif":
            err = gif.Encode(w, child.ToImage(), nil)
        }

        return
    },
}
