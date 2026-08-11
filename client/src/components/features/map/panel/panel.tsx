import {
  isSystemExtended
  
  
  
} from "@/api/types"
import type {AnySystemExtended, SchemaPlanet, SchemaWaypoint, AgentExtended, SchemaResourceDeposit } from "@/api/types";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { SystemPanel } from "./system-panel"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PlanetPanel } from "./planet-panel"
import { WaypointPanel } from "./waypoint-panel"
import { WAYPOINT_PARAMS } from "../constants"
import { AsteroidMineDialog } from "./asteroid-mine-dialog"
import { useState } from "react"
import { PlanetMineDialog } from "./planet-mine-dialog"

interface PanelProps {
  system: AnySystemExtended | null
  currentAgent: AgentExtended
  selectedPlanet?: SchemaPlanet
  selectedWaypoint?: SchemaWaypoint
  onClose: () => void
  onSystemCenterCamera: () => void
  onSystemOpen: () => void
  onSelectPlanet: (p: SchemaPlanet) => void
  onSystemWarp: () => void
  onSelectWaypoint: (w: SchemaWaypoint) => void
  onSelectNone: () => void
  onPlanetNavigate: (p: SchemaPlanet) => void
  onWaypointNavigate: (w: SchemaWaypoint) => void
}

export function Panel({
  system,
  currentAgent,
  selectedPlanet,
  selectedWaypoint,
  onClose,
  onSystemCenterCamera,
  onSystemOpen,
  onSelectPlanet,
  onSystemWarp,
  onSelectWaypoint,
  onSelectNone,
  onWaypointNavigate,
  onPlanetNavigate,
}: PanelProps) {
  const fullSystem = isSystemExtended(system)
  const [isAsteroidMineDialogOpened, setIsAsteroidMineDialogOpened] =
    useState(false)

  const [planetMineDialogDeposits, setPlanetMineDialogDeposits] = useState<
    SchemaResourceDeposit[] | null
  >(null)

  return (
    <>
      <aside
        className={`fixed top-0 right-0 z-50 h-screen w-96 transform border-l border-border bg-card/95 backdrop-blur-md transition-transform duration-300 ${
          system ? "translate-x-0" : "translate-x-full"
        }`}
      >
        <ScrollArea className="h-full">
          <Tabs
            value={
              selectedWaypoint
                ? `waypoint-${selectedWaypoint.id}`
                : selectedPlanet
                  ? `planet-${selectedPlanet.orbit}`
                  : "system"
            }
            onValueChange={(v: string) => {
              if (v === "system") {
                onSelectNone()
              } else if (v.startsWith("planet-")) {
                const pl =
                  fullSystem &&
                  system.system.planets.find(
                    (p) => p.orbit === Number(v.split("-")[1])
                  )
                if (pl) {
                  onSelectPlanet(pl)
                }
              } else if (v.startsWith("waypoint-")) {
                const wp =
                  fullSystem &&
                  system.system.waypoints.find(
                    (w) => w.id === Number(v.split("-")[1])
                  )
                if (wp) {
                  onSelectWaypoint(wp)
                }
              }
            }}
          >
            <ScrollArea className="w-full whitespace-nowrap">
              <TabsList className="w-max">
                <TabsTrigger value="system">
                  System: {system?.system.name}
                </TabsTrigger>
                {fullSystem &&
                  system.system.planets.map((p) => (
                    <TabsTrigger key={p.orbit} value={`planet-${p.orbit}`}>
                      {p.name}
                    </TabsTrigger>
                  ))}
                {fullSystem &&
                  system.system.waypoints.map((w) => (
                    <TabsTrigger key={w.id} value={`waypoint-${w.id}`}>
                      {WAYPOINT_PARAMS[w.type].name + ` #${w.id}`}
                    </TabsTrigger>
                  ))}
              </TabsList>

              <ScrollBar orientation="horizontal" />
            </ScrollArea>

            <TabsContent value="system">
              {system && (
                <SystemPanel
                  currentAgent={currentAgent}
                  system={system}
                  onCenterCamera={onSystemCenterCamera}
                  onClose={onClose}
                  onSystemOpen={onSystemOpen}
                  onWarp={onSystemWarp}
                />
              )}
            </TabsContent>
            {fullSystem &&
              system.system.planets.map((p) => (
                <TabsContent key={p.orbit} value={`planet-${p.orbit}`}>
                  <PlanetPanel
                    currentAgent={currentAgent}
                    system={system}
                    planet={p}
                    onClose={onClose}
                    onNavigate={onPlanetNavigate}
                    onPlanetMine={() => setPlanetMineDialogDeposits(p.deposits)}
                  />
                </TabsContent>
              ))}
            {fullSystem &&
              system.system.waypoints.map((w) => (
                <TabsContent key={w.id} value={`waypoint-${w.id}`}>
                  <WaypointPanel
                    currentAgent={currentAgent}
                    system={system}
                    waypoint={w}
                    onClose={onClose}
                    onNavigate={onWaypointNavigate}
                    onAsteroidMine={() => setIsAsteroidMineDialogOpened(true)}
                  />
                </TabsContent>
              ))}
          </Tabs>
        </ScrollArea>
        {currentAgent && (
          <>
            <AsteroidMineDialog
              agentId={currentAgent.agent.id}
              onClose={() => setIsAsteroidMineDialogOpened(false)}
              open={isAsteroidMineDialogOpened}
            />
            <PlanetMineDialog
              agentId={currentAgent.agent.id}
              onClose={() => setPlanetMineDialogDeposits(null)}
              deposits={planetMineDialogDeposits}
            />
          </>
        )}
      </aside>
    </>
  )
}
