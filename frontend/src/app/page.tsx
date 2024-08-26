"use client";

import { Auth } from "@/components/auth";
import useLocalStorage from "@/hooks/useLocalStorage";
import { getMe } from "@/lib/api";
import dynamic from "next/dynamic";
import { useEffect, useState } from "react";

const Game = dynamic(() => import("@/components/game"), {
  ssr: false,
});

export default function Home() {
  const [loading, setLoading] = useState(true);
  const [unAuthorized, setUnAuthorized] = useState(false);
  const [token, setToken] = useLocalStorage("access-token", "");

  useEffect(() => {
    async function action() {
      if (!token) {
        setUnAuthorized(true);
        setLoading(false);
        return;
      }
      const me = await getMe(token);
      if (!me) {
        setUnAuthorized(true);
        setLoading(false);
        return;
      }

      setLoading(false);
    }
    action();
  }, []);

  if (loading) return <h1>Loading</h1>;
  if (unAuthorized) return <Auth />;

  return <Game />;
}
