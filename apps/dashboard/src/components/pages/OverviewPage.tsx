"use client";

import React from "react";
import { useDashboardOverview, useDashboardCapabilities, useExecuteOperation } from "@/services/dashboard";
import { useDashboardSocket, ConnectionStatus } from "@/hooks/useDashboardSocket";
import { Card } from "@/components/ui/Card";
import { clsx } from "clsx";
import type { CapabilityStatus, DashboardOverview, DashboardStatus } from "@/types/dashboard";
import { MetricCard } from "@/components/pages/MetricCard";
import { Wifi, WifiOff, RefreshCw, Power, RotateCw, Moon } from "lucide-react";
import { ConfirmDialog } from "@/components/ConfirmDialog";

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
  const { data: capabilities } = useDashboardCapabilities();
  const executeMutation = useExecuteOperation();

  const [toast, setToast] = React.useState<string | null>(null);

  const showToast = (message: string) => {
    setToast(message);
    setTimeout(() => setToast(null), 4000);
  };

  const handleExecute = async (action: string) => {
    try {
      const result = await executeMutation.mutateAsync({ action, confirm: true });
      showToast(`${action}: ${result.message}`);
    } catch (error) {
      showToast(error instanceof Error ? error.message : "Operation failed");
    }
  };

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

  const capabilityStatus: Record<string, CapabilityStatus> = {};
  if (capabilities?.capabilities) {
    for (const cap of capabilities.capabilities) {
      capabilityStatus[cap.capability] = cap;
    }
  }

  const capabilityCards = [
    {
      key: "monitoring",
      label: "Monitoring",
      status: capabilityStatus["monitoring"],
    },
    {
      key: "remote_desktop",
      label: "Remote Desktop",
      status: capabilityStatus["remote_desktop"],
    },
    {
      key: "vpn",
      label: "VPN",
      status: capabilityStatus["vpn"],
    },
    {
      key: "guardian",
      label: "Guardian",
      status: capabilityStatus["guardian"],
    },
    {
      key: "observer",
      label: "Observer",
      status: capabilityStatus["observer"],
    },
  ];

  return (
    <div className="space-y-6">
      {toast && (
        <div className="fixed right-4 top-4 z-50 rounded-lg border border-white/10 bg-slate-900 px-4 py-3 text-sm text-white shadow-xl">
          {toast}
        </div>
      )}
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

      <Card title="Quick Actions" description="Certified remote operations. Each requires confirmation.">
        <div className="flex flex-wrap gap-3">
          <QuickActionButton
            label="Restart"
            icon={<RotateCw className="h-4 w-4" />}
            action="restart"
            description="Reboot this node"
            onExecute={handleExecute}
            loading={executeMutation.isPending}
          />
          <QuickActionButton
            label="Shutdown"
            icon={<Power className="h-4 w-4" />}
            action="shutdown"
            description="Power off this node"
            onExecute={handleExecute}
            loading={executeMutation.isPending}
          />
          <QuickActionButton
            label="Sleep"
            icon={<Moon className="h-4 w-4" />}
            action="sleep"
            description="Suspend this node"
            onExecute={handleExecute}
            loading={executeMutation.isPending}
          />
        </div>
      </Card>

      <Card title="Capabilities" description="What this node can currently do.">
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {capabilityCards.map((item) => {
            const status = item.status;
            const isAvailable = status?.available ?? false;
            const state = status?.state ?? "unknown";
            const details = status?.details ?? "";

            const stateColor =
              state === "active" || state === "ready" || state === "connected"
                ? "text-emerald-400"
                : state === "missing" || state === "unavailable"
                  ? "text-rose-400"
                  : "text-amber-400";

            const icon = isAvailable ? "●" : "○";

            return (
              <div
                key={item.key}
                className="rounded-lg border border-white/5 bg-white/5 px-3 py-2"
              >
                <p className="text-sm font-medium text-white">{item.label}</p>
                <p className={`text-xs ${stateColor}`}>
                  {icon} {state}
                </p>
                {details && (
                  <p className="text-xs text-slate-400">{details}</p>
                )}
              </div>
            );
          })}
        </div>
      </Card>

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
                    {alert.value != null ? `${alert.value.toFixed(1)} ${alert.metric}` : alert.metric}
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

function QuickActionButton({
  label,
  icon,
  action,
  description,
  onExecute,
  loading = false,
}: {
  label: string;
  icon: React.ReactNode;
  action: string;
  description: string;
  onExecute: (action: string) => void;
  loading?: boolean;
}) {
  const [open, setOpen] = React.useState(false);

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="flex items-center gap-2 rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-white hover:bg-white/10"
      >
        {icon}
        {label}
      </button>
      <ConfirmDialog
        open={open}
        title={`Confirm ${label}`}
        description={description}
        confirmLabel={label}
        onConfirm={() => onExecute(action)}
        onCancel={() => setOpen(false)}
        loading={loading}
      />
    </>
  );
}
