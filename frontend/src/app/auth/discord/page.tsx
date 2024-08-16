"use client";

import { request } from "@/lib/api";
import { useCookies } from "next-client-cookies";
import { redirect, useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function DiscordAuth() {
  const cookies = useCookies();
  const [loaded, setLoaded] = useState(false);
  const params = useSearchParams();
  const code = params.get("code");

  // const state = router.query.state;

  useEffect(() => {
    async function action() {
      if (code && !loaded) {
        const body = JSON.stringify({
          code: code,
          state: "123",
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
        setLoaded(true);
        const json: { access_token: string } = await response.json();
        cookies.set("access-token", json.access_token);
        return;
      }
    }

    action();
  }, []);

  if (code && !loaded) {
    return <div>Loading...</div>;
  }

  redirect("/");
}
