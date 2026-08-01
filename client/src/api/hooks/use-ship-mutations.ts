import { useMutation } from "@tanstack/react-query"

import {
  changeActiveShip,
  dockShip,
  orbitShip,
  renameShip,
} from "@/api/mutations/ships"
import type { SchemaRenameMyShipRequest } from "@/api/types"

export function useRenameShipMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      id,
      body,
    }: {
      agentId: string
      id: string
      body: SchemaRenameMyShipRequest
    }) => renameShip(agentId, id, body),
  })
}

export function useChangeActiveShipMutation() {
  return useMutation({
    mutationFn: ({ agentId, id }: { agentId: string; id: string }) =>
      changeActiveShip(agentId, id),
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
