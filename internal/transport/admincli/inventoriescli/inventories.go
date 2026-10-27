package inventoriescli

import (
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/game/inventory"
	"github.com/sqlmerr/astragalaxy/internal/game/items"
)

func NewCommand(itemsService *items.Service, inventoryService *inventory.Service) *cobra.Command {
	var rootInventoriesCmd = &cobra.Command{
		Use:   "inventories",
		Short: "Inventory management module",
	}

	rootInventoriesCmd.AddCommand(newAddItemCmd(itemsService))
	rootInventoriesCmd.AddCommand(newAddResourceCmd(inventoryService))

	return rootInventoriesCmd
}
