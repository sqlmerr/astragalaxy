import { useState } from "react"
import type { SyntheticEvent } from "react"
import { useNavigate } from "@tanstack/react-router"
import { KeyRound, Rocket, ShieldCheck, UserRound } from "lucide-react"

import { useLoginMutation } from "@/api/hooks"
import { getErrorMessage } from "@/api/errors"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"
import { Textarea } from "@/components/ui/textarea"
import { useAuth } from "@/components/features/auth/auth-provider"
import { toast } from "@/components/ui/toast"

export function LoginForm() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [jwtToken, setJwtToken] = useState("")
  const [jwtError, setJwtError] = useState<string | null>(null)
  const { signIn } = useAuth()
  const navigate = useNavigate()

  const passwordLogin = useLoginMutation()

  function handlePasswordLogin(event: SyntheticEvent) {
    event.preventDefault()
    passwordLogin.mutate(
      { username, password },
      {
        onSuccess: (response) => {
          if (response.access_token) {
            signIn(response.access_token)
            navigate({ to: "/", replace: true })
          }
        },
        onError: (error) => {
          toast.add({
            type: "error",
            title: "Authentication failed",
            description: getErrorMessage(error),
          })
        },
      }
    )
  }

  function handleJwtLogin(event: SyntheticEvent) {
    event.preventDefault()
    setJwtError(null)

    if (!jwtToken.trim()) {
      setJwtError("Please enter a JWT token")
      return
    }

    signIn(jwtToken.trim())
    void navigate({ to: "/", replace: true })
  }

  return (
    <div className="grid w-full max-w-4xl grid-cols-1 gap-4 md:grid-cols-2">
      <Card>
        <CardHeader className="gap-2">
          <div className="flex items-center gap-2 text-primary">
            <UserRound className="size-4" aria-hidden="true" />
            <span className="text-[10px] font-bold tracking-[0.2em] uppercase">
              Command access
            </span>
          </div>
          <CardTitle className="text-base">
            Sign in to your command deck
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form className="grid gap-4" onSubmit={handlePasswordLogin}>
            <label
              className="grid gap-2 text-xs font-semibold tracking-wide text-muted-foreground"
              htmlFor="username"
            >
              Username
              <Input
                id="username"
                autoComplete="username"
                value={username}
                onChange={(event) => setUsername(event.target.value)}
                placeholder="SPACE-RANGER-37"
                required
              />
            </label>
            <label
              className="grid gap-2 text-xs font-semibold tracking-wide text-muted-foreground"
              htmlFor="password"
            >
              Password
              <Input
                id="password"
                type="password"
                autoComplete="current-password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                placeholder="Enter your password"
                required
              />
            </label>
            <Button
              className="mt-1 w-full"
              type="submit"
              disabled={passwordLogin.isPending}
            >
              <Rocket data-icon="inline-start" aria-hidden="true" />
              {passwordLogin.isPending ? "Connecting" : "Sign In"}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card className="surface-panel">
        <CardHeader className="gap-2">
          <div className="flex items-center gap-2 text-primary">
            <KeyRound className="size-4" aria-hidden="true" />
            <span className="text-[10px] font-bold tracking-[0.2em] uppercase">
              Token access
            </span>
          </div>
          <CardTitle className="text-base">Continue with JWT</CardTitle>
        </CardHeader>
        <CardContent>
          <Separator className="mb-5" />
          <form className="grid gap-4" onSubmit={handleJwtLogin}>
            <label
              className="grid gap-2 text-xs font-semibold tracking-wide text-muted-foreground"
              htmlFor="jwt-token"
            >
              JWT Token
              <Textarea
                id="jwt-token"
                value={jwtToken}
                onChange={(event) => {
                  setJwtToken(event.target.value)
                  setJwtError(null)
                }}
                placeholder="Paste your JWT token"
                spellCheck={false}
                required
              />
            </label>
            <Button className="w-full" type="submit" variant="outline">
              <ShieldCheck data-icon="inline-start" aria-hidden="true" />
              Continue with JWT
            </Button>
          </form>
          {jwtError ? (
            <p
              className="mt-4 border-l-2 border-destructive px-3 text-xs leading-relaxed text-destructive"
              role="alert"
            >
              {jwtError}
            </p>
          ) : null}
        </CardContent>
      </Card>
    </div>
  )
}
