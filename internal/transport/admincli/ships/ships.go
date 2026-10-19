package shipscli

import (
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/game/ships"
)

func NewCommand(shipsService *ships.Service) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "ships",
		Short: "Ships management module",
	}

	rootCmd.AddCommand(newSpawnShipCmd(*shipsService))

	return rootCmd
}
