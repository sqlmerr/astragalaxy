import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react"
import type { ReactNode } from "react"

import { api } from "@/api/client"
import type { SchemaUser } from "@/api/types"

const AUTH_TOKEN_KEY = "astragalaxy.auth.token"

interface AuthContextValue {
  token: string | null
  user: SchemaUser | null
  isAuthenticated: boolean
  isReady: boolean
  signIn: (token: string) => void
  signOut: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<SchemaUser | null>(null)
  const [isReady, setIsReady] = useState(false)
  const validatingRef = useRef(false)

  const signOut = useCallback(() => {
    window.localStorage.removeItem(AUTH_TOKEN_KEY)
    setToken(null)
    setUser(null)
  }, [])

  const signIn = useCallback((nextToken: string) => {
    window.localStorage.setItem(AUTH_TOKEN_KEY, nextToken)
    setToken(nextToken)
  }, [])

  useEffect(() => {
    if (!token) {
      setUser(null)
      setIsReady(true)
      return
    }
    validatingRef.current = true

    api
      .GET("/api/v1/auth/me", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then(({ data, error }) => {
        if (!validatingRef.current) return
        if (error) {
          signOut()
        } else {
          setUser(data)
        }
      })
      .catch(() => {
        if (!validatingRef.current) return
        signOut()
      })
      .finally(() => {
        if (validatingRef.current) {
          setIsReady(true)
          validatingRef.current = false
        }
      })

    return () => {
      validatingRef.current = false
    }
  }, [token])

  useEffect(() => {
    const stored = localStorage.getItem(AUTH_TOKEN_KEY)
    setToken(stored)
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      user,
      isAuthenticated: token !== null && user !== null,
      isReady,
      signIn,
      signOut,
    }),
    [token, user, isReady, signIn, signOut]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider")
  }

  return context
}
