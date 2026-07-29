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

interface AgentActionsMenuProps {
  btnClassName?: string
  agent: AgentWithShip
}

export function AgentActionsMenu({
  btnClassName,
  agent,
}: AgentActionsMenuProps) {
  const actions = [
    {
      groupLabel: "Ship",
      items: [
        {
          label: "Dock",
          visible: agent.ship.status === "ORBIT",
          action: async () => {},
        },
        {
          label: "Orbit",
          visible: agent.ship.status === "DOCKED",
          action: async () => {},
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
            <DropdownMenuGroup>
              <DropdownMenuLabel>{a.groupLabel}</DropdownMenuLabel>
              {a.items.map((i) => (
                <DropdownMenuItem onClick={i.action} disabled={!i.visible}>
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
