import {
  useFacilitiesQuery,
  useItemsQuery,
  useRecipesQuery,
  useResourcesQuery,
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

  const isPending =
    isRecipesPending ||
    isItemsPending ||
    isResourcesPending ||
    isFaciltiesPending

  const isError =
    isRecipesError || isItemsError || isResourcesError || isFacilitiesError

  return {
    recipes: recipes?.data || [],
    items: items?.data || [],
    resources: resources?.data || [],
    facilities: facilities?.data || [],
    isPending,
    isError,
  }
}
