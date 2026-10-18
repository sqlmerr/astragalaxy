package inventories

import (
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/game/items"
)

func NewCommand(itemsService *items.Service) *cobra.Command {
	var rootInventoriesCmd = &cobra.Command{
		Use:   "inventories",
		Short: "Inventory management module",
	}

	rootInventoriesCmd.AddCommand(newAddItemCmd(itemsService))

	return rootInventoriesCmd
}
