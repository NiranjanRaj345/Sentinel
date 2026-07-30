import { useQuery } from "@tanstack/react-query";
import type { DashboardOverview } from "@/types/dashboard";

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
