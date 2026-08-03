"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useObserverStatus, useObserverEnvironment } from "@/services/dashboard";
import type { ObserverStatus, ObserverEnvironment } from "@/types/dashboard";

function formatUptime(uptimeSeconds: number): string {
  const hours = Math.floor(uptimeSeconds / 3600);
  const minutes = Math.floor((uptimeSeconds % 3600) / 60);
  const seconds = uptimeSeconds % 60;

  if (hours > 0) {
    return `${hours}h ${minutes}m ${seconds}s`;
  }
  if (minutes > 0) {
    return `${minutes}m ${seconds}s`;
  }
  return `${seconds}s`;
}

function formatLastSeen(lastSeen: string): string {
  const date = new Date(lastSeen);
  return date.toLocaleString();
}

export default function ObserverRoute() {
  const { data: status, isLoading: statusLoading, isError: statusError } = useObserverStatus();
  const { data: environment, isLoading: envLoading, isError: envError } = useObserverEnvironment();

  const observerStatus = status as ObserverStatus | undefined;
  const observerEnv = environment as ObserverEnvironment | undefined;

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Observer</h1>
          <p className="text-sm text-slate-400">Environmental monitoring and local status display.</p>
        </div>

        {statusLoading && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
            Loading Observer status...
          </div>
        )}

        {statusError && (
          <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            Failed to load Observer status.
          </div>
        )}

        {observerStatus && (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div className="rounded-lg border border-white/10 bg-white/5 p-4">
              <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Status</p>
              <p className="mt-2 flex items-center gap-2 text-sm font-medium text-white">
                <span
                  className={`inline-block h-2 w-2 rounded-full ${
                    observerStatus.status === "online" ? "bg-emerald-400" : "bg-red-400"
                  }`}
                />
                {observerStatus.status === "online" ? "Online" : "Offline"}
              </p>
            </div>

            <div className="rounded-lg border border-white/10 bg-white/5 p-4">
              <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Firmware</p>
              <p className="mt-2 text-sm font-medium text-white">{observerStatus.firmware}</p>
            </div>

            <div className="rounded-lg border border-white/10 bg-white/5 p-4">
              <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Uptime</p>
              <p className="mt-2 text-sm font-medium text-white">{formatUptime(observerStatus.uptime)}</p>
            </div>

            <div className="rounded-lg border border-white/10 bg-white/5 p-4">
              <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Last Seen</p>
              <p className="mt-2 text-sm font-medium text-white">{formatLastSeen(observerStatus.lastSeen)}</p>
            </div>
          </div>
        )}

        {envLoading && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
            Loading environment data...
          </div>
        )}

        {envError && (
          <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            Failed to load environment data.
          </div>
        )}

        {observerEnv && (
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="rounded-lg border border-white/10 bg-white/5 p-4">
              <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Room Temperature</p>
              <p className="mt-2 text-sm font-medium text-white">{observerEnv.temperature.toFixed(1)}°C</p>
            </div>

            <div className="rounded-lg border border-white/10 bg-white/5 p-4">
              <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Humidity</p>
              <p className="mt-2 text-sm font-medium text-white">{observerEnv.humidity.toFixed(1)}%</p>
            </div>
          </div>
        )}
      </div>
    </DashboardShell>
  );
}
