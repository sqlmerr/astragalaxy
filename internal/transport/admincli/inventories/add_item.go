package inventories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/game/items"
)

func newAddItemCmd(itemsService *items.Service) *cobra.Command {
	var metadata string

	cmd := &cobra.Command{
		Use:   "add_item <inventory id> <item id>",
		Short: "Adds item to an inventory",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			inventoryID, err := uuid.Parse(args[0])
			if err != nil {
				return fmt.Errorf("invalid inventory id (%s): %w", args[0], err)
			}
			itemID := args[1]

			var rawMetadata json.RawMessage
			if metadata != "" {
				if !json.Valid([]byte(metadata)) {
					return fmt.Errorf("metadata must be valid JSON")
				}
				rawMetadata = json.RawMessage(metadata)
			}

			item, err := itemsService.GrantItem(ctx, inventoryID, itemID, rawMetadata)
			if err != nil {
				return fmt.Errorf("failed to add item: %w", err)
			}
			fmt.Printf("added item %q to inventory %s\n(id: %s, type: %s)\n", itemID, inventoryID, item.ID, item.ItemType)

			return nil
		},
	}

	cmd.Flags().StringVarP(&metadata, "metadata", "m", "", "item metadata as JSON (optional)")

	return cmd
}
