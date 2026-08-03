"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useDashboardOverview } from "@/services/dashboard";

function severityColor(severity: string) {
  switch (severity) {
    case "critical":
      return "bg-red-500/10 text-red-300";
    case "warning":
      return "bg-amber-500/10 text-amber-300";
    case "info":
    default:
      return "bg-sky-500/10 text-sky-300";
  }
}

export default function AlertsRoute() {
  const { data, isLoading, isError, refetch } = useDashboardOverview();
  const alerts = data?.activeAlerts ?? [];

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Alerts</h1>
          <p className="text-sm text-slate-400">
            Active alert rules currently triggered on this node.
          </p>
        </div>

        {isLoading && (
          <div className="space-y-3">
            {Array.from({ length: 4 }).map((_, index) => (
              <div key={index} className="h-20 animate-pulse rounded-lg bg-white/5" />
            ))}
          </div>
        )}

        {isError && (
          <div className="flex items-center justify-between rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            <span>Failed to load alerts.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        <div className="space-y-3">
          {alerts.map((alert) => (
            <div
              key={alert.ruleId}
              className="flex items-center justify-between rounded-lg border border-white/10 bg-white/5 px-4 py-3"
            >
              <div>
                <p className="text-sm font-medium text-white">{alert.ruleName}</p>
                <p className="mt-1 text-xs text-slate-400">
                  {alert.value != null ? `${alert.value.toFixed(1)} ${alert.metric}` : alert.metric}
                </p>
              </div>
              <span className={`rounded-full px-2 py-1 text-xs font-medium ${severityColor(alert.severity)}`}>
                {alert.severity}
              </span>
            </div>
          ))}

          {!isLoading && alerts.length === 0 && !isError && (
            <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
              No active alerts.
            </div>
          )}
        </div>
      </div>
    </DashboardShell>
  );
}
