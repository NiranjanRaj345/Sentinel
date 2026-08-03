"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useNodes, useNodeDetail } from "@/services/dashboard";
import { useToastStore } from "@/stores/toast";
import type { Node } from "@/types/dashboard";
import Link from "next/link";
import { useParams } from "next/navigation";

const STATUS_COLORS: Record<string, string> = {
  online: "bg-emerald-500/10 text-emerald-300",
  offline: "bg-red-500/10 text-red-300",
  unknown: "bg-white/5 text-slate-300",
};

function NodeCard({ node }: { node: Node }) {
  const lastSeen = new Date(node.lastSeen);
  const timeAgo = Math.floor((Date.now() - lastSeen.getTime()) / 1000);
  const timeAgoText =
    timeAgo < 60 ? `${timeAgo}s ago` : `${Math.floor(timeAgo / 60)}m ago`;

  return (
    <Link
      href={`/nodes/${encodeURIComponent(node.id)}`}
      className="flex items-center justify-between rounded-lg border border-white/10 bg-white/5 px-4 py-3 transition-colors hover:bg-white/10"
    >
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-white/10 text-sm font-semibold text-white">
          {node.name.slice(0, 2).toUpperCase()}
        </div>
        <div>
          <p className="text-sm font-medium text-white">{node.name}</p>
          <p className="mt-1 text-xs text-slate-400">
            {node.hostname} · {node.platform} · v{node.version}
          </p>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <div className="text-right text-xs text-slate-400">
          <p>{timeAgoText}</p>
          <p className="mt-0.5">{node.address}</p>
        </div>
        <span
          className={`rounded-full px-2 py-1 text-xs font-medium ${
            STATUS_COLORS[node.status] ?? STATUS_COLORS.unknown
          }`}
        >
          {node.status}
        </span>
      </div>
    </Link>
  );
}

function NodeDetail({ nodeId }: { nodeId: string }) {
  const { data, isLoading, isError } = useNodeDetail(nodeId);
  const node = data as Node | undefined;

  if (isLoading) {
    return (
      <div className="space-y-3">
        <div className="h-24 animate-pulse rounded-lg bg-white/5" />
        <div className="h-24 animate-pulse rounded-lg bg-white/5" />
      </div>
    );
  }

  if (isError || !node) {
    return (
      <div className="rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
        Failed to load node details.
      </div>
    );
  }

  const lastSeen = new Date(node.lastSeen);
  const created = new Date(node.createdAt);

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-white/10 text-lg font-semibold text-white">
          {node.name.slice(0, 2).toUpperCase()}
        </div>
        <div>
          <h2 className="text-lg font-semibold text-white">{node.name}</h2>
          <p className="text-sm text-slate-400">{node.hostname}</p>
        </div>
      </div>

      <div className="grid gap-3 sm:grid-cols-2">
        <div className="rounded-lg border border-white/10 bg-white/5 p-3 text-sm text-slate-300">
          <p className="text-xs text-slate-400">Status</p>
          <p className="mt-1 capitalize">{node.status}</p>
        </div>
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
      </div>
    </div>
  );
}

export default function NodesRoute() {
  const nodesQuery = useNodes();
  const { data, isLoading, isError, refetch } = nodesQuery;
  const nodes = (data?.nodes ?? []) as Node[];
  const addToast = useToastStore((state) => state.addToast);
  const params = useParams();
  const selectedId = params?.id as string | undefined;

  React.useEffect(() => {
    if (isError) {
      addToast("Failed to load nodes", "error");
    }
  }, [isError, addToast]);

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Nodes</h1>
          <p className="text-sm text-slate-400">
            Managed devices registered with Mission Control.
          </p>
        </div>

        {isLoading && (
          <div className="space-y-3">
            {Array.from({ length: 3 }).map((_, index) => (
              <div key={index} className="h-20 animate-pulse rounded-lg bg-white/5" />
            ))}
          </div>
        )}

        {isError && (
          <div className="flex items-center justify-between rounded-lg border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-200">
            <span>Failed to load nodes.</span>
            <button
              onClick={() => refetch()}
              className="rounded-md bg-white/10 px-3 py-1.5 text-xs font-medium text-white hover:bg-white/20"
            >
              Retry
            </button>
          </div>
        )}

        {!isLoading && nodes.length === 0 && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4 text-sm text-slate-300">
            No nodes registered yet. Register a node to see it here.
          </div>
        )}

        <div className="space-y-3">
          {nodes.map((node: Node) => (
            <NodeCard key={node.id} node={node} />
          ))}
        </div>

        {selectedId && (
          <div className="rounded-lg border border-white/10 bg-white/5 p-4">
            <div className="mb-3 flex items-center justify-between">
              <h3 className="text-sm font-medium text-white">Node Detail</h3>
              <Link
                href="/nodes"
                className="text-xs text-slate-300 hover:text-white"
              >
                Close
              </Link>
            </div>
            <NodeDetail nodeId={selectedId} />
          </div>
        )}
      </div>
    </DashboardShell>
  );
}
