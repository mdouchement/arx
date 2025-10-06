package archive

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mholt/archives"
	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	var short bool

	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "List archive content",
		Example: "arx ls archive.tar.zst",
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			format, stream, err := archives.Identify(context.Background(), args[0], f)
			if err != nil {
				return fmt.Errorf("identify: %w", err)
			}

			extractor, ok := format.(archives.Extractor)
			if !ok {
				return errors.New("unsupported archive")
			}

			var total int64

			err = extractor.Extract(context.Background(), stream, func(ctx context.Context, info archives.FileInfo) error {
				total++

				if short {
					fmt.Println(info.NameInArchive)
					return nil
				}

				fmt.Printf("%s\t%d\t%s\t%s\n",
					info.Mode(),
					info.Size(),
					info.ModTime(),
					info.NameInArchive,
				)

				return nil
			})
			if err != nil {
				return fmt.Errorf("list: %w", err)
			}

			fmt.Println("total", total)
			return nil
		},
	}
	c.Flags().BoolVarP(&short, "short", "s", false, "Only print filenames")

	return c
}
