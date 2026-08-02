import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ActivityEvent, AutomationExecution, AutomationExecutionsResponse, CapabilitiesResponse, DashboardOverview, EventsResponse, HistoryResponse, OperationResult, Rule, RulesResponse, Resource, ResourcesResponse, ServiceItem, ServicesResponse } from "@/types/dashboard";

const API_BASE = (process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080").replace(/\/$/, "");

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

async function handleResponseError(label: string, response: Response): Promise<never> {
  let body = "";
  try {
    body = await response.text();
  } catch {
    body = `<unreadable ${response.status}>`;
  }
  const error = new Error(`${label}: ${response.status} ${response.statusText}${body ? ` - ${body}` : ""}`);
  if (response.status === 401 || response.status === 403) {
    await handleUnauthorized();
  }
  throw error;
}

async function fetchJSON<T>(label: string, url: string): Promise<T> {
  const response = await fetch(url, {
    cache: "no-store",
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    await handleResponseError(label, response);
  }

  return response.json() as Promise<T>;
}

export function apiBase() {
  return API_BASE;
}

async function fetchOverview(): Promise<DashboardOverview> {
  return fetchJSON("overview", `${API_BASE}/dashboard/overview`);
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
  return fetchJSON("history", `${API_BASE}/dashboard/history?period=${encodeURIComponent(period)}`);
}

export function useDashboardHistory(period: string) {
  return useQuery({
    queryKey: ["dashboard", "history", period],
    queryFn: () => fetchHistory(period),
    enabled: Boolean(period),
  });
}

async function fetchCapabilities(): Promise<CapabilitiesResponse> {
  return fetchJSON("capabilities", `${API_BASE}/dashboard/capabilities`);
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

  if (!response.ok) {
    await handleResponseError("operation", response);
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
  return fetchJSON("events", `${API_BASE}/events/recent`);
}

async function fetchRules(): Promise<RulesResponse> {
  return fetchJSON("rules", `${API_BASE}/rules`);
}

export function useRules() {
  return useQuery({
    queryKey: ["rules"],
    queryFn: fetchRules,
    staleTime: 30_000,
  });
}

async function fetchServices(): Promise<ServicesResponse> {
  return fetchJSON("services", `${API_BASE}/services`);
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
            await handleResponseError("service action", res);
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
  return fetchJSON("resources", `${API_BASE}/resources`);
}

export function useResources() {
  return useQuery({
    queryKey: ["resources"],
    queryFn: fetchResources,
    staleTime: 15_000,
  });
}

async function fetchAutomationExecutions(): Promise<AutomationExecutionsResponse> {
  return fetchJSON("automation", `${API_BASE}/automation/executions`);
}

export function useAutomationExecutions() {
  return useQuery({
    queryKey: ["automation", "executions"],
    queryFn: fetchAutomationExecutions,
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
