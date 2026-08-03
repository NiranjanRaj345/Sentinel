"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useRecoveryExecutions, useExecuteRecovery } from "@/services/dashboard";
import { useToastStore } from "@/stores/toast";
import type { RecoveryExecution } from "@/types/dashboard";

export default function RecoveryRoute() {
  const { data, isLoading, isError, refetch } = useRecoveryExecutions();
  const executeMutation = useExecuteRecovery();
  const executions = (data?.executions ?? []) as RecoveryExecution[];
  const addToast = useToastStore((state) => state.addToast);

  React.useEffect(() => {
    if (executeMutation.isError) {
      addToast(
        executeMutation.error instanceof Error
          ? executeMutation.error.message
          : "Recovery execution failed",
        "error"
      );
    } else if (executeMutation.isSuccess) {
      addToast("Recovery execution started", "success");
    }
  }, [executeMutation.isError, executeMutation.isSuccess, executeMutation.error, addToast]);

  const handleExecute = () => {
    executeMutation.mutate({ policyId: "desktop_recovery", target: "" });
  };

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="flex items-center justify-between space-y-1">
          <div className="space-y-1">
            <h1 className="text-2xl font-semibold text-white">Recovery</h1>
            <p className="text-sm text-slate-400">Recovery orchestration and execution history.</p>
          </div>
          <button
            type="button"
            onClick={handleExecute}
            disabled={executeMutation.isPending}
            className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-white/20 disabled:opacity-60"
          >
            Run Desktop Recovery
          </button>
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
            <span>Failed to load recovery history.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        <div className="space-y-3">
          {executions.map((execution) => (
            <div
              key={execution.id}
              className="flex items-center justify-between rounded-lg border border-white/10 bg-white/5 px-4 py-3"
            >
              <div>
                <p className="text-sm font-medium text-white">{execution.policyId}</p>
                <p className="mt-1 text-xs text-slate-400">
                  Step {execution.current} · Attempt {execution.attempts} · {execution.message}
                </p>
                <p className="mt-1 text-xs text-slate-500">
                  Started: {new Date(execution.startedAt).toLocaleString()}
                </p>
              </div>
              <span
                className={`rounded-full px-2 py-1 text-xs font-medium ${
                  execution.status === "succeeded"
                    ? "bg-emerald-500/10 text-emerald-300"
                    : execution.status === "failed"
                      ? "bg-red-500/10 text-red-300"
                      : "bg-amber-500/10 text-amber-300"
                }`}
              >
                {execution.status}
              </span>
            </div>
          ))}

          {!isLoading && executions.length === 0 && (
            <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
              No recovery executions yet.
            </div>
          )}
        </div>
      </div>
    </DashboardShell>
  );
}
