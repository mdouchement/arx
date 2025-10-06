package archive

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mholt/archives"
	"github.com/spf13/cobra"
)

func ArchiveCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "archive",
		Aliases: []string{"a"},
		Short:   "Create archive",
		Example: "arx a archive.tar.zst folder/",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			format, _, err := archives.Identify(context.Background(), args[0], nil)
			if err != nil {
				return fmt.Errorf("could not identify archive format: %w", err)
			}

			archiver, ok := format.(archives.Archiver)
			if !ok {
				return errors.New("unsupported archive format")
			}

			filenames := map[string]string{}
			for _, arg := range args[1:] {
				filenames[arg] = ""
			}

			files, err := archives.FilesFromDisk(context.Background(), nil, filenames)
			if err != nil {
				return fmt.Errorf("files from disk: %w", err)
			}

			f, err := os.Create(args[0])
			if err != nil {
				return fmt.Errorf("could not create archive: %w", err)
			}
			defer f.Close()

			err = archiver.Archive(context.Background(), f, files)
			if err != nil {
				return fmt.Errorf("could not archive: %w", err)
			}

			return f.Sync()
		},
	}

	return c
}
