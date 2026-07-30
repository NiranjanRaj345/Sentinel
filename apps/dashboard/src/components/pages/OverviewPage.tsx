"use client";

import React from "react";
import { useDashboardOverview } from "@/services/dashboard";
import { useDashboardSocket, ConnectionStatus } from "@/hooks/useDashboardSocket";
import { Card } from "@/components/ui/Card";
import { clsx } from "clsx";
import type { DashboardOverview, DashboardStatus } from "@/types/dashboard";
import { MetricCard } from "@/components/pages/MetricCard";
import { Wifi, WifiOff, RefreshCw } from "lucide-react";

const statusConfig: Record<
  DashboardStatus,
  { label: string; color: string }
> = {
  healthy: { label: "Healthy", color: "text-emerald-400" },
  warning: { label: "Warning", color: "text-amber-400" },
  critical: { label: "Critical", color: "text-rose-400" },
  offline: { label: "Offline", color: "text-slate-400" },
};

const connectionConfig: Record<
  ConnectionStatus,
  { label: string; color: string; icon: React.ReactNode }
> = {
  live: { label: "Live", color: "text-emerald-400", icon: <Wifi className="h-4 w-4" /> },
  connecting: {
    label: "Connecting",
    color: "text-amber-400",
    icon: <RefreshCw className="h-4 w-4 animate-spin" />,
  },
  reconnecting: {
    label: "Reconnecting",
    color: "text-amber-400",
    icon: <RefreshCw className="h-4 w-4 animate-spin" />,
  },
  disconnected: {
    label: "Disconnected",
    color: "text-rose-400",
    icon: <WifiOff className="h-4 w-4" />,
  },
};

function formatUptime(uptimeMs: number) {
  const seconds = Math.floor(uptimeMs / 1000);
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  const parts = [];
  if (days > 0) parts.push(`${days}d`);
  if (hours > 0) parts.push(`${hours}h`);
  if (minutes > 0 || parts.length === 0) parts.push(`${minutes}m`);

  return parts.join(" ");
}

function formatBytes(value: number) {
  if (value === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(value) / Math.log(1024));
  const size = value / Math.pow(1024, i);
  return `${size.toFixed(1)} ${units[i]}`;
}

export function OverviewPage() {
  const { data, isLoading, isError } = useDashboardOverview();
  const { status } = useDashboardSocket();

  if (isLoading) {
    return (
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-white">Overview</h1>
        <p className="text-sm text-slate-400">Loading dashboard...</p>
      </div>
    );
  }

  if (isError || !data) {
    return (
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-white">Overview</h1>
        <p className="text-sm text-rose-400">
          Failed to load dashboard overview. Is the backend reachable?
        </p>
      </div>
    );
  }

  const nodeStatus = statusConfig[data.status] ?? statusConfig.offline;
  const connection = connectionConfig[status];

  return (
    <div className="space-y-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold text-white">Overview</h1>
        <p className="text-sm text-slate-400">
          {data.nodeName} · Sentinel v{data.version}
        </p>
      </div>

      <Card
        title="Node Status"
        description="Primary health signal for the mission dashboard."
      >
        <div className="grid gap-4 md:grid-cols-3">
          <div>
            <p className="text-xs text-slate-400">Status</p>
            <p className={`text-lg font-semibold ${nodeStatus.color}`}>
              {nodeStatus.label}
            </p>
          </div>
          <div>
            <p className="text-xs text-slate-400">Connection</p>
            <p className={`text-lg font-semibold ${connection.color}`}>
              {connection.label}
            </p>
          </div>
          <div>
            <p className="text-xs text-slate-400">Uptime</p>
            <p className="text-lg font-semibold text-white">
              {formatUptime(data.uptimeMs)}
            </p>
          </div>
          <div>
            <p className="text-xs text-slate-400">Last Collection</p>
            <p className="text-lg font-semibold text-white">
              {new Date(data.lastCollection).toLocaleTimeString()}
            </p>
          </div>
        </div>
      </Card>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          title="CPU"
          value={`${data.snapshot.cpu.usage_percent.toFixed(1)}%`}
          description="Current utilization"
        />
        <MetricCard
          title="Memory"
          value={`${data.snapshot.memory.usage_percent.toFixed(1)}%`}
          description={`${formatBytes(data.snapshot.memory.used_bytes)} of ${formatBytes(data.snapshot.memory.total_bytes)}`}
        />
        <MetricCard
          title="Disk"
          value={
            data.snapshot.disks.length > 0
              ? `${data.snapshot.disks[0].usage_percent.toFixed(1)}%`
              : "n/a"
          }
          description={
            data.snapshot.disks.length > 0
              ? data.snapshot.disks[0].mountpoint
              : "No disks detected"
          }
        />
        <MetricCard
          title="Network"
          value={`${(data.snapshot.network.io.bytes_sent / 1024 / 1024).toFixed(1)} MB/s`}
          description="Last known throughput"
        />
      </div>

      <Card title="Active Alerts" description="Rules currently triggered.">
        {data.activeAlerts.length === 0 ? (
          <p className="text-sm text-emerald-400">No active alerts.</p>
        ) : (
          <div className="space-y-2">
            {data.activeAlerts.map((alert) => (
              <div
                key={alert.ruleId}
                className="flex items-center justify-between rounded-lg border border-white/5 bg-white/5 px-3 py-2"
              >
                <div>
                  <p className="text-sm font-medium text-white">
                    {alert.ruleName}
                  </p>
                  <p className="text-xs text-slate-400">
                    {alert.value.toFixed(1)} {alert.metric}
                  </p>
                </div>
                <span
                  className={clsx(
                    "text-xs font-semibold capitalize",
                    alert.severity === "critical"
                      ? "text-rose-400"
                      : alert.severity === "warning"
                        ? "text-amber-400"
                        : "text-sky-400"
                  )}
                >
                  {alert.severity}
                </span>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
