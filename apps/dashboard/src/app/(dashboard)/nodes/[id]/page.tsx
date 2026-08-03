"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useNodeDetail } from "@/services/dashboard";
import { useToastStore } from "@/stores/toast";
import type { Node } from "@/types/dashboard";
import Link from "next/link";
import { useParams } from "next/navigation";

const STATUS_COLORS: Record<string, string> = {
  online: "bg-emerald-500/10 text-emerald-300",
  offline: "bg-red-500/10 text-red-300",
  unknown: "bg-white/5 text-slate-300",
};

export default function NodeDetailRoute() {
  const params = useParams();
  const nodeId = params?.id as string | undefined;
  const { data, isLoading, isError, refetch } = useNodeDetail(nodeId ?? "");
  const node = data as Node | undefined;
  const addToast = useToastStore((state) => state.addToast);

  React.useEffect(() => {
    if (isError) {
      addToast("Failed to load node details", "error");
    }
  }, [isError, addToast]);

  if (isLoading) {
    return (
      <DashboardShell>
        <div className="space-y-3">
          <div className="h-8 w-48 animate-pulse rounded bg-white/5" />
          <div className="h-64 animate-pulse rounded-lg bg-white/5" />
        </div>
      </DashboardShell>
    );
  }

  if (isError || !node) {
    return (
      <DashboardShell>
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Link
              href="/nodes"
              className="text-sm text-slate-300 hover:text-white"
            >
              ← Nodes
            </Link>
          </div>
          <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            Failed to load node details.
            <button
              onClick={() => refetch()}
              className="ml-3 rounded-md bg-white/10 px-2 py-1 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        </div>
      </DashboardShell>
    );
  }

  const lastSeen = new Date(node.lastSeen);
  const created = new Date(node.createdAt);

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="flex items-center gap-2">
          <Link
            href="/nodes"
            className="text-sm text-slate-300 hover:text-white"
          >
            ← Nodes
          </Link>
        </div>

        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-white/10 text-lg font-semibold text-white">
            {node.name.slice(0, 2).toUpperCase()}
          </div>
          <div>
            <h1 className="text-2xl font-semibold text-white">{node.name}</h1>
            <p className="text-sm text-slate-400">{node.hostname}</p>
          </div>
          <span
            className={`ml-auto rounded-full px-2 py-1 text-xs font-medium ${
              STATUS_COLORS[node.status] ?? STATUS_COLORS.unknown
            }`}
          >
            {node.status}
          </span>
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          <div className="rounded-lg border border-white/10 bg-white/5 p-3 text-sm text-slate-300">
            <p className="text-xs text-slate-400">Platform</p>
            <p className="mt-1 capitalize">{node.platform}</p>
          </div>
          <div className="rounded-lg border border-white/10 bg-white/5 p-3 text-sm text-slate-300">
            <p className="text-xs text-slate-400">Version</p>
            <p className="mt-1">{node.version}</p>
          </div>
          <div className="rounded-lg border border-white/10 bg-white/5 p-3 text-sm text-slate-300">
            <p className="text-xs text-slate-400">Address</p>
            <p className="mt-1">{node.address}</p>
          </div>
          <div className="rounded-lg border border-white/10 bg-white/5 p-3 text-sm text-slate-300">
            <p className="text-xs text-slate-400">Last Seen</p>
            <p className="mt-1">{lastSeen.toLocaleString()}</p>
          </div>
          <div className="rounded-lg border border-white/10 bg-white/5 p-3 text-sm text-slate-300">
            <p className="text-xs text-slate-400">Registered</p>
            <p className="mt-1">{created.toLocaleString()}</p>
          </div>
          <div className="rounded-lg border border-white/10 bg-white/5 p-3 text-sm text-slate-300">
            <p className="text-xs text-slate-400">Node ID</p>
            <p className="mt-1 font-mono text-xs">{node.id}</p>
          </div>
        </div>
      </div>
    </DashboardShell>
  );
}
