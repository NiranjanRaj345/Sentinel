"use client";

import { useEffect, useRef, useCallback, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { DashboardOverview } from "@/types/dashboard";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

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

    const ws = new WebSocket(`${API_BASE.replace("http", "ws")}/dashboard/stream`);
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

    ws.onclose = () => {
      setStatus("disconnected");
      const delay = RECONNECT_DELAYS[
        Math.min(reconnectIndex.current, RECONNECT_DELAYS.length - 1)
      ];
      reconnectIndex.current++;
      setStatus("reconnecting");
      timeoutRef.current = setTimeout(() => {
        connect();
      }, delay);
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
      wsRef.current?.close();
    };
  }, [connect]);

  return { status, connect };
}
