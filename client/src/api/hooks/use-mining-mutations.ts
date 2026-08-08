import { useMutation } from "@tanstack/react-query"

import { mineAsteroid, minePlanet } from "@/api/mutations/mining"
import type {
  SchemaMineAsteroidRequest,
  SchemaMinePlanetRequest,
} from "@/api/types"

export function useMineAsteroidMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      body,
    }: {
      agentId: string
      body: SchemaMineAsteroidRequest
    }) => mineAsteroid(agentId, body),
  })
}

export function useMinePlanetMutation() {
  return useMutation({
    mutationFn: ({
      agentId,
      body,
    }: {
      agentId: string
      body: SchemaMinePlanetRequest
    }) => minePlanet(agentId, body),
  })
}
