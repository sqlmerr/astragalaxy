import type { SchemaPlanet, SystemExtended } from "@/api/types"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { SystemPanel } from "./system-panel"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PlanetPanel } from "./planet-panel"

interface PanelProps {
  system: SystemExtended | null
  selectedPlanet?: SchemaPlanet
  onClose: () => void
  onSystemCenterCamera: () => void
  onSystemOpen: () => void
  onPlanetCenterCamera: () => void
  onSelectPlanet: (p: SchemaPlanet | null) => void
}

export function Panel({
  system,
  selectedPlanet,
  onClose,
  onSystemCenterCamera,
  onSystemOpen,
  onPlanetCenterCamera,
  onSelectPlanet,
}: PanelProps) {
  return (
    <aside
      className={`fixed top-0 right-0 z-50 h-screen w-96 transform border-l border-border bg-card/95 backdrop-blur-md transition-transform duration-300 ${
        system ? "translate-x-0" : "translate-x-full"
      }`}
    >
      <ScrollArea className="h-full">
        <Tabs
          value={!selectedPlanet ? "system" : `planet-${selectedPlanet.orbit}`}
          onValueChange={(v: string) => {
            if (v === "system") {
              onSelectPlanet(null)
            } else {
              const pl = system?.system.planets.find(
                (p) => p.orbit === Number(v.split("-")[1])
              )
              if (pl) {
                onSelectPlanet(pl)
              }
            }
          }}
        >
          <ScrollArea className="w-full whitespace-nowrap">
            <TabsList className="w-max">
              <TabsTrigger value="system">System</TabsTrigger>
              {system?.system.planets.map((p) => (
                <TabsTrigger key={p.orbit} value={`planet-${p.orbit}`}>
                  Planet {p.orbit + 1}
                </TabsTrigger>
              ))}
            </TabsList>

            <ScrollBar orientation="horizontal" />
          </ScrollArea>

          <TabsContent value="system">
            {system && (
              <SystemPanel
                system={system}
                onCenterCamera={onSystemCenterCamera}
                onClose={onClose}
                onSystemOpen={onSystemOpen}
              />
            )}
          </TabsContent>
          {system?.system.planets.map((p) => (
            <TabsContent key={p.orbit} value={`planet-${p.orbit}`}>
              <PlanetPanel
                planet={p}
                onCenterCamera={onPlanetCenterCamera}
                onClose={onClose}
              />
            </TabsContent>
          ))}
        </Tabs>
      </ScrollArea>
    </aside>
  )
}
