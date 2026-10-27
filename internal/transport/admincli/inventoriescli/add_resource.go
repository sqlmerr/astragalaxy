package inventoriescli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/game/inventory"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func newAddResourceCmd(service *inventory.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add_resource <inventory id> <resource id> <amount>",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inventoryID, err := uuid.Parse(args[0])
			if err != nil {
				return fmt.Errorf("invalid inventory ID: %w", err)
			}

			resourceID := model.ResourceType(args[1])

			amount, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid amount (must be integer): %w", err)
			}

			res, err := service.GrantResource(context.Background(), inventoryID, resourceID, amount)
			if err != nil {
				return err
			}

			if res.ResourceType == "" {
				fmt.Printf("deleted resource %q from inventory %s\n", resourceID, inventoryID)
				return nil
			}

			fmt.Printf("added %d resource of type %q to inventory %s\n", amount, resourceID, inventoryID)
			return nil
		},
	}

	return cmd
}
