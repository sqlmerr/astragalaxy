import { BookOpen, Factory, Package, Pickaxe, Wrench } from "lucide-react"

import type {
  SchemaFacilityData,
  SchemaItemData,
  SchemaRecipeData,
  SchemaResourceData,
} from "@/api/types"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Spinner } from "@/components/ui/spinner"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useData } from "@/hooks/use-data"

interface CatalogueDialogProps {
  open: boolean
  onClose: () => void
}

interface CatalogueGridProps {
  children: React.ReactNode
  isPending: boolean
  isError: boolean
  isEmpty: boolean
}

function CatalogueGrid({
  children,
  isPending,
  isError,
  isEmpty,
}: CatalogueGridProps) {
  if (isPending) {
    return (
      <div className="flex min-h-48 items-center justify-center">
        <Spinner />
      </div>
    )
  }

  if (isError) {
    return (
      <p className="flex min-h-48 items-center justify-center text-muted-foreground">
        Failed to load catalogue data.
      </p>
    )
  }

  if (isEmpty) {
    return (
      <p className="flex min-h-48 items-center justify-center text-muted-foreground">
        No entries found.
      </p>
    )
  }

  return (
    <div className="grid content-start justify-start gap-3 sm:grid-cols-2 xl:grid-cols-3">
      {children}
    </div>
  )
}

function ItemCatalogueCard({ item }: { item: SchemaItemData }) {
  return (
    <Card size="sm" className="text-left">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 break-all">
          <Package className="size-4 shrink-0 text-primary" />
          {item.id}
        </CardTitle>
      </CardHeader>
      <CardContent className="text-muted-foreground">
        Provides facility:{" "}
        <span className="text-foreground">{item.provides_facility}</span>
      </CardContent>
    </Card>
  )
}

function ResourceCatalogueCard({ resource }: { resource: SchemaResourceData }) {
  return (
    <Card size="sm" className="text-left">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 break-all">
          <Pickaxe className="size-4 shrink-0 text-primary" />
          {resource.id}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-wrap gap-1.5">
        {resource.tags.length > 0 ? (
          resource.tags.map((tag) => <Badge key={tag}>{tag}</Badge>)
        ) : (
          <span className="text-muted-foreground">No tags</span>
        )}
      </CardContent>
    </Card>
  )
}

function RecipeCatalogueCard({ recipe }: { recipe: SchemaRecipeData }) {
  return (
    <Card size="sm" className="text-left">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 break-all">
          <Wrench className="size-4 shrink-0 text-primary" />
          {recipe.id}
        </CardTitle>
        <p className="text-muted-foreground">
          {recipe.required_facility}, tier {recipe.min_tier}, {recipe.duration}s
        </p>
      </CardHeader>
      <CardContent className="grid gap-2 text-xs">
        <div>
          <span className="text-muted-foreground">Inputs: </span>
          {recipe.inputs
            .map((input) => `${input.resource_id} x${input.amount}`)
            .join(", ") || "None"}
        </div>
        <div>
          <span className="text-muted-foreground">Outputs: </span>
          {recipe.outputs
            .map((output) => `${output.resource_id} x${output.amount}`)
            .join(", ") || "None"}
        </div>
      </CardContent>
    </Card>
  )
}

function FacilityCatalogueCard({ facility }: { facility: SchemaFacilityData }) {
  return (
    <Card size="sm" className="text-left">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 break-all">
          <Factory className="size-4 shrink-0 text-primary" />
          {facility.id}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-wrap gap-1.5">
        <Badge>{facility.type}</Badge>
        <Badge variant="outline">Tier {facility.tier}</Badge>
        <Badge variant="outline">x{facility.time_multiplier} time</Badge>
      </CardContent>
    </Card>
  )
}

export function CatalogueDialog({ open, onClose }: CatalogueDialogProps) {
  const { items, resources, recipes, facilities, isPending, isError } =
    useData()

  return (
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && onClose()}>
      <DialogContent className="h-[min(44rem,calc(100vh-2rem))] max-w-[calc(100%-2rem)] grid-rows-[auto_minmax(0,1fr)] gap-5 sm:max-w-4xl">
        <DialogHeader className="pr-10">
          <DialogTitle className="flex items-center gap-2">
            <BookOpen className="size-5 text-primary" />
            Catalogue
          </DialogTitle>
        </DialogHeader>
        <Tabs defaultValue="items" className="min-h-0 w-full">
          <TabsList className="justify-start self-start">
            <TabsTrigger value="items">Items ({items.length})</TabsTrigger>
            <TabsTrigger value="resources">
              Resources ({resources.length})
            </TabsTrigger>
            <TabsTrigger value="recipes">
              Recipes ({recipes.length})
            </TabsTrigger>
            <TabsTrigger value="facilities">
              Facilities ({facilities.length})
            </TabsTrigger>
          </TabsList>
          <div className="min-h-0 w-full flex-1 overflow-y-auto p-1">
            <TabsContent value="items" className="w-full">
              <CatalogueGrid
                isPending={isPending}
                isError={isError}
                isEmpty={!items.length}
              >
                {items.map((item) => (
                  <ItemCatalogueCard key={item.id} item={item} />
                ))}
              </CatalogueGrid>
            </TabsContent>
            <TabsContent value="resources" className="w-full">
              <CatalogueGrid
                isPending={isPending}
                isError={isError}
                isEmpty={!resources.length}
              >
                {resources.map((resource) => (
                  <ResourceCatalogueCard
                    key={resource.id}
                    resource={resource}
                  />
                ))}
              </CatalogueGrid>
            </TabsContent>
            <TabsContent value="recipes" className="w-full">
              <CatalogueGrid
                isPending={isPending}
                isError={isError}
                isEmpty={!recipes.length}
              >
                {recipes.map((recipe) => (
                  <RecipeCatalogueCard key={recipe.id} recipe={recipe} />
                ))}
              </CatalogueGrid>
            </TabsContent>
            <TabsContent value="facilities" className="w-full">
              <CatalogueGrid
                isPending={isPending}
                isError={isError}
                isEmpty={!facilities.length}
              >
                {facilities.map((facility) => (
                  <FacilityCatalogueCard
                    key={facility.id}
                    facility={facility}
                  />
                ))}
              </CatalogueGrid>
            </TabsContent>
          </div>
        </Tabs>
      </DialogContent>
    </Dialog>
  )
}
