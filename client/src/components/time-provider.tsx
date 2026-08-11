import {
  createContext,
  useContext,
  useEffect,
  useState
  
} from "react"
import type {PropsWithChildren} from "react";

const TimeContext = createContext(Date.now())

export function TimeProvider({ children }: PropsWithChildren) {
  const [now, setNow] = useState(Date.now())

  useEffect(() => {
    const id = setInterval(() => {
      setNow(Date.now())
    }, 1000)

    return () => clearInterval(id)
  }, [])

  return <TimeContext.Provider value={now}>{children}</TimeContext.Provider>
}

export const useNow = () => useContext(TimeContext)
