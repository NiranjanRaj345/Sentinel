import { DashboardShell } from "@/components/layout/DashboardShell";

export default function HistoryRoute() {
  return (
    <DashboardShell>
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-white">History</h1>
        <p className="text-sm text-slate-400">
          Time-series history from SQLite will be surfaced here.
        </p>
      </div>
    </DashboardShell>
  );
}
