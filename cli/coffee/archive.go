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
	Short: "Archive editing commands",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var archiveDecompressCmd = &cobra.Command{
	Use:   "decompress",
	Short: "Decompresses an archive",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			in  []byte
			out []byte
		)

		in, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return
		}

		out, err = jagex.DecompressFileArchive(in)
		if err != nil {
			return
		}

		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		_, err = w.Write(out)
		return
	},
}
