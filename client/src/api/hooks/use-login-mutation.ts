import { useMutation } from "@tanstack/react-query"

import { loginWithPassword } from "@/api/mutations/auth"

export function useLoginMutation() {
  return useMutation({
    mutationFn: loginWithPassword,
  })
}
