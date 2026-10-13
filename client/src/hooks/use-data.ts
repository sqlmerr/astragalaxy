import {
  useFacilitiesQuery,
  useItemsQuery,
  useRecipesQuery,
  useResourcesQuery,
  useShipClassesQuery,
} from "@/api/hooks/use-data-queries"

export function useData() {
  const {
    data: recipes,
    isPending: isRecipesPending,
    isError: isRecipesError,
  } = useRecipesQuery()

  const {
    data: items,
    isPending: isItemsPending,
    isError: isItemsError,
  } = useItemsQuery()

  const {
    data: resources,
    isPending: isResourcesPending,
    isError: isResourcesError,
  } = useResourcesQuery()

  const {
    data: facilities,
    isPending: isFaciltiesPending,
    isError: isFacilitiesError,
  } = useFacilitiesQuery()

  const {
    data: shipClasses,
    isPending: isShipClassesPending,
    isError: isShipClassesError,
  } = useShipClassesQuery()

  const isPending =
    isRecipesPending ||
    isItemsPending ||
    isResourcesPending ||
    isFaciltiesPending ||
    isShipClassesPending

  const isError =
    isRecipesError ||
    isItemsError ||
    isResourcesError ||
    isFacilitiesError ||
    isShipClassesError

  return {
    recipes: recipes?.data || [],
    items: items?.data || [],
    resources: resources?.data || [],
    facilities: facilities?.data || [],
    shipClasses: shipClasses?.data || [],
    isPending,
    isError,
  }
}
