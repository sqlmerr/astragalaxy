import { useMutation } from "@tanstack/react-query"

import {
  changeActiveShip,
  dockShip,
  orbitShip,
  renameShip,
} from "@/api/mutations/ships"
import type { SchemaRenameMyShipRequest } from "@/api/types"

export function useRenameShipMutation(agentID: string) {
  return useMutation({
    mutationFn: ({
      id,
      body,
    }: {
      id: string
      body: SchemaRenameMyShipRequest
    }) => renameShip(agentID, id, body),
  })
}

export function useChangeActiveShipMutation(agentID: string) {
  return useMutation({
    mutationFn: (id: string) => changeActiveShip(agentID, id),
  })
}

export function useDockShipMutation() {
  return useMutation({
    mutationFn: ({ agentID }: { agentID: string }) => dockShip(agentID),
  })
}

export function useOrbitShipMutation() {
  return useMutation({
    mutationFn: ({ agentID }: { agentID: string }) => orbitShip(agentID),
  })
}
