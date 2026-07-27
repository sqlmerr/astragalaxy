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

export function useDockShipMutation(agentID: string) {
  return useMutation({
    mutationFn: () => dockShip(agentID),
  })
}

export function useOrbitShipMutation(agentID: string) {
  return useMutation({
    mutationFn: () => orbitShip(agentID),
  })
}
