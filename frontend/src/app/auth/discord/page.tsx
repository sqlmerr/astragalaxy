"use client";

import useLocalStorage from "@/hooks/useLocalStorage";
import { request } from "@/lib/api";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function DiscordAuth() {
  const [isLoading, setIsLoading] = useState(false);
  const params = useSearchParams();
  const code = params.get("code");
  const [token, setToken] = useLocalStorage("access-token", "");

  useEffect(() => {
    async function action() {
      if (code && !isLoading) {
        setIsLoading(true);
        const body = JSON.stringify({
          code: code,
        });

        console.log(code, body);
        const response = await request("/auth/discord", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: body,
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }
        setIsLoading(false);
        const json: { access_token: string } = await response.json();
        setToken(json.access_token);
        return;
      }
    }

    action();
  }, []);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  // redirect("/");
}
