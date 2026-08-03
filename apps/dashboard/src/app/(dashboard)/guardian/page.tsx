"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useGuardianStatus, useGuardianPower, useGuardianReset } from "@/services/dashboard";
import { useToastStore } from "@/stores/toast";
import type { GuardianStatus } from "@/types/dashboard";

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

export default function GuardianRoute() {
  const { data: status, isLoading, isError, refetch } = useGuardianStatus();
  const powerMutation = useGuardianPower();
  const resetMutation = useGuardianReset();
  const addToast = useToastStore((state) => state.addToast);

  const guardianStatus = status as GuardianStatus | undefined;

  React.useEffect(() => {
    if (powerMutation.isError) {
      addToast(
        powerMutation.error instanceof Error ? powerMutation.error.message : "Power action failed",
        "error"
      );
    } else if (powerMutation.isSuccess) {
      addToast("Power action sent", "success");
    }
  }, [powerMutation.isError, powerMutation.isSuccess, powerMutation.error, addToast]);

  React.useEffect(() => {
    if (resetMutation.isError) {
      addToast(
        resetMutation.error instanceof Error ? resetMutation.error.message : "Reset action failed",
        "error"
      );
    } else if (resetMutation.isSuccess) {
      addToast("Reset action sent", "success");
    }
  }, [resetMutation.isError, resetMutation.isSuccess, resetMutation.error, addToast]);

  const handlePower = (action: "press" | "release") => {
    powerMutation.mutate(action);
  };

  const handleReset = (action: "press" | "release") => {
    resetMutation.mutate(action);
  };

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Guardian</h1>
          <p className="text-sm text-slate-400">Hardware recovery controller</p>
        </div>

        {isLoading && (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 4 }).map((_, index) => (
              <div key={index} className="h-24 animate-pulse rounded-lg bg-white/5" />
            ))}
          </div>
        )}

        {isError && (
          <div className="flex items-center justify-between rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            <span>Failed to load Guardian status.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        {guardianStatus && (
          <>
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div className="rounded-lg border border-white/10 bg-white/5 p-4">
                <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Status</p>
                <p className="mt-2 flex items-center gap-2 text-sm font-medium text-white">
                  <span
                    className={`inline-block h-2 w-2 rounded-full ${
                      guardianStatus.status === "online" ? "bg-emerald-400" : "bg-red-400"
                    }`}
                  />
                  {guardianStatus.status === "online" ? "Online" : "Offline"}
                </p>
              </div>

              <div className="rounded-lg border border-white/10 bg-white/5 p-4">
                <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Firmware</p>
                <p className="mt-2 text-sm font-medium text-white">{guardianStatus.firmware}</p>
              </div>

              <div className="rounded-lg border border-white/10 bg-white/5 p-4">
                <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Uptime</p>
                <p className="mt-2 text-sm font-medium text-white">{formatUptime(guardianStatus.uptime)}</p>
              </div>

              <div className="rounded-lg border border-white/10 bg-white/5 p-4">
                <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Last Seen</p>
                <p className="mt-2 text-sm font-medium text-white">{formatLastSeen(guardianStatus.lastSeen)}</p>
              </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="rounded-lg border border-white/10 bg-white/5 p-4">
                <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Power Relay</p>
                <p className="mt-2 text-xs text-slate-300">
                  Button state: {guardianStatus.powerButton ? "pressed" : "released"}
                </p>
                <div className="mt-3 flex gap-2">
                  <button
                    type="button"
                    onClick={() => handlePower("press")}
                    disabled={powerMutation.isPending}
                    className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-white/20 disabled:opacity-60"
                  >
                    Press
                  </button>
                  <button
                    type="button"
                    onClick={() => handlePower("release")}
                    disabled={powerMutation.isPending}
                    className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-white/20 disabled:opacity-60"
                  >
                    Release
                  </button>
                </div>
              </div>

              <div className="rounded-lg border border-white/10 bg-white/5 p-4">
                <p className="text-xs font-medium uppercase tracking-wide text-slate-400">Reset Relay</p>
                <p className="mt-2 text-xs text-slate-300">
                  Button state: {guardianStatus.resetButton ? "pressed" : "released"}
                </p>
                <div className="mt-3 flex gap-2">
                  <button
                    type="button"
                    onClick={() => handleReset("press")}
                    disabled={resetMutation.isPending}
                    className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-white/20 disabled:opacity-60"
                  >
                    Press
                  </button>
                  <button
                    type="button"
                    onClick={() => handleReset("release")}
                    disabled={resetMutation.isPending}
                    className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-white/20 disabled:opacity-60"
                  >
                    Release
                  </button>
                </div>
              </div>
            </div>
          </>
        )}
      </div>
    </DashboardShell>
  );
}
