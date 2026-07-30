import { useQuery } from "@tanstack/react-query";
import type { CapabilitiesResponse, DashboardOverview, HistoryResponse } from "@/types/dashboard";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

async function fetchOverview(): Promise<DashboardOverview> {
  const response = await fetch(`${API_BASE}/dashboard/overview`, {
    cache: "no-store",
    headers: {
      Accept: "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to load dashboard overview: ${response.status}`);
  }

  return response.json();
}

export function useDashboardOverview(refreshMs = 2000) {
  return useQuery({
    queryKey: ["dashboard", "overview"],
    queryFn: fetchOverview,
    refetchInterval: refreshMs,
    refetchIntervalInBackground: false,
    staleTime: 0,
  });
}

async function fetchHistory(period: string): Promise<HistoryResponse> {
  const response = await fetch(
    `${API_BASE}/dashboard/history?period=${encodeURIComponent(period)}`,
    {
      cache: "no-store",
      headers: {
        Accept: "application/json",
      },
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to load history: ${response.status}`);
  }

  return response.json();
}

export function useDashboardHistory(period: string) {
  return useQuery({
    queryKey: ["dashboard", "history", period],
    queryFn: () => fetchHistory(period),
    enabled: Boolean(period),
  });
}

async function fetchCapabilities(): Promise<CapabilitiesResponse> {
  const response = await fetch(`${API_BASE}/dashboard/capabilities`, {
    cache: "no-store",
    headers: {
      Accept: "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to load capabilities: ${response.status}`);
  }

  return response.json();
}

export function useDashboardCapabilities() {
  return useQuery({
    queryKey: ["dashboard", "capabilities"],
    queryFn: fetchCapabilities,
    staleTime: 60_000,
  });
}
