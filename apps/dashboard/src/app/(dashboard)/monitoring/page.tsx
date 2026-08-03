"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useDashboardOverview, useDashboardCapabilities } from "@/services/dashboard";
import { Card } from "@/components/ui/Card";

function formatBytes(value: number) {
  if (value === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(value) / Math.log(1024));
  const size = value / Math.pow(1024, i);
  return `${size.toFixed(1)} ${units[i]}`;
}

export default function MonitoringRoute() {
  const { data, isLoading, isError, refetch } = useDashboardOverview();
  const { data: capabilities } = useDashboardCapabilities();

  const monitoringCapability = capabilities?.capabilities?.find((capability) => capability.capability === "monitoring");

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Monitoring</h1>
          <p className="text-sm text-slate-400">
            Live system telemetry collected by the node agent.
          </p>
        </div>

        {isLoading && (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 4 }).map((_, index) => (
              <div key={index} className="h-32 animate-pulse rounded-lg bg-white/5" />
            ))}
          </div>
        )}

        {isError && (
          <div className="flex items-center justify-between rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            <span>Failed to load monitoring data.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        {data && (
          <>
            <Card title="Processor" description="Current CPU utilization.">
              <p className="text-3xl font-semibold text-white">{data.snapshot.cpu.usage_percent.toFixed(1)}%</p>
            </Card>

            <Card title="Memory" description="RAM utilization and availability.">
              <div className="space-y-2">
                <p className="text-3xl font-semibold text-white">{data.snapshot.memory.usage_percent.toFixed(1)}%</p>
                <p className="text-xs text-slate-400">
                  {formatBytes(data.snapshot.memory.used_bytes)} of {formatBytes(data.snapshot.memory.total_bytes)}
                </p>
              </div>
            </Card>

            <Card title="Disk" description="Primary disk utilization.">
              {data.snapshot.disks.length > 0 ? (
                <div className="space-y-2">
                  <p className="text-3xl font-semibold text-white">{data.snapshot.disks[0].usage_percent.toFixed(1)}%</p>
                  <p className="text-xs text-slate-400">{data.snapshot.disks[0].mountpoint}</p>
                </div>
              ) : (
                <p className="text-sm text-slate-400">No disk data available.</p>
              )}
            </Card>

            <Card title="Network" description="Last known network I/O.">
              <div className="space-y-2">
                <p className="text-3xl font-semibold text-white">{data.snapshot.network.hostname}</p>
                <p className="text-xs text-slate-400">
                  {data.snapshot.network.interfaces.length} interfaces
                </p>
              </div>
            </Card>

            <Card title="Capabilities" description="Monitoring subsystem status.">
              <p className={`text-sm font-medium ${monitoringCapability?.available ? "text-emerald-400" : "text-amber-400"}`}>
                {monitoringCapability?.state ?? "unknown"}
              </p>
            </Card>
          </>
        )}
      </div>
    </DashboardShell>
  );
}
