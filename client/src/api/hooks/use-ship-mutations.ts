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
      id,
      body,
    }: {
      id: string
      body: SchemaRenameMyShipRequest
    }) => renameShip(id, body),
  })
}

export function useChangeActiveShipMutation() {
  return useMutation({
    mutationFn: changeActiveShip,
  })
}

export function useDockShipMutation() {
  return useMutation({
    mutationFn: dockShip,
  })
}

export function useOrbitShipMutation() {
  return useMutation({
    mutationFn: orbitShip,
  })
}
