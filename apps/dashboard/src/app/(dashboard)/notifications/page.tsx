"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useNotifications } from "@/services/dashboard";
import { RefreshCw } from "lucide-react";

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

function statusColor(status: string) {
  switch (status) {
    case "sent":
      return "bg-emerald-500/10 text-emerald-300";
    case "failed":
      return "bg-red-500/10 text-red-300";
    case "pending":
    default:
      return "bg-amber-500/10 text-amber-300";
  }
}

function formatTime(value: string) {
  const date = new Date(value);
  return date.toLocaleString([], {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function SkeletonRow() {
  return (
    <div className="h-16 animate-pulse rounded-lg bg-white/5" />
  );
}

export default function NotificationsRoute() {
  const { data, isLoading, isError, refetch } = useNotifications();
  const notifications = data?.notifications ?? [];

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h1 className="text-2xl font-semibold text-white">Notifications</h1>
            <p className="text-sm text-slate-400">
              Recent notification delivery attempts across the node.
            </p>
          </div>
          <button
            onClick={() => refetch()}
            className="flex items-center gap-2 rounded-lg bg-white/10 px-3 py-1.5 text-sm font-medium text-white hover:bg-white/20"
          >
            <RefreshCw className="h-4 w-4" />
            Refresh
          </button>
        </div>

        {isLoading && (
          <div className="space-y-3">
            {Array.from({ length: 6 }).map((_, index) => (
              <SkeletonRow key={index} />
            ))}
          </div>
        )}

        {isError && (
          <div className="flex items-center justify-between rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            <span>Failed to load notifications.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        <div className="space-y-3">
          {notifications.map((notification) => (
            <div
              key={notification.id}
              className="flex items-center justify-between rounded-lg border border-white/10 bg-white/5 px-4 py-3"
            >
              <div className="flex-1">
                <p className="text-sm font-medium text-white">{notification.title}</p>
                <p className="mt-1 text-xs text-slate-400">{notification.message}</p>
                <p className="mt-1 text-xs text-slate-500">
                  {formatTime(notification.createdAt)} · {notification.source}
                  {notification.provider ? ` · via ${notification.provider}` : ""}
                </p>
              </div>
              <div className="ml-4 flex items-center gap-2">
                <span className={`rounded-full px-2 py-1 text-xs font-medium ${severityColor(notification.severity)}`}>
                  {notification.severity}
                </span>
                <span className={`rounded-full px-2 py-1 text-xs font-medium ${statusColor(notification.status)}`}>
                  {notification.status}
                </span>
              </div>
            </div>
          ))}

          {!isLoading && notifications.length === 0 && !isError && (
            <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
              No notifications yet.
            </div>
          )}
        </div>
      </div>
    </DashboardShell>
  );
}
