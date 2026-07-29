import { ClientOnly, useNavigate } from "@tanstack/react-router"
import { LogOut, UserRound } from "lucide-react"

import {
  useActiveShipQuery,
  useNavigateWarpMutation,
  useShipRadarQuery,
} from "@/api/hooks"
import { Button } from "@/components/ui/button"
import { AgentRoster } from "@/components/features/agents/agent-roster"
import { useAuth } from "@/components/features/auth/auth-provider"
import {
  GalaxyMap,
  type GalaxyMapRef,
} from "@/components/features/map/galaxy/galaxy-map"
import type {
  AgentWithShip,
  SchemaAgent,
  SchemaErrorResponse,
  SchemaPlanet,
  SchemaSystem,
  SystemExtended,
} from "@/api/types"
import { useCallback, useEffect, useRef, useState } from "react"
import { Panel } from "@/components/features/map/panel/panel"
import { useAgentsWithShips } from "@/components/features/auth/use-agents-with-ships"
import { toast } from "@/components/ui/toast"
import {
  SystemMap,
  type SystemMapRef,
} from "@/components/features/map/system/system-map"
import { useQueryClient } from "@tanstack/react-query"
import { queryKeys } from "@/api/query-keys"
import { handleError } from "@/components/features/utils"

export function StarMapLayout() {
  const queryClient = useQueryClient()
  const { signOut, user, currentAgentID, isReady } = useAuth()
  const navigate = useNavigate()
  const [selectedSystem, setSelectedSystem] = useState<SystemExtended | null>(
    null
  )
  const [selectedPlanet, setSelectedPlanet] = useState<SchemaPlanet | null>(
    null
  )
  const [openedSystem, setOpenedSystem] = useState<SystemExtended | null>(null)

  const galaxyMapRef = useRef<GalaxyMapRef>(null)
  const systemMapRef = useRef<SystemMapRef>(null)

  const warpMutation = useNavigateWarpMutation()

  function handleSignOut() {
    signOut()
    void navigate({ to: "/login", replace: true })
  }

  function closeSystem() {
    setSelectedSystem(null)
    setSelectedPlanet(null)
    setOpenedSystem(null)
    galaxyMapRef.current?.closeSystem()
  }

  async function warp() {
    if (!selectedSystem || !currentAgentID) {
      return
    }

    await warpMutation.mutateAsync(
      {
        agentId: currentAgentID,
        body: { x: selectedSystem.system.x, y: selectedSystem.system.y },
      },
      {
        onError(err) {
          handleError(err, "Failed to warp")
        },
        onSuccess(data) {
          // TODO: cooldown
          toast.add({
            type: "success",
            title: "Success",
            description: "Successfully warped to another system",
          })
          queryClient.invalidateQueries({
            queryKey: queryKeys.agents.my,
          })
        },
      }
    )
  }

  useEffect(() => {
    closeSystem()
  }, [currentAgentID])

  const {
    data: ship,
    isPending: isShipPending,
    isError: isShipError,
  } = useActiveShipQuery(currentAgentID || undefined)

  const {
    data: agentsWithShips,
    isPending: isAgentsPending,
    isError: isAgentsError,
  } = useAgentsWithShips()

  const {
    data: systemsData,
    isPending: isSystemsPending,
    isError: isSystemsError,
  } = useShipRadarQuery(currentAgentID || undefined)

  const systems: SchemaSystem[] = systemsData?.data ?? []

  const agentLocations: Map<string, AgentWithShip[]> = new Map()
  for (const a of agentsWithShips) {
    const key = `${a.ship.system_x}:${a.ship.system_y}`
    const l = agentLocations.get(key) ?? []
    l.push(a)
    agentLocations.set(key, l)
  }

  const extendedSystems: SystemExtended[] = []
  systems.forEach((s) => {
    const agentsInSystem = agentLocations.get(`${s.x}:${s.y}`) || []

    extendedSystems.push({
      system: s,
      agents: agentsInSystem,
    })
  })

  useEffect(() => {
    if (isSystemsError) {
      toast.add({
        type: "error",
        title: "Error loading radar data",
      })
    }
    if (isShipError) {
      toast.add({
        type: "error",
        title: "Error loading agent ship data",
      })
    }
    if (isAgentsError) {
      toast.add({
        type: "error",
        title: "Error loading agents data",
      })
    }
  }, [isSystemsError, isShipError, isAgentsError, toast])

  const ready =
    isReady &&
    !isSystemsPending &&
    !isShipPending &&
    !isSystemsError &&
    !isShipError &&
    currentAgentID &&
    ship &&
    !isAgentsError &&
    !isAgentsPending

  return (
    <div className="relative h-screen overflow-hidden bg-background">
      <ClientOnly>
        {ready &&
          (!openedSystem ? (
            <GalaxyMap
              ref={galaxyMapRef}
              onSystemClick={setSelectedSystem}
              ship={ship}
              systems={extendedSystems}
            />
          ) : (
            <SystemMap
              ref={systemMapRef}
              system={openedSystem}
              agents={agentsWithShips}
              onPlanetClick={setSelectedPlanet}
            />
          ))}
      </ClientOnly>

      <div
        className="star-field pointer-events-none absolute inset-0 z-10"
        aria-hidden="true"
      />

      <div className="relative z-20">
        <AgentRoster />
        <Panel
          system={selectedSystem}
          currentAgent={agentsWithShips.find(
            (a) => a.agent.id === currentAgentID
          )!}
          selectedPlanet={selectedPlanet || undefined}
          onClose={closeSystem}
          onSystemCenterCamera={() =>
            selectedSystem &&
            galaxyMapRef.current?.centerOnSystem(selectedSystem.system)
          }
          onSystemOpen={() => selectedSystem && setOpenedSystem(selectedSystem)}
          onSelectPlanet={(p) => {
            setSelectedPlanet(p)
            systemMapRef.current?.selectPlanet(p)
          }}
          onSystemWarp={warp}
        />

        <div className="fixed top-4 right-4 flex items-center gap-2 lg:top-6 lg:right-6">
          <div className="hidden items-center gap-2 border border-border bg-card/80 px-3 py-2 backdrop-blur-sm sm:flex">
            <UserRound className="size-3.5 text-primary" />
            <span className="text-xs font-semibold tracking-wide">
              {user?.username ?? "Commander"}
            </span>
          </div>

          <Button variant="outline" size="sm" onClick={handleSignOut}>
            <LogOut />
            <span className="hidden sm:inline">Log out</span>
          </Button>
        </div>
      </div>
    </div>
  )
}
