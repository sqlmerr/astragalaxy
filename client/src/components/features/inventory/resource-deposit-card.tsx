import { Pickaxe } from "lucide-react"

import type { SchemaResourceDeposit } from "@/api/types"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"

interface ResourceDepositCard {
  deposit: SchemaResourceDeposit
}

export function ResourceDepositCard({ deposit }: ResourceDepositCard) {
  return (
    <article className="group border border-transparent p-3 transition-colors hover:border-primary/25 hover:bg-primary/4">
      <div className="flex items-start gap-3">
        <Avatar className="transition-transform duration-300 group-hover:scale-105">
          <AvatarFallback>
            <Pickaxe className="size-4" aria-hidden="true" />
          </AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <h4 className="truncate text-sm font-semibold tracking-wide">
            {deposit.resource}
          </h4>
          <p className="mt-1 truncate text-[11px] text-muted-foreground">
            Amount: {deposit.amount} Richness: {deposit.richness}
          </p>
        </div>
      </div>
    </article>
  )
}
