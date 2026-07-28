import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react"
import type { ReactNode } from "react"

import type { SchemaAgent, SchemaUser } from "@/api/types"
import { useMeQuery, useMyAgentsQuery } from "@/api/hooks"

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
  const [currentAgentID, setCurrentAgentIDState] = useState<string | null>(
    () => {
      if (typeof window === "undefined") return null
      return window.localStorage.getItem(CURRENT_AGENT_KEY)
    }
  )

  const signOut = useCallback(() => {
    window.localStorage.removeItem(AUTH_TOKEN_KEY)
    window.localStorage.removeItem(CURRENT_AGENT_KEY)
    setToken(null)
    // setAgents(null)
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
    const stored = localStorage.getItem(AUTH_TOKEN_KEY)
    setToken(stored)
  }, [])
  const {
    data: userData,
    isPending: isMePending,
    isError: isMeError,
  } = useMeQuery(!!token)
  const {
    data: agents = { data: [] },
    isPending: isAgentsPending,
    isError: isAgentsError,
  } = useMyAgentsQuery(!!userData)

  useEffect(() => {
    if (isMeError) {
      signOut()
      return
    }
  }, [isMeError, isMePending, userData, signOut])

  useEffect(() => {
    if (isAgentsPending || isAgentsError) return

    if (!currentAgentID || !agents.data.some((a) => a.id === currentAgentID)) {
      if (agents.data.length > 0) {
        setCurrentAgentID(agents.data[0].id)
      }
    }
  }, [
    isAgentsPending,
    isAgentsError,
    agents.data,
    currentAgentID,
    setCurrentAgentID,
  ])

  // if (isMeError) {
  //   signOut()
  // } else if (!isMeError && !isMePending && userData) {
  //   setIsReady(true)
  // }

  // if (!isAgentsPending && !isAgentsError) {
  //   if (!(currentAgentID && agents.data.some((a) => a.id == currentAgentID)))
  //     if (agents.data.length > 0) {
  //       setCurrentAgentID(agents.data[0].id)
  //     }
  // }

  const isReady = token === null || (!isMePending && !isMeError)

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      user: userData || null,
      agents: agents.data,
      currentAgentID,
      isAuthenticated: token !== null && userData !== null,
      isReady,
      signIn,
      signOut,
      setCurrentAgentID,
      currentAgent: agents.data.find((v) => v.id == currentAgentID) || null,
    }),
    [
      token,
      userData,
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
