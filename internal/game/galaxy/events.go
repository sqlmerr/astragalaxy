package galaxy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"github.com/sqlmerr/astragalaxy/internal/data"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type beaconMetadata struct {
	PayoutAddress string `json:"payout_address"`
}

const payoutAmount = 100 // units

func (s *Service) HandleBeaconAgentWarped(ctx context.Context, ev model.AgentWarpedEvent) error {
	fmt.Println("NEW EVENT")
	waypoints, err := s.store.DynamicWaypoints().GetDynamicWaypointsBySystem(ctx, ev.SystemX, ev.SystemY)
	if err != nil {
		return fmt.Errorf("get waypoints in system: %w", err)
	}
	beacons := lo.Filter(waypoints, func(i model.DynamicWaypoint, _ int) bool {
		return i.Type == model.DynamicWaypointCartographyBeacon
	})

	err = s.store.ExecTx(ctx, func(tx data.Store) error {
		for _, b := range beacons {
			var m beaconMetadata
			if err := json.Unmarshal(b.Metadata, &m); err != nil {
				return fmt.Errorf("unmarshal beacon metadata: %w", err)
			}

			if b.OwnerID == ev.AgentID { // TODO: add anti abuse system
				return nil
			}

			_, err = tx.Wallets().AddFunds(ctx, m.PayoutAddress, payoutAmount, 0)
			if err != nil {
				return fmt.Errorf("add funds for payout: %w", err)
			}
		}

		if ev.IsNewDiscovery { // reward for discovering new system
			fmt.Println("AAAA REWARD")
			_, err = s.packageSender.SendSystemPackage(ctx, tx, ev.AgentID, []model.Item{
				{
					ItemType: model.ItemContract,
					Metadata: json.RawMessage("{\"units\": 100, \"quarks\": 3}"),
				},
			}, nil)
			if err != nil {
				return err
			}
			fmt.Println("AAAA SUCCESS")
		}

		return nil
	})

	return err
}
