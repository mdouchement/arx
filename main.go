package main

import (
	"fmt"
	"os"

	"github.com/mdouchement/arx/archive"
	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{
		Use:     "arx",
		Short:   "Archiver eXtended",
		Version: Version(),
		Args:    cobra.NoArgs,
	}
	c.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Version for arx",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(c.Version)
		},
	})
	c.AddCommand(archive.ArchiveCommand())
	c.AddCommand(archive.ExtractCommand())
	c.AddCommand(archive.ListCommand())

	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
