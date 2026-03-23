package command

import (
	"fmt"
	biz_omiai "omiai-server/internal/biz/omiai"

	"github.com/spf13/cobra"
)

func (s *Script) Migrate() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Auto migrate database schemas",
		Long:  "Run GORM AutoMigrate for all defined models",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting database migration...")
			err := s.db.AutoMigrate(
				&biz_omiai.Client{},
				&biz_omiai.ClientInteraction{},
				&biz_omiai.ClientCoinRecord{},
			)
			if err != nil {
				fmt.Printf("Migration failed: %v\n", err)
				return
			}
			fmt.Println("Database migration completed successfully!")
		},
	}
}
