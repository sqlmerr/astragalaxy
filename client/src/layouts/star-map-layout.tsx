import { ClientOnly, useNavigate } from "@tanstack/react-router"
import { LogOut, UserRound } from "lucide-react"

import {
  useActiveShipQuery,
  useNavigatePlanetMutation,
  useNavigateWarpMutation,
  useNavigateWaypointMutation,
  useCurrentSystemQuery,
  useShipRadarQuery,
} from "@/api/hooks"
import { Button } from "@/components/ui/button"
import { AgentRoster } from "@/components/features/agents/agent-roster"
import { useAuth } from "@/components/features/auth/auth-provider"
import {
  GalaxyMap
  
} from "@/components/features/map/galaxy/galaxy-map"
import type {GalaxyMapRef} from "@/components/features/map/galaxy/galaxy-map";
import {
  
  
  
  isSystemExtended
  
  
  
} from "@/api/types"
import type {SchemaWaypoint, AgentExtended, SchemaPlanet, AnySystemExtended, ShortSystemExtended, SystemExtended} from "@/api/types";
import { useEffect, useRef, useState } from "react"
import { Panel } from "@/components/features/map/panel/panel"
import { useAgents } from "@/components/features/auth/use-agents"
import { toast } from "@/components/ui/toast"
import { SystemMap } from "@/components/features/map/system/system-map"
import { useQueryClient } from "@tanstack/react-query"
import { queryKeys } from "@/api/query-keys"
import { currentSystemQueryOptions } from "@/api/queries/systems"
import { useErrorHandler } from "@/errors/utils"

