import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { AgentWithShip } from "@/api/types"
import { useDockShipMutation, useOrbitShipMutation } from "@/api/hooks"
import { useErrorHandler } from "@/errors/utils"
import { toast } from "@/components/ui/toast"
import { useQueryClient } from "@tanstack/react-query"
import { queryKeys } from "@/api/query-keys"

interface AgentActionsMenuProps {
  btnClassName?: string
  agent: AgentWithShip
}

export function AgentActionsMenu({
  btnClassName,
  agent,
}: AgentActionsMenuProps) {
  const queryClient = useQueryClient()
  const handleError = useErrorHandler()
  const dockMutation = useDockShipMutation()
  const orbitMutation = useOrbitShipMutation()

  const actions = [
    {
      groupLabel: "Ship",
      items: [
        {
          label: "Dock",
          visible: agent.ship.status === "ORBIT",
          action: async () => {
            await dockMutation.mutateAsync(
              { agentID: agent.agent.id },
              {
                onError(err) {
                  handleError(err, "Failed to dock ship")
                },
                onSuccess(data) {
                  // TOOD: cooldowns
                  toast.add({
                    type: "success",
                    title: "Success",
                    description: "Sucessfully docked ship",
                  })
                  queryClient.invalidateQueries({
                    queryKey: queryKeys.ships.active(agent.agent.id),
                  })
                },
              }
            )
          },
        },
        {
          label: "Orbit",
          visible: agent.ship.status === "DOCKED",
          action: async () => {
            await orbitMutation.mutateAsync(
              { agentID: agent.agent.id },
              {
                onError(err) {
                  handleError(err, "Failed to orbit ship")
                },
                onSuccess(data) {
                  // TOOD: cooldowns
                  toast.add({
                    type: "success",
                    title: "Success",
                    description: "Sucessfully orbitted ship",
                  })
                  queryClient.invalidateQueries({
                    queryKey: queryKeys.ships.active(agent.agent.id),
                  })
                },
              }
            )
          },
        },
      ],
    },
  ]

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={<Button className={btnClassName} variant={"outline"} />}
      >
        Actions
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        {actions.map((a) => (
          <>
            <DropdownMenuGroup key={a.groupLabel}>
              <DropdownMenuLabel>{a.groupLabel}</DropdownMenuLabel>
              {a.items.map((i) => (
                <DropdownMenuItem
                  key={i.label}
                  onClick={i.action}
                  disabled={!i.visible}
                >
                  {i.label}
                </DropdownMenuItem>
              ))}
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
          </>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
