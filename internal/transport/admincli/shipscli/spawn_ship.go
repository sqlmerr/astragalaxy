package shipscli

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/sqlmerr/astragalaxy/internal/game/ships"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

func newSpawnShipCmd(shipsService ships.Service) *cobra.Command {
	var (
		shipType                string
		shipName                string
		shipModules             []string
		inventoryItemSlots      int
		inventoryResourceVolume int
		shipClass               string
	)

	cmd := &cobra.Command{
		Use:     "spawn <agent id>",
		Args:    cobra.ExactArgs(1),
		Short:   "Spawn new ship for agent",
		Example: "admincli ships spawn_ship 2cc37bd1-a16c-487b-b13a-181e33fc5757 --class S --inv-item-slots 45 --modules portable_smelter,portable_printer --name \"Super mega ship ultra best ship ever\" --type trader",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			agentID, err := uuid.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse agent id UUID: %w", err)
			}

			shipType, err := model.NewShipType(shipType)
			if err != nil {
				return fmt.Errorf("parse ship type: %w", err)
			}

			if inventoryItemSlots < 1 {
				return fmt.Errorf("inv-item-slots value must not be less than 1")
			}

			if inventoryResourceVolume < 1 {
				return fmt.Errorf("inv-resource-volume value must not be less than 1")
			}

			ship, err := shipsService.SpawnShip(ctx, ships.SpawnShipSpec{
				AgentID: agentID,
				Type:    shipType,
				Name:    shipName,
				Modules: lo.Map(shipModules, func(item string, _ int) model.ShipModuleType { return model.ShipModuleType(item) }),
				Inventory: model.Inventory{
					MaxItemSlots:      inventoryItemSlots,
					MaxResourceVolume: inventoryResourceVolume,
				},
				Class: model.ShipClassID(shipClass),
			}, nil)
			if err != nil {
				return fmt.Errorf("create ship: %w", err)
			}

			fmt.Printf("Spawned new ship\nShip ID: %s\nOwner (agent): %s\nName: %s\n", ship.ID, agentID, ship.Name)

			return nil
		},
	}

	cmd.Flags().StringVar(&shipType, "type", "scout", "ship type (scout, trader, miner)")
	cmd.Flags().StringVar(&shipName, "name", "Ship", "ship name")
	cmd.Flags().StringSliceVar(&shipModules, "modules", []string{}, "ship modules separated by comma")
	cmd.Flags().IntVar(&inventoryItemSlots, "inv-item-slots", 15, "Maximum amount of available slots for ship's inventory. Min: 1")
	cmd.Flags().IntVar(&inventoryResourceVolume, "inv-resource-volume", 3000, "Maximum resource volume for ship's inventory. Min: 1")
	cmd.Flags().StringVar(&shipClass, "class", "C", "ship class")
	// TODO: enter ship location?

	return cmd
}
