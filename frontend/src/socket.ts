"use client";

import { io } from "socket.io-client";

export const socket = io("localhost:8000", {
  transports: ["websocket"],
  withCredentials: true,
  extraHeaders: {
    authorization: "Bearer " + localStorage.getItem("access-token"),
  },
  path: "/websockets/ws",
  upgrade: true,
  transportOptions: {
    websocket: {
      extraHeaders: {
        authorization: "Bearer " + localStorage.getItem("access-token"),
      },
    },
  },
});
