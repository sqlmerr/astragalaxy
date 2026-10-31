import { useQuery } from "@tanstack/react-query"
import { myWalletsQueryOptions } from "../queries/wallets"

export function useMyWalletsQuery(agentId: string, enabled: boolean) {
  return useQuery({ ...myWalletsQueryOptions(agentId), enabled })
}
