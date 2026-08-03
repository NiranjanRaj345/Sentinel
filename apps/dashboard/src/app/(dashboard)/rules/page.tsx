"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useRules } from "@/services/dashboard";

function formatCondition(condition: { field: string; operator: string; value: string }) {
  const operatorLabel = {
    equals: "==",
    not_equals: "!=",
    contains: "contains",
  }[condition.operator] ?? condition.operator;

  return `${condition.field} ${operatorLabel} ${condition.value}`;
}

export default function RulesRoute() {
  const { data, isLoading, isError, refetch } = useRules();
  const rules = data?.rules ?? [];

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Rules</h1>
          <p className="text-sm text-slate-400">
            Read-only view of active decision rules evaluated from events.
          </p>
        </div>

        {isLoading && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
            Loading rules...
          </div>
        )}

        {isError && (
          <div className="flex items-center justify-between rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            <span>Failed to load rules.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        <div className="space-y-3">
          {rules.map((rule) => (
            <div
              key={rule.id}
              className="rounded-lg border border-white/10 bg-white/5 p-4"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-sm font-semibold text-white">{rule.name}</h2>
                  <p className="mt-1 text-xs text-slate-400">ID: {rule.id}</p>
                </div>
                <span
                  className={`rounded-full px-2 py-1 text-xs font-medium ${
                    rule.enabled
                      ? "bg-emerald-500/10 text-emerald-300"
                      : "bg-slate-500/10 text-slate-300"
                  }`}
                >
                  {rule.enabled ? "Enabled" : "Disabled"}
                </span>
              </div>

              <div className="mt-4 grid gap-3 sm:grid-cols-2">
                <div className="rounded-md border border-white/10 bg-slate-950/40 p-3">
                  <p className="text-xs font-medium text-slate-300">Trigger</p>
                  <p className="mt-1 text-sm text-white">{rule.trigger}</p>
                </div>
                <div className="rounded-md border border-white/10 bg-slate-950/40 p-3">
                  <p className="text-xs font-medium text-slate-300">Actions</p>
                  <p className="mt-1 text-sm text-white">
                    {rule.actions.length > 0 ? rule.actions.join(", ") : "—"}
                  </p>
                </div>
              </div>

              {rule.conditions.length > 0 && (
                <div className="mt-4 rounded-md border border-white/10 bg-slate-950/40 p-3">
                  <p className="text-xs font-medium text-slate-300">Conditions</p>
                  <ul className="mt-2 space-y-1">
                    {rule.conditions.map((condition, index) => (
                      <li
                        key={index}
                        className="text-sm text-slate-200"
                      >
                        {formatCondition(condition)}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ))}

          {!isLoading && rules.length === 0 && (
            <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
              No rules configured yet.
            </div>
          )}
        </div>
      </div>
    </DashboardShell>
  );
}
