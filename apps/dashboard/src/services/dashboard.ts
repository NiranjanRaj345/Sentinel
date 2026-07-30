import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ActivityEvent, CapabilitiesResponse, DashboardOverview, EventsResponse, HistoryResponse, OperationResult, Rule, RulesResponse } from "@/types/dashboard";

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

async function handleUnauthorized() {
  if (typeof window !== "undefined") {
    localStorage.removeItem("sentinel_token");
  }
}

async function fetchOverview(): Promise<DashboardOverview> {
  const response = await fetch(`${API_BASE}/dashboard/overview`, {
    cache: "no-store",
    headers: getAuthHeaders(),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

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
      headers: getAuthHeaders(),
    }
  );

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

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
    headers: getAuthHeaders(),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

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

async function executeOperation(action: string, confirm: boolean): Promise<OperationResult> {
  const response = await fetch(`${API_BASE}/operations`, {
    method: "POST",
    cache: "no-store",
    headers: getAuthHeaders(),
    body: JSON.stringify({ action, confirm }),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Operation failed: ${response.status}`);
  }

  return response.json();
}

export function useExecuteOperation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ action, confirm }: { action: string; confirm: boolean }) =>
      executeOperation(action, confirm),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });
}

async function fetchRecentEvents(): Promise<EventsResponse> {
  const response = await fetch(`${API_BASE}/events/recent`, {
    cache: "no-store",
    headers: getAuthHeaders(),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    throw new Error(`Failed to load events: ${response.status}`);
  }

  return response.json();
}

async function fetchRules(): Promise<RulesResponse> {
  const response = await fetch(`${API_BASE}/rules`, {
    cache: "no-store",
    headers: getAuthHeaders(),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    throw new Error(`Failed to load rules: ${response.status}`);
  }

  return response.json();
}

export function useRules() {
  return useQuery({
    queryKey: ["rules"],
    queryFn: fetchRules,
    staleTime: 30_000,
  });
}

export function useRecentEvents(limit = 100) {
  return useQuery({
    queryKey: ["events", "recent", limit],
    queryFn: fetchRecentEvents,
    staleTime: 30_000,
  });
}
