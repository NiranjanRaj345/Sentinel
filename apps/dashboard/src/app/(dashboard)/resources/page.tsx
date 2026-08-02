"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useResources } from "@/services/dashboard";

const HEALTH_COLORS: Record<string, string> = {
  healthy: "bg-emerald-500/10 text-emerald-300",
  degraded: "bg-amber-500/10 text-amber-300",
  unavailable: "bg-red-500/10 text-red-300",
  unknown: "bg-white/5 text-slate-300",
};

const TYPE_LABELS: Record<string, string> = {
  remote_desktop: "Remote Desktop",
  vpn: "VPN",
  container_runtime: "Containers",
  media_server: "Media",
  database: "Database",
  monitoring: "Monitoring",
  application: "Application",
};

export default function ResourcesRoute() {
  const { data, isLoading, isError } = useResources();
  const resources = data?.resources ?? [];

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Resources</h1>
          <p className="text-sm text-slate-400">
            Infrastructure health from the user&apos;s perspective.
          </p>
        </div>

        {isLoading && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
            Loading resources...
          </div>
        )}

        {isError && (
          <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            Failed to load resources.
          </div>
        )}

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {resources.map((resource) => (
            <div
              key={resource.name}
              className="rounded-lg border border-white/10 bg-white/5 p-4"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-semibold text-white">{resource.name}</p>
                  <p className="mt-1 text-xs text-slate-400">
                    {TYPE_LABELS[resource.type] ?? resource.type}
                  </p>
                </div>
                <span
                  className={`rounded-full px-2 py-1 text-xs font-medium ${
                    HEALTH_COLORS[resource.health] ?? HEALTH_COLORS.unknown
                  }`}
                >
                  {resource.health}
                </span>
              </div>

              <div className="mt-4 space-y-2">
                <div className="flex items-center justify-between text-xs">
                  <span className="text-slate-400">Status</span>
                  <span className="text-slate-200">{resource.status}</span>
                </div>
                {resource.message && (
                  <div className="text-xs text-slate-400">
                    {resource.message}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>

        {!isLoading && resources.length === 0 && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
            No resources discovered on this node.
          </div>
        )}
      </div>
    </DashboardShell>
  );
}
