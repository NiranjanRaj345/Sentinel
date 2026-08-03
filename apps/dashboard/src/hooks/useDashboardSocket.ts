"use client";

import { useEffect, useRef, useCallback, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { DashboardOverview } from "@/types/dashboard";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

function getAuthHeaders(): HeadersInit {
  const headers: HeadersInit = {
    Accept: "application/json",
  };

  if (typeof window !== "undefined") {
    const token = localStorage.getItem("sentinel_token");
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
  }

  return headers;
}

export type ConnectionStatus = "connecting" | "live" | "disconnected" | "reconnecting";

const RECONNECT_DELAYS = [1000, 2000, 5000, 10000, 30000];

export function useDashboardSocket() {
  const queryClient = useQueryClient();
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectIndex = useRef(0);
  const timeoutRef = useRef<ReturnType<typeof setTimeout>>();
  const [status, setStatus] = useState<ConnectionStatus>("disconnected");

  const connect = useCallback(() => {
    setStatus("connecting");

    const token = typeof window !== "undefined" ? localStorage.getItem("sentinel_token") : null;
    const wsUrl = new URL(`${API_BASE.replace("http", "ws")}/dashboard/stream`);
    if (token) {
      wsUrl.searchParams.set("token", token);
    }
    const ws = new WebSocket(wsUrl.toString());
    wsRef.current = ws;

    ws.onopen = () => {
      setStatus("live");
      reconnectIndex.current = 0;
    };

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (message.type === "overview" && message.data) {
          queryClient.setQueryData(["dashboard", "overview"], message.data);
        }
      } catch {
        // ignore malformed messages
      }
    };

    ws.onclose = (event) => {
      setStatus("disconnected");
      if (!event.wasClean) {
        const delay = RECONNECT_DELAYS[
          Math.min(reconnectIndex.current, RECONNECT_DELAYS.length - 1)
        ];
        reconnectIndex.current++;
        setStatus("reconnecting");
        timeoutRef.current = setTimeout(() => {
          connect();
        }, delay);
      }
    };

    ws.onerror = () => {
      ws.close();
    };
  }, [queryClient]);

  useEffect(() => {
    connect();

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      reconnectIndex.current = 0;
      wsRef.current?.close();
    };
  }, [connect]);

  return { status, connect };
}