export function StarMapLayout() {
  const queryClient = useQueryClient()
  const handleError = useErrorHandler()
  const { signOut, user, currentAgentID, isReady } = useAuth()
  const navigate = useNavigate()
  const [selectedSystem, setSelectedSystem] =
    useState<AnySystemExtended | null>(null)
  const [selectedPlanet, setSelectedPlanet] = useState<SchemaPlanet | null>(
    null
  )
  const [selectedWaypoint, setSelectedWaypoint] =
    useState<SchemaWaypoint | null>(null)
  const [openedSystem, setOpenedSystem] = useState<SystemExtended | null>(null)

  const galaxyMapRef = useRef<GalaxyMapRef>(null)

  const warpMutation = useNavigateWarpMutation()
  const waypointNavigateMutation = useNavigateWaypointMutation()
  const planetNavigateMutation = useNavigatePlanetMutation()

  const {
    data: agentsWithShips,
    isPending: isAgentsPending,
    isError: isAgentsError,
    setCooldown,
  } = useAgents()

  function handleSignOut() {
    signOut()
    void navigate({ to: "/login", replace: true })
  }

  function closeSystem() {
    setSelectedSystem(null)
    setSelectedPlanet(null)
    setSelectedWaypoint(null)
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
        async onSuccess(data) {
          toast.add({
            type: "success",
            title: "Success",
            description: "Successfully warped to another system",
          })
          queryClient.invalidateQueries({
            queryKey: queryKeys.ships.my(currentAgentID),
          })
          queryClient.invalidateQueries({
            queryKey: queryKeys.systems.current(currentAgentID),
          })
          queryClient.invalidateQueries({
            queryKey: queryKeys.ships.radar(currentAgentID),
          })
          const currentSystem = await queryClient.fetchQuery(
            currentSystemQueryOptions(currentAgentID)
          )
          setSelectedSystem({
            system: currentSystem,
            agents: selectedSystem.agents,
          })
          setCooldown(currentAgentID, data.cooldown)
        },
      }
    )
  }

  async function waypointNavigate(w: SchemaWaypoint) {
    if (!selectedSystem || !currentAgentID || !selectedWaypoint) {
      return
    }

    await waypointNavigateMutation.mutateAsync(
      {
        agentId: currentAgentID,
        body: { id: w.id },
      },
      {
        onError(err) {
          handleError(err, "Failed to navigate to waypoint")
        },
        onSuccess(data) {
          toast.add({
            type: "success",
            title: "Success",
            description: `Successfully navigated to waypoint with id ${w.id}`,
          })
          queryClient.invalidateQueries({
            queryKey: queryKeys.ships.my(currentAgentID),
          })
          setCooldown(currentAgentID, data.cooldown)
        },
      }
    )
  }

  async function planetNavigate(p: SchemaPlanet) {
    if (!selectedSystem || !currentAgentID || !selectedPlanet) {
      return
    }

    await planetNavigateMutation.mutateAsync(
      {
        agentId: currentAgentID,
        body: { orbit: p.orbit },
      },
      {
        onError(err) {
          handleError(err, "Failed to navigate to planet")
        },
        onSuccess(data) {
          toast.add({
            type: "success",
            title: "Success",
            description: `Successfully navigated to planet with orbit ${p.orbit}`,
          })
          queryClient.invalidateQueries({
            queryKey: queryKeys.ships.my(currentAgentID),
          })
          setCooldown(currentAgentID, data.cooldown)
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
    data: systemsData,
    isPending: isSystemsPending,
    isError: isSystemsError,
  } = useShipRadarQuery(currentAgentID || undefined)

  const {
    data: currentSystem,
    isPending: isCurrentSystemPending,
    isError: isCurrentSystemError,
  } = useCurrentSystemQuery(currentAgentID || undefined)

  const systems = systemsData?.data ?? []

  const agentLocations: Map<string, AgentExtended[]> = new Map()
  for (const a of agentsWithShips) {
    const key = `${a.ship.system_x}:${a.ship.system_y}`
    const l = agentLocations.get(key) ?? []
    l.push(a)
    agentLocations.set(key, l)
  }

  const extendedSystems: ShortSystemExtended[] = []
  systems.forEach((s) => {
    const agentsInSystem = agentLocations.get(`${s.x}:${s.y}`) || []

    extendedSystems.push({
      system: s,
      agents: agentsInSystem,
    })
  })

  const currentExtendedSystem: SystemExtended | null = currentSystem
    ? {
        system: currentSystem,
        agents:
          agentLocations.get(`${currentSystem.x}:${currentSystem.y}`) ?? [],
      }
    : null

  function selectSystem(system: ShortSystemExtended) {
    if (
      currentExtendedSystem &&
      system.system.x === currentExtendedSystem.system.x &&
      system.system.y === currentExtendedSystem.system.y
    ) {
      setSelectedSystem(currentExtendedSystem)
      return
    }

    setSelectedSystem(system)
  }

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
    if (isCurrentSystemError) {
      toast.add({
        type: "error",
        title: "Error loading current system",
      })
    }
  }, [isSystemsError, isShipError, isAgentsError, isCurrentSystemError, toast])

  const ready =
    isReady &&
    !isSystemsPending &&
    !isCurrentSystemPending &&
    !isShipPending &&
    !isSystemsError &&
    !isShipError &&
    !isCurrentSystemError &&
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
              onSystemClick={selectSystem}
              ship={ship}
              systems={extendedSystems}
            />
          ) : (
            <SystemMap
              system={openedSystem}
              agents={agentsWithShips}
              selectedPlanet={selectedPlanet}
              setSelectedPlanet={setSelectedPlanet}
              selectedWaypoint={selectedWaypoint}
              setSelectedWaypoint={setSelectedWaypoint}
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
          selectedWaypoint={selectedWaypoint || undefined}
          onClose={closeSystem}
          onSystemCenterCamera={() =>
            selectedSystem &&
            galaxyMapRef.current?.centerOnSystem(selectedSystem.system)
          }
          onSystemOpen={() =>
            selectedSystem &&
            isSystemExtended(selectedSystem) &&
            setOpenedSystem(selectedSystem)
          }
          onSelectPlanet={(p) => {
            setSelectedPlanet(p)
            setSelectedWaypoint(null)
          }}
          onSystemWarp={warp}
          onSelectWaypoint={(w) => {
            setSelectedWaypoint(w)
            setSelectedPlanet(null)
          }}
          onSelectNone={() => {
            setSelectedPlanet(null)
            setSelectedWaypoint(null)
          }}
          onWaypointNavigate={waypointNavigate}
          onPlanetNavigate={planetNavigate}
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
