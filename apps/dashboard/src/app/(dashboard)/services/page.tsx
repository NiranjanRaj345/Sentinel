"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useServices } from "@/services/dashboard";
import { useToastStore } from "@/stores/toast";
import type { ServiceItem } from "@/types/dashboard";

const STATUS_COLORS: Record<string, string> = {
  running: "bg-emerald-500/10 text-emerald-300",
  stopped: "bg-slate-500/10 text-slate-300",
  failed: "bg-red-500/10 text-red-300",
  unknown: "bg-white/5 text-slate-300",
};

export default function ServicesRoute() {
  const { servicesQuery, actionMutation } = useServices();
  const { data, isLoading, isError, refetch } = servicesQuery;
  const services = (data?.services ?? []) as ServiceItem[];
  const addToast = useToastStore((state) => state.addToast);

  React.useEffect(() => {
    if (actionMutation.isError) {
      addToast(
        actionMutation.error instanceof Error
          ? actionMutation.error.message
          : "Service action failed",
        "error"
      );
    } else if (actionMutation.isSuccess) {
      addToast("Service action completed", "success");
    }
  }, [actionMutation.isError, actionMutation.isSuccess, actionMutation.error, addToast]);

  const runAction = (name: string, action: string) => {
    actionMutation.mutate({ action, name });
  };

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Services</h1>
          <p className="text-sm text-slate-400">
            Manage system services on this node.
          </p>
        </div>

        {isLoading && (
          <div className="space-y-3">
            {Array.from({ length: 4 }).map((_, index) => (
              <div key={index} className="h-16 animate-pulse rounded-lg bg-white/5" />
            ))}
          </div>
        )}

        {isError && (
          <div className="flex items-center justify-between rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            <span>Failed to load services.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        <div className="space-y-3">
          {services.map((svc: ServiceItem) => (
            <div
              key={svc.name}
              className="flex items-center justify-between rounded-lg border border-white/10 bg-white/5 px-4 py-3"
            >
              <div>
                <p className="text-sm font-medium text-white">{svc.name}</p>
                <p className="mt-1 text-xs text-slate-400">
                  {svc.message ?? "—"}
                </p>
              </div>

              <div className="flex items-center gap-2">
                <span
                  className={`rounded-full px-2 py-1 text-xs font-medium ${
                    STATUS_COLORS[svc.status] ?? STATUS_COLORS.unknown
                  }`}
                >
                  {svc.status}
                </span>

                <div className="flex gap-2">
                  <button
                    onClick={() => runAction(svc.name, "start")}
                    disabled={actionMutation.isPending}
                    className="rounded-md bg-emerald-500/10 px-2 py-1 text-xs font-medium text-emerald-200 hover:bg-emerald-500/20 disabled:opacity-60"
                  >
                    Start
                  </button>
                  <button
                    onClick={() => runAction(svc.name, "stop")}
                    disabled={actionMutation.isPending}
                    className="rounded-md bg-red-500/10 px-2 py-1 text-xs font-medium text-red-200 hover:bg-red-500/20 disabled:opacity-60"
                  >
                    Stop
                  </button>
                  <button
                    onClick={() => runAction(svc.name, "restart")}
                    disabled={actionMutation.isPending}
                    className="rounded-md bg-white/10 px-2 py-1 text-xs font-medium text-white hover:bg-white/20 disabled:opacity-60"
                  >
                    Restart
                  </button>
                </div>
              </div>
            </div>
          ))}

          {!isLoading && services.length === 0 && (
            <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
              No services found on this node.
            </div>
          )}
        </div>
      </div>
    </DashboardShell>
  );
}
