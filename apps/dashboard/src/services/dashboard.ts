import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ActivityEvent, AutomationExecution, AutomationExecutionsResponse, CapabilitiesResponse, DashboardOverview, EventsResponse, HistoryResponse, OperationResult, Rule, RulesResponse, Resource, ResourcesResponse, ServiceItem, ServicesResponse } from "@/types/dashboard";

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

async function fetchServices(): Promise<ServicesResponse> {
  const response = await fetch(`${API_BASE}/services`, {
    cache: "no-store",
    headers: getAuthHeaders(),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    throw new Error(`Failed to load services: ${response.status}`);
  }

  return response.json();
}

export function useServices() {
  const queryClient = useQueryClient();

  return {
    servicesQuery: useQuery({
      queryKey: ["services"],
      queryFn: fetchServices,
      staleTime: 15_000,
    }),
    actionMutation: useMutation({
      mutationFn: ({ action, name }: { action: string; name: string }) =>
        fetch(`${API_BASE}/services`, {
          method: "POST",
          cache: "no-store",
          headers: getAuthHeaders(),
          body: JSON.stringify({ action, name }),
        }).then(async (res) => {
          if (!res.ok) {
            const text = await res.text();
            throw new Error(text || `Service action failed: ${res.status}`);
          }
          return res.json() as Promise<ServiceItem>;
        }),
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ["services"] });
      },
    }),
  };
}

async function fetchResources(): Promise<ResourcesResponse> {
  const response = await fetch(`${API_BASE}/resources`, {
    cache: "no-store",
    headers: getAuthHeaders(),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    throw new Error(`Failed to load resources: ${response.status}`);
  }

  return response.json();
}

export function useResources() {
  return useQuery({
    queryKey: ["resources"],
    queryFn: fetchResources,
    staleTime: 15_000,
  });
}

export function useRecentEvents(limit = 100) {
  return useQuery({
    queryKey: ["events", "recent", limit],
    queryFn: fetchRecentEvents,
    staleTime: 30_000,
  });
}

async function fetchAutomationExecutions(): Promise<AutomationExecutionsResponse> {
  const response = await fetch(`${API_BASE}/automation/executions`, {
    cache: "no-store",
    headers: getAuthHeaders(),
  });

  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    throw new Error(`Failed to load automation executions: ${response.status}`);
  }

  return response.json();
}

export function useAutomationExecutions() {
  return useQuery({
    queryKey: ["automation", "executions"],
    queryFn: fetchAutomationExecutions,
    staleTime: 30_000,
  });
}
