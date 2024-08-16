// "use server";

import User from "./types/user";

export async function request(url: string, options?: RequestInit) {
  return await fetch("http://localhost:8000" + url, options);
}

export async function getMe(token: string) {
  const response = await request("/auth/me", {
    headers: {
      Authorization: "Bearer " + token,
    },
  });

  if (response.status !== 200) {
    return null;
  }

  return (await response.json()) as User;
}
