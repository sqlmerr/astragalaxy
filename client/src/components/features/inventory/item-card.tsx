import { Package } from "lucide-react"

import type { SchemaItem } from "@/api/types"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { useUseItemMutation } from "@/api/hooks/use-inventory-mutations"
import { useQueryClient } from "@tanstack/react-query"
import { useErrorHandler } from "@/errors/utils"
import { toast } from "@/components/ui/toast"
import { queryKeys } from "@/api/query-keys"
import { useAgents } from "../auth/use-agents"

interface ItemCardProps {
  item: SchemaItem
  inventoryAgentId: string
}

export function ItemCard({ item, inventoryAgentId }: ItemCardProps) {
  const { setCooldown } = useAgents()
  const queryClient = useQueryClient()
  const useItemMutation = useUseItemMutation()
  const errorHandler = useErrorHandler()

  async function useItem() {
    if (!inventoryAgentId) return

    await useItemMutation.mutateAsync(
      {
        agentId: inventoryAgentId,
        itemId: item.id,
      },
      {
        onError(err) {
          errorHandler(err, "Failed to use item")
        },
        onSuccess(data) {
          queryClient.invalidateQueries({ queryKey: queryKeys.inventories.all })
          setCooldown(inventoryAgentId, data.cooldown)
          toast.add({
            type: "success",
            title: "Success",
            description: "Successfully used item",
          })
        },
      }
    )
  }

  return (
    <article className="group border border-transparent p-3 transition-colors hover:border-primary/25 hover:bg-primary/4">
      <div className="flex items-start gap-3">
        <Avatar className="transition-transform duration-300 group-hover:scale-105">
          <AvatarFallback>
            <Package className="size-4" aria-hidden="true" />
          </AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <h4 className="truncate text-sm font-semibold tracking-wide">
            {item.item_type}
          </h4>
          <p className="mt-1 truncate text-[11px] text-muted-foreground">
            ID: {item.id.slice(0, 8)}...
          </p>
        </div>
        <Button
          type="button"
          size="sm"
          onClick={useItem}
          disabled={useItemMutation.isPending}
        >
          {useItemMutation.isPending ? "Using..." : "Use"}
        </Button>
      </div>
    </article>
  )
}
