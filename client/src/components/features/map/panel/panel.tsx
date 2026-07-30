import type { SchemaPlanet, SchemaWaypoint, SystemExtended } from "@/api/types"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { SystemPanel } from "./system-panel"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PlanetPanel } from "./planet-panel"
import type { AgentWithShip } from "@/api/types"
import { WaypointPanel } from "./waypoint-panel"
import { WAYPOINT_PARAMS } from "../constants"

interface PanelProps {
  system: SystemExtended | null
  currentAgent: AgentWithShip
  selectedPlanet?: SchemaPlanet
  selectedWaypoint?: SchemaWaypoint
  onClose: () => void
  onSystemCenterCamera: () => void
  onSystemOpen: () => void
  onSelectPlanet: (p: SchemaPlanet) => void
  onSystemWarp: () => void
  onSelectWaypoint: (w: SchemaWaypoint) => void
  onSelectNone: () => void
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
}: PanelProps) {
  return (
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
              const pl = system?.system.planets.find(
                (p) => p.orbit === Number(v.split("-")[1])
              )
              if (pl) {
                onSelectPlanet(pl)
              }
            } else if (v.startsWith("waypoint-")) {
              const wp = system?.system.waypoints.find(
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
              {system?.system.planets.map((p) => (
                <TabsTrigger key={p.orbit} value={`planet-${p.orbit}`}>
                  {p.name}
                </TabsTrigger>
              ))}
              {system?.system.waypoints.map((w) => (
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
          {system?.system.planets.map((p) => (
            <TabsContent key={p.orbit} value={`planet-${p.orbit}`}>
              <PlanetPanel
                currentAgent={currentAgent}
                system={system}
                planet={p}
                onClose={onClose}
              />
            </TabsContent>
          ))}
          {system?.system.waypoints.map((w) => (
            <TabsContent key={w.id} value={`waypoint-${w.id}`}>
              <WaypointPanel
                currentAgent={currentAgent}
                system={system}
                waypoint={w}
                onClose={onClose}
              />
            </TabsContent>
          ))}
        </Tabs>
      </ScrollArea>
    </aside>
  )
}
