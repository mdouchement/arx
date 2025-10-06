package archive

import (
	"archive/tar"
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/bodgit/sevenzip"
	"github.com/mholt/archives"
	"github.com/nwaples/rardecode/v2"
	"github.com/spf13/cobra"
)

const layout = "2006-01-02 15:04:05 -0700 MST"

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
			padding := 8

			err = extractor.Extract(context.Background(), stream, func(_ context.Context, info archives.FileInfo) error {
				total++

				if short {
					fmt.Println(info.NameInArchive)
					return nil
				}

				padding = max(padding, len(strconv.FormatInt(info.Size(), 10)))
				template := fmt.Sprintf("%%%dd", padding)
				fsize := fmt.Sprintf(template, info.Size())

				switch h := info.Header.(type) {
				case zip.FileHeader:
					fmt.Printf("%s\t%d\t%s\t%s\t%s\n",
						info.Mode(),
						h.Method,
						fsize,
						info.ModTime().Format(layout),
						h.Name,
					)
				case *tar.Header:
					fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
						info.Mode(),
						h.Uname,
						h.Gname,
						fsize,
						info.ModTime().Format(layout),
						h.Name,
					)
				case *rardecode.FileHeader:
					fmt.Printf("%s\t%d\t%s\t%s\t%s\n",
						info.Mode(),
						int(h.HostOS),
						fsize,
						info.ModTime().Format(layout),
						h.Name,
					)
				case sevenzip.FileHeader:
					fmt.Printf("%s\t%s\t%s\t%s\n",
						info.Mode(),
						fsize,
						info.ModTime().Format(layout),
						f.Name(),
					)
				default:
					fmt.Printf("%s\t%s\t%s\t?/%s\n",
						info.Mode(),
						fsize,
						info.ModTime().Format(layout),
						f.Name(), // we don't know full path from this
					)
				}

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
