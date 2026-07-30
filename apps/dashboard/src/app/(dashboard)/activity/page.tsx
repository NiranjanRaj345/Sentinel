"use client";

import React from "react";
import { useRecentEvents } from "@/services/dashboard";
import { Card } from "@/components/ui/Card";
import { RefreshCw } from "lucide-react";

function severityColor(severity: string) {
  switch (severity) {
    case "critical":
      return "text-rose-400";
    case "warning":
      return "text-amber-400";
    case "info":
    default:
      return "text-sky-400";
  }
}

function formatTime(value: string) {
  const date = new Date(value);
  return date.toLocaleString([], {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function SkeletonRow() {
  return (
    <div className="h-12 animate-pulse rounded-lg bg-white/5" />
  );
}

export default function ActivityPage() {
  const { data, isLoading, isError, refetch } = useRecentEvents(100);

  const events = data?.events ?? [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Activity</h1>
          <p className="text-sm text-slate-400">
            Unified event stream across operations, alerts, and system state.
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

      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 6 }).map((_, index) => (
            <SkeletonRow key={index} />
          ))}
        </div>
      ) : isError || events.length === 0 ? (
        <Card
          title={isError ? "Unable to load activity" : "No events yet"}
          description={
            isError
              ? "Something went wrong while fetching the event stream."
              : "Events will appear here once operations, alerts, or system state changes occur."
          }
        >
          {isError && (
            <button
              onClick={() => refetch()}
              className="mt-3 flex items-center gap-2 rounded-lg bg-white/10 px-3 py-1.5 text-sm font-medium text-white hover:bg-white/20"
            >
              <RefreshCw className="h-4 w-4" />
              Retry
            </button>
          )}
        </Card>
      ) : (
        <div className="space-y-3">
          {events.map((event) => (
            <div
              key={event.id}
              className="flex items-start gap-4 rounded-lg border border-white/5 bg-white/5 px-4 py-3"
            >
              <div className="mt-0.5 h-2 w-2 shrink-0 rounded-full bg-emerald-400" />
              <div className="flex-1">
                <div className="flex items-center justify-between gap-4">
                  <p className="text-sm font-medium text-white">{event.title}</p>
                  <span className={`text-xs font-semibold capitalize ${severityColor(event.severity)}`}>
                    {event.severity}
                  </span>
                </div>
                <p className="mt-1 text-xs text-slate-400">{event.message}</p>
                <p className="mt-1 text-xs text-slate-500">{formatTime(event.createdAt)}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
