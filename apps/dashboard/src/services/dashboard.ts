import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  ActivityEvent,
  AutomationExecution,
  AutomationExecutionsResponse,
  CapabilitiesResponse,
  DashboardOverview,
  EventsResponse,
  HistoryResponse,
  OperationResult,
  Rule,
  RulesResponse,
  Resource,
  ResourcesResponse,
  ServiceItem,
  ServicesResponse,
  GuardianStatus,
  GuardianActionResponse,
  RecoveryExecution,
  RecoveryExecutionsResponse,
  ExecuteRecoveryRequest,
  ExecuteRecoveryResponse,
  ObserverStatus,
  ObserverEnvironment,
  Notification,
  NotificationsResponse,
} from "@/types/dashboard";

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
  try {
    const response = await fetch(url, {
      cache: "no-store",
      headers: getAuthHeaders(),
    });

    if (!response.ok) {
      await handleResponseError(label, response);
    }

    return response.json() as Promise<T>;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`${label}: ${error.message}`);
    }
    throw new Error(`${label}: network request failed`);
  }
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

async function fetchGuardianStatus(): Promise<GuardianStatus> {
  return fetchJSON("guardian status", `${API_BASE}/guardian/status`);
}

export function useGuardianStatus(refreshMs = 5000) {
  return useQuery({
    queryKey: ["guardian", "status"],
    queryFn: fetchGuardianStatus,
    refetchInterval: refreshMs,
    refetchIntervalInBackground: false,
    staleTime: 0,
  });
}

async function sendGuardianPower(action: "press" | "release"): Promise<GuardianActionResponse> {
  const response = await fetch(`${API_BASE}/guardian/power`, {
    method: "POST",
    cache: "no-store",
    headers: getAuthHeaders(),
    body: JSON.stringify({ action }),
  });

  if (!response.ok) {
    await handleResponseError("guardian power", response);
  }

  return response.json() as Promise<GuardianActionResponse>;
}

export function useGuardianPower() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (action: "press" | "release") => sendGuardianPower(action),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["guardian", "status"] });
    },
  });
}

async function sendGuardianReset(action: "press" | "release"): Promise<GuardianActionResponse> {
  const response = await fetch(`${API_BASE}/guardian/reset`, {
    method: "POST",
    cache: "no-store",
    headers: getAuthHeaders(),
    body: JSON.stringify({ action }),
  });

  if (!response.ok) {
    await handleResponseError("guardian reset", response);
  }

  return response.json() as Promise<GuardianActionResponse>;
}

export function useGuardianReset() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (action: "press" | "release") => sendGuardianReset(action),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["guardian", "status"] });
    },
  });
}

async function fetchRecoveryExecutions(): Promise<RecoveryExecutionsResponse> {
  return fetchJSON("recovery", `${API_BASE}/recovery/recent`);
}

export function useRecoveryExecutions(staleMs = 30_000) {
  return useQuery({
    queryKey: ["recovery", "executions"],
    queryFn: fetchRecoveryExecutions,
    staleTime: staleMs,
  });
}

async function executeRecovery(policyId: string, target: string): Promise<ExecuteRecoveryResponse> {
  const response = await fetch(`${API_BASE}/recovery/execute`, {
    method: "POST",
    cache: "no-store",
    headers: getAuthHeaders(),
    body: JSON.stringify({ policyId, target }),
  });

  if (!response.ok) {
    await handleResponseError("recovery execute", response);
  }

  return response.json() as Promise<ExecuteRecoveryResponse>;
}

export function useExecuteRecovery() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ policyId, target }: { policyId: string; target: string }) =>
      executeRecovery(policyId, target),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["recovery", "executions"] });
    },
  });
}

async function fetchObserverStatus(): Promise<ObserverStatus> {
  return fetchJSON("observer status", `${API_BASE}/observer/status`);
}

export function useObserverStatus(refreshMs = 5000) {
  return useQuery({
    queryKey: ["observer", "status"],
    queryFn: fetchObserverStatus,
    refetchInterval: refreshMs,
    refetchIntervalInBackground: false,
    staleTime: 0,
  });
}

async function fetchObserverEnvironment(): Promise<ObserverEnvironment> {
  return fetchJSON("observer environment", `${API_BASE}/observer/environment`);
}

export function useObserverEnvironment(refreshMs = 5000) {
  return useQuery({
    queryKey: ["observer", "environment"],
    queryFn: fetchObserverEnvironment,
    refetchInterval: refreshMs,
    refetchIntervalInBackground: false,
    staleTime: 0,
  });
}

async function fetchNotifications(): Promise<NotificationsResponse> {
  return fetchJSON("notifications", `${API_BASE}/notifications/recent`);
}

export function useNotifications(staleMs = 30_000) {
  return useQuery({
    queryKey: ["notifications"],
    queryFn: fetchNotifications,
    staleTime: staleMs,
  });
}
