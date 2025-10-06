package archive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
	"github.com/spf13/cobra"
)

func ExtractCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "extract",
		Aliases: []string{"e", "x"},
		Short:   "Extract archive",
		Example: "arx e archive.tar.zst [my-folder/]",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(c *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			dst := "."
			if len(args) > 1 {
				dst = args[1]
			}

			format, stream, err := archives.Identify(context.Background(), args[0], f)
			if err != nil {
				return fmt.Errorf("identify: %w", err)
			}

			extractor, ok := format.(archives.Extractor)
			if !ok {
				return errors.New("unsupported archive")
			}

			return extractor.Extract(context.Background(), stream, func(ctx context.Context, info archives.FileInfo) error {
				destination := filepath.Join(dst, info.NameInArchive)

				if info.IsDir() {
					return os.MkdirAll(destination, info.Mode())
				}

				f, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
				if err != nil {
					return err
				}
				defer f.Close()

				r, err := info.Open()
				if err != nil {
					return err
				}
				defer r.Close()

				_, err = io.Copy(f, r)
				if err != nil {
					return err
				}

				return f.Sync()
			})
		},
	}

	return c
}
