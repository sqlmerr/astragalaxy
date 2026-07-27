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
import type { SchemaAgent, SchemaUser } from "@/api/types"

const AUTH_TOKEN_KEY = "astragalaxy.auth.token"
const CURRENT_AGENT_KEY = "astragalaxy.auth.agentID"

interface AuthContextValue {
  token: string | null
  user: SchemaUser | null
  agents: SchemaAgent[] | null
  currentAgentID: string | null
  isAuthenticated: boolean
  isReady: boolean
  signIn: (token: string) => void
  signOut: () => void
  setCurrentAgentID: (agentID: string) => void
  currentAgent: SchemaAgent | null
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<SchemaUser | null>(null)
  const [agents, setAgents] = useState<SchemaAgent[] | null>(null)
  const [currentAgentID, setCurrentAgentIDState] = useState<string | null>(
    () => {
      if (typeof window === "undefined") return null
      return window.localStorage.getItem(CURRENT_AGENT_KEY)
    }
  )
  const [isReady, setIsReady] = useState(false)
  const validatingRef = useRef(false)

  const signOut = useCallback(() => {
    window.localStorage.removeItem(AUTH_TOKEN_KEY)
    window.localStorage.removeItem(CURRENT_AGENT_KEY)
    setToken(null)
    setUser(null)
    setAgents(null)
    setCurrentAgentIDState(null)
  }, [])

  const signIn = useCallback((nextToken: string) => {
    window.localStorage.setItem(AUTH_TOKEN_KEY, nextToken)
    setToken(nextToken)
  }, [])

  const setCurrentAgentID = useCallback((agentID: string) => {
    window.localStorage.setItem(CURRENT_AGENT_KEY, agentID)
    setCurrentAgentIDState(agentID)
  }, [])

  useEffect(() => {
    if (!token) {
      setUser(null)
      setAgents(null)
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
    if (!user) return

    api
      .GET("/api/v1/agents/my", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then(({ data, error }) => {
        if (!error && data?.data) {
          setAgents(data.data)

          if (
            currentAgentID &&
            data.data.some((a) => a.id === currentAgentID)
          ) {
            return
          }

          if (data.data.length > 0) {
            setCurrentAgentID(data.data[0].id)
          }
        }
      })
      .catch(() => {})
  }, [user])

  useEffect(() => {
    const stored = localStorage.getItem(AUTH_TOKEN_KEY)
    setToken(stored)
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      user,
      agents,
      currentAgentID,
      isAuthenticated: token !== null && user !== null,
      isReady,
      signIn,
      signOut,
      setCurrentAgentID,
      currentAgent: agents?.find((v) => v.id == currentAgentID) || null,
    }),
    [
      token,
      user,
      agents,
      currentAgentID,
      isReady,
      signIn,
      signOut,
      setCurrentAgentID,
    ]
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
