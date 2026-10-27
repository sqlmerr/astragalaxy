package inventoriescli

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/game/inventory"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func newAddResourceCmd(service *inventory.Service) *cobra.Command {
	var substract bool
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
			if substract {
				amount = -amount
			}

			res, err := service.GrantResource(context.Background(), inventoryID, resourceID, amount)
			if err != nil {
				return err
			}

			if res.ResourceType == "" {
				fmt.Printf("deleted resource %q from inventory %s\n", resourceID, inventoryID)
				return nil
			}

			if !substract {
				fmt.Printf("substracted %d resource of type %q from inventory %s\n", int(math.Abs(float64(amount))), resourceID, inventoryID)
			} else {
				fmt.Printf("added %d resource of type %q to inventory %s\n", int(math.Abs(float64(amount))), resourceID, inventoryID)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&substract, "substract", "s", false, "is true then amount will be lower than 0")

	return cmd
}
