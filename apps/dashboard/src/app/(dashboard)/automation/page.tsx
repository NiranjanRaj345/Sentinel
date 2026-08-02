"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useAutomationExecutions } from "@/services/dashboard";

export default function AutomationRoute() {
  const { data, isLoading, isError } = useAutomationExecutions();
  const executions = data?.executions ?? [];

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Automation</h1>
          <p className="text-sm text-slate-400">
            Recent rule-triggered automation executions.
          </p>
        </div>

        {isLoading && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
            Loading executions...
          </div>
        )}

        {isError && (
          <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            Failed to load executions.
          </div>
        )}

        <div className="space-y-3">
          {executions.map((execution) => (
            <div
              key={execution.id}
              className="flex items-center justify-between rounded-lg border border-white/10 bg-white/5 px-4 py-3"
            >
              <div>
                <p className="text-sm font-medium text-white">{execution.ruleName}</p>
                <p className="mt-1 text-xs text-slate-400">
                  {execution.action} · {execution.message}
                </p>
              </div>
              <span
                className={`rounded-full px-2 py-1 text-xs font-medium ${
                  execution.success
                    ? "bg-emerald-500/10 text-emerald-300"
                    : "bg-red-500/10 text-red-300"
                }`}
              >
                {execution.success ? "Success" : "Failed"}
              </span>
            </div>
          ))}

          {!isLoading && executions.length === 0 && (
            <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
              No automation executions yet.
            </div>
          )}
        </div>
      </div>
    </DashboardShell>
  );
}
