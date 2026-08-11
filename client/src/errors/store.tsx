import { createContext, useContext, useState  } from "react"
import type {ReactNode} from "react";
import type { AppError } from "./types"

interface ErrorContextValue {
  error: AppError | null
  open: (error: AppError) => void
  close: () => void
}

const ErrorContext = createContext<ErrorContextValue | undefined>(undefined)

export function ErrorProvider({ children }: { children: ReactNode }) {
  const [error, setError] = useState<AppError | null>(null)

  return (
    <ErrorContext.Provider
      value={{ error, open: setError, close: () => setError(null) }}
    >
      {children}
    </ErrorContext.Provider>
  )
}

export function useError() {
  const context = useContext(ErrorContext)

  if (!context) {
    throw new Error("useError must be used within ErrorProvider")
  }

  return context
}
