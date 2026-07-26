import { useMutation } from "@tanstack/react-query"

import { registerUser } from "@/api/mutations/register"

export function useRegisterMutation() {
  return useMutation({
    mutationFn: registerUser,
  })
}
