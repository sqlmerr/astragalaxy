import { useRecipesQuery } from "@/api/hooks/use-data-queries"

export function useData() {
  const {
    data: recipes,
    isPending: isRecipesPending,
    isError: isRecipesError,
  } = useRecipesQuery()

  return {
    recipes: recipes?.data || [],
    isPending: isRecipesPending,
    isError: isRecipesError,
  }
}
