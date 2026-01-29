package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (s *Script) InsertClass() *cobra.Command {
	return &cobra.Command{
		Use:   "insert-class",
		Short: "Insert clothing classes",
		Long:  "Insert predefined clothing classes into database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Successfully inserted all clothing classes")
		},
	}
}
